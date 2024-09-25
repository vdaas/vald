package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unsafe"

	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/http/client"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

const (
	BASE_URL      = "https://vald.vdaas.org"
	PREFIX_PROP   = `property="og:url" content="`
	PREFIX_SRC    = "src="
	PREFIX_SRCSET = "srcset="
	PREFIX_HREF   = "href="

	BASE_REGEXP = `[\w!\?/\+\-_~=;:\.,\*&@#\$%\(\)'\[\]]+`
	URL_REGEXP  = `https?://[\w!\?/\+\-_~=;\.,\*&@#\$%\(\)'\[\]]+`
)

var (
	format     = flag.String("format", "html", "file format(html)")
	path       = flag.String("path", "./", "directory or file path")
	ignorePath = flag.String("ignore-path", "", "ignore path to check")

	ignoreLinks = []string{
		"javascript:void(0)",
		"mailto:vald@vdaas.org",
		"https://github.com",
		"https://twitter.com/vdaas_vald",
		"https://vdaas-vald.medium.com",
	}

	reProp   = regexp.MustCompile(PREFIX_PROP + BASE_REGEXP)
	reSrc    = regexp.MustCompile(PREFIX_SRC + BASE_REGEXP)
	reSrcSet = regexp.MustCompile(PREFIX_SRCSET + BASE_REGEXP)
	reHref   = regexp.MustCompile(PREFIX_HREF + BASE_REGEXP)

	url = ""
)

func getFiles(dir string) []string {
	filePaths := []string{}
	reIgnore := regexp.MustCompile(*ignorePath)
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(info.Name(), *format) || reIgnore.FindString(path) != "" {
			return nil
		}
		filePaths = append(filePaths, path)
		return nil
	})
	if err != nil {
		log.Error(err.Error())
	}
	return filePaths
}

func convertToURL(s string) string {
	b := bytes.NewBuffer(make([]byte, 0, 100))
	if strings.HasPrefix(s, "#") {
		b.WriteString(url)
		b.WriteString(s)
		return b.String()
	}
	if strings.HasPrefix(s, "/") {
		b.WriteString(BASE_URL)
		b.WriteString(s)
		return b.String()
	}
	if strings.HasPrefix(s, ".") {
		b.WriteString(BASE_URL)
		b.WriteString("/docs/")
		b.WriteString(s)
		return b.String()
	}
	return s
}

func isBlackList(url string, bList []string) bool {
	for _, list := range bList {
		if strings.Contains(url, list) {
			return true
		}
	}
	return false
}

func exec(url string, cli *http.Client) int {
	resp, err := cli.Get(url)
	if err != nil {
		log.Error(err)
		return -1
	}
	defer resp.Body.Close()
	return resp.StatusCode
}

func main() {
	flag.Parse()

	// file paths
	var paths []string
	info, err := os.Stat(*path)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	if info.IsDir() {
		paths = getFiles(*path)
	} else {
		paths = []string{*path}
	}

	// check
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cli, err := client.New(
		client.WithIdleConnTimeout("2s"),
		// client.WithResponseHeaderTimeout("2s"),
	)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	eg, _ := errgroup.New(ctx)
	// eg.SetLimit(egLimit)
	countAll := 0
	successAll := 0
	mu := sync.Mutex{}
	for _, path := range paths {
		b, err := file.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}
		str := *(*string)(unsafe.Pointer(&b))
		// get origin url
		url = strings.TrimPrefix(reProp.FindString(str), PREFIX_PROP)

		count := 0
		success := 0

		// get links and convert to correct URL
		urls := []map[string]string{}
		u := reSrc.FindAllString(str, -1)
		for _, elem := range u {
			e := strings.TrimPrefix(elem, PREFIX_SRC)
			url := convertToURL(e)
			if !isBlackList(url, ignoreLinks) {
				urls = append(urls, map[string]string{e: url})
				count++
			}
		}
		u = reHref.FindAllString(str, -1)
		for _, elem := range u {
			e := strings.TrimPrefix(elem, PREFIX_HREF)
			url := convertToURL(e)
			if !isBlackList(url, ignoreLinks) {
				urls = append(urls, map[string]string{e: url})
				count++
			}
		}
		u = reSrcSet.FindAllString(str, -1)
		for _, elem := range u {
			e := strings.TrimPrefix(elem, PREFIX_SRCSET)
			url := convertToURL(e)
			if !isBlackList(url, ignoreLinks) {
				urls = append(urls, map[string]string{e: url})
				count++
			}
		}
		errUrls := map[string]int{}
		fmt.Printf("%s (url: %s)\n", path, url)
		for _, url := range urls {
			eg.Go(func() error {
				for k, v := range url {
					mu.Lock()
					if _, ok := errUrls[k]; ok {
						continue
					}
					mu.Unlock()
					code := exec(v, cli)
					if code == 200 {
						success++
					} else {
						log.Warnf("[%d] %s", code, v)
						mu.Lock()
						errUrls[k] = code
						mu.Unlock()
					}
				}
				return nil
			})
		}
		err = eg.Wait()
		if err != nil {
			log.Error(err.Error())
		}
		if len(errUrls) > 0 {
			fmt.Printf("error links:\n")
			for ref, code := range errUrls {
				fmt.Printf("%s => %d\n", ref, code)
			}
		}
		countAll += count
		successAll += success
		fmt.Printf("[result] all: %d, ok: %d, fail: %d\n\n", count, success, count-success)
	}
	fmt.Printf("\n[summary] all: %d, OK: %d, NG: %d", countAll, successAll, countAll-successAll)
	return
}
