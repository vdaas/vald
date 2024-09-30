// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
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
	"sync/atomic"

	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net"
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
	format      = flag.String("format", "html", "file format(html)")
	path        = flag.String("path", "./", "directory or file path")
	ignorePath  = flag.String("ignore-path", "", "ignore path to check")
	ignoreLinks = []string{
		"javascript:void(0)",
		"mailto:vald@vdaas.org",
		"https://github.com",
		"https://twitter.com/vdaas_vald",
		"https://vdaas-vald.medium.com",
		"https://join.slack.com/t/vald-community",
	}

	reProp   = regexp.MustCompile(PREFIX_PROP + BASE_REGEXP)
	reSrc    = regexp.MustCompile(PREFIX_SRC + BASE_REGEXP)
	reSrcSet = regexp.MustCompile(PREFIX_SRCSET + BASE_REGEXP)
	reHref   = regexp.MustCompile(PREFIX_HREF + BASE_REGEXP)
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

func convertToURL(s, url string) string {
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
		log.Errorf(err.Error())
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
		return -1
	}
	resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return resp.StatusCode
	}
	return 200
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

	der, err := net.NewDialer(
		net.WithEnableDNSCache(),
		net.WithEnableDialerDualStack(),
		net.WithDNSCacheExpiration("1h"),
	)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	cli, err := client.New(
		// enable http2
		client.WithForceAttemptHTTP2(true),
		client.WithEnableKeepalives(false),
		// stream max connection
		client.WithMaxConnsPerHost(len(paths)*30),
		client.WithMaxIdleConnsPerHost(len(paths)*30),
		client.WithMaxIdleConns(len(paths)*30),
		client.WithDialContext(der.GetDialer()),
	)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	eg, _ := errgroup.New(ctx)

	countAll := 0
	successAll := 0
	failAll := 0
	mu := sync.Mutex{}

	type result struct {
		url      string
		count    int
		success  int
		fail     int
		errLinks map[string]int
	}
	// result for each file path
	res := map[string]result{}
	// map of external link to avoid DOS
	exLinks := map[string]int{}
	for _, path := range paths {
		b, err := file.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}
		// str := *(*string)(unsafe.Pointer(&b))
		str := string(b)
		// get origin url
		originURL := strings.TrimPrefix(reProp.FindString(str), PREFIX_PROP)
		// init counter
		r := result{
			url:      originURL,
			count:    0,
			success:  0,
			fail:     0,
			errLinks: map[string]int{},
		}
		// Get links and convert to correct URL
		// Key is raw link written in file, value is correct link expression.
		urls := []map[string]string{}
		u := reSrc.FindAllString(str, -1)
		for _, elem := range u {
			e := strings.TrimPrefix(elem, PREFIX_SRC)
			targetUrl := convertToURL(e, originURL)
			if !isBlackList(targetUrl, ignoreLinks) {
				urls = append(urls, map[string]string{e: targetUrl})
				r.count++
			}
		}
		u = reHref.FindAllString(str, -1)
		for _, elem := range u {
			e := strings.TrimPrefix(elem, PREFIX_HREF)
			targetUrl := convertToURL(e, originURL)
			if !isBlackList(targetUrl, ignoreLinks) {
				urls = append(urls, map[string]string{e: targetUrl})
				r.count++
			}
		}
		u = reSrcSet.FindAllString(str, -1)
		for _, elem := range u {
			e := strings.TrimPrefix(elem, PREFIX_SRCSET)
			targetUrl := convertToURL(e, originURL)
			if !isBlackList(targetUrl, ignoreLinks) {
				urls = append(urls, map[string]string{e: targetUrl})
				r.count++
			}
		}
		fmt.Printf("checking...%s (url: %s)\n", path, originURL)
		var (
			success int32 = 0
			fail    int32 = 0
		)
		for _, url := range urls {
			eg.Go(func() error {
				for k, v := range url {
					mu.Lock()
					if _, ok := r.errLinks[k]; ok {
						atomic.AddInt32(&fail, 1)
						mu.Unlock()
						continue
					}
					var code int
					if strings.Contains(v, BASE_URL) {
						code = exec(v, cli)
					} else {
						if c, ok := exLinks[v]; ok {
							code = c
						} else {
							code = exec(v, cli)
							exLinks[v] = code
						}
						mu.Unlock()
					}
					if code == 200 {
						atomic.AddInt32(&success, 1)
					} else {
						log.Warnf("[%d] %s", code, v)
						r.errLinks[k] = code
						atomic.AddInt32(&fail, 1)
					}
				}
				return nil
			})
		}
		err = eg.Wait()
		if err != nil {
			log.Error(err.Error())
		}
		r.success = int(atomic.LoadInt32(&success))
		r.fail = int(atomic.LoadInt32(&fail))

		countAll += r.count
		successAll += r.success
		failAll += r.fail
		res[path] = r
	}

	for k, v := range res {
		fmt.Printf("[%s]\n%s\n", k, v.url)
		for link, code := range v.errLinks {
			fmt.Printf("%s => %d\n", link, code)
		}
		fmt.Printf("count: %d, ok: %d, fail: %d\n\n", v.count, v.success, v.fail)
	}
	fmt.Printf("\n[summary] all: %d, OK: %d, NG: %d\n", countAll, successAll, failAll)
	if countAll != successAll {
		os.Exit(1)
	}
	return
}
