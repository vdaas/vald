//
// Copyright (C) 2019-2019 kpango (Yusuke Kato)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
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
	fi, err := f.Stat()
	if err != nil {
		f.Close()
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
	case ".go", ".proto":
		d.Escape = slushEscape
	}

	lf := true
	bf := false
	sc := bufio.NewScanner(f)
	once := sync.Once{}
	for sc.Scan() {
		line := sc.Text()
		if filepath.Ext(path) == ".go" && strings.HasPrefix(line, "// +build") ||
			filepath.Ext(path) == ".py" && strings.HasPrefix(line, "# -*-") {
			bf = true
			buf.WriteString(line)
			buf.WriteString("\n")
			buf.WriteString("\n")
			continue
		}

		if lf && strings.HasPrefix(line, d.Escape) {
			continue
		} else if !bf {
			once.Do(func() {
				apache.Execute(buf, d)
			})
			lf = false
		}
		if !lf {
			buf.WriteString(line)
			buf.WriteString("\n")
		}
		bf = false
	}
	f.Close()
	os.RemoveAll(path)
	f, err = os.Create(path)
	if err != nil {
		f.Close()
		return errors.Errorf("filepath %s, could not open", path)
	}
	f.Write(buf.Bytes())
	f.Close()

	return nil
}
