package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
)

var (
	apache = template.Must(template.New("Apache License").Parse(`{{.Escape}}
{{.Escape}} Copyright (C) 2019-{{.Year}} {{.NickName}} ({{.FullName}})
{{.Escape}}
{{.Escape}} Licensed under the Apache License, Version 2.0 (the "License");
{{.Escape}} you may not use this file except in compliance with the License.
{{.Escape}} You may obtain a copy of the License at
{{.Escape}}
{{.Escape}}    http://www.apache.org/licenses/LICENSE-2.0
{{.Escape}}
{{.Escape}} Unless required by applicable law or agreed to in writing, software
{{.Escape}} distributed under the License is distributed on an "AS IS" BASIS,
{{.Escape}} WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
{{.Escape}} See the License for the specific language governing permissions and
{{.Escape}} limitations under the License.
{{.Escape}}
`))
	slushEscape = "//"
	sharpEscape = "#"
)

type Data struct {
	Escape   string
	NickName string
	FullName string
	Year     int
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal(errors.New("invalid argument"))
	}
	for _, path := range dirwalk(os.Args[1]) {
		fmt.Println(path)
		readAndRewrite(path)
	}
}

func dirwalk(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var paths []string
	for _, file := range files {
		if file.IsDir() {
			if !strings.Contains(file.Name(), "vendor") && !strings.Contains(file.Name(), ".git") {
				paths = append(paths, dirwalk(filepath.Join(dir, file.Name()))...)
			}
			continue
		}
		switch filepath.Ext(file.Name()) {
		case
			".cfg",
			".drawio",
			".git",
			".gitignore",
			".gitkeep",
			".gitmodules",
			".helmignore",
			".html",
			".json",
			".lock",
			".md",
			".md5",
			".mod",
			".png",
			".sum",
			".ssv",
			".svg",
			".tpl",
			".txt",
			".whitesource",
			"LICENSE",
			"Pipefile":
		default:
			switch file.Name() {
			case
				"LICENSE",
				"Pipefile",
				"grp",
				"src",
				"obj",
				"prf",
				"tre":
			default:
				path, err := filepath.Abs(filepath.Join(dir, file.Name()))
				if err != nil {
					log.Fatal("error")
				}
				paths = append(paths, path)
			}
		}
	}

	return paths
}

func readAndRewrite(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return errors.Errorf("filepath %s, could not open", path)
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return errors.Errorf("filepath %s, could not open", path)
	}
	buf := bytes.NewBuffer(make([]byte, 0, fi.Size()))
	d := Data{
		NickName: "kpango",
		FullName: "Yusuke Kato",
		Year:     time.Now().Year(),
		Escape:   sharpEscape,
	}
	switch filepath.Ext(path) {
	case ".go", ".py":
		d.Escape = slushEscape
	}
	apache.Execute(buf, d)

	lf := true
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		if lf && strings.HasPrefix(line, d.Escape) {
			continue
		} else {
			lf = false
		}
		if !lf {
			buf.WriteString(line)
			buf.WriteString("\n")
		}
	}
	f.Truncate(fi.Size())
	io.Copy(f, buf)

	return nil
}
