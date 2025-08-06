//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
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
	"io/fs"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync"
)

var (
	apache = template.Must(template.New("Apache License").Parse(`{{.Escape}}
{{.Escape}} Copyright (C) 2019-{{.Year}} {{.Maintainer}}
{{.Escape}}
{{.Escape}} Licensed under the Apache License, Version 2.0 (the "License");
{{.Escape}} You may not use this file except in compliance with the License.
{{.Escape}} You may obtain a copy of the License at
{{.Escape}}
{{.Escape}}    https://www.apache.org/licenses/LICENSE-2.0
{{.Escape}}
{{.Escape}} Unless required by applicable law or agreed to in writing, software
{{.Escape}} distributed under the License is distributed on an "AS IS" BASIS,
{{.Escape}} WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
{{.Escape}} See the License for the specific language governing permissions and
{{.Escape}} limitations under the License.
{{.Escape}}
`))
	docker = template.Must(template.New("Apache License").Parse(`{{.Escape}} syntax = docker/dockerfile:latest
{{.Escape}} check=error=true
{{.Escape}}
{{.Escape}} Copyright (C) 2019-{{.Year}} {{.Maintainer}}
{{.Escape}}
{{.Escape}} Licensed under the Apache License, Version 2.0 (the "License");
{{.Escape}} You may not use this file except in compliance with the License.
{{.Escape}} You may obtain a copy of the License at
{{.Escape}}
{{.Escape}}    https://www.apache.org/licenses/LICENSE-2.0
{{.Escape}}
{{.Escape}} Unless required by applicable law or agreed to in writing, software
{{.Escape}} distributed under the License is distributed on an "AS IS" BASIS,
{{.Escape}} WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
{{.Escape}} See the License for the specific language governing permissions and
{{.Escape}} limitations under the License.
{{.Escape}}
`))

	googleProtoApache = template.Must(template.New("Google Proto Apache License").Parse(`{{.Escape}}
{{.Escape}} Copyright (C) {{.Year}} Google LLC
{{.Escape}} Modified by {{.Maintainer}}
{{.Escape}}
{{.Escape}} Licensed under the Apache License, Version 2.0 (the "License");
{{.Escape}} You may not use this file except in compliance with the License.
{{.Escape}} You may obtain a copy of the License at
{{.Escape}}
{{.Escape}}    https://www.apache.org/licenses/LICENSE-2.0
{{.Escape}}
{{.Escape}} Unless required by applicable law or agreed to in writing, software
{{.Escape}} distributed under the License is distributed on an "AS IS" BASIS,
{{.Escape}} WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
{{.Escape}} See the License for the specific language governing permissions and
{{.Escape}} limitations under the License.
{{.Escape}}
`))

	goStandard = template.Must(template.New("Go License").Parse(`{{.Escape}}
{{.Escape}} Copyright (c) 2009-{{.Year}} The Go Authors. All rights resered.
{{.Escape}} Modified by {{.Maintainer}}
{{.Escape}}
{{.Escape}} Redistribution and use in source and binary forms, with or without
{{.Escape}} modification, are permitted provided that the following conditions are
{{.Escape}} met:
{{.Escape}}
{{.Escape}}    * Redistributions of source code must retain the above copyright
{{.Escape}} notice, this list of conditions and the following disclaimer.
{{.Escape}}    * Redistributions in binary form must reproduce the above
{{.Escape}} copyright notice, this list of conditions and the following disclaimer
{{.Escape}} in the documentation and/or other materials provided with the
{{.Escape}} distribution.
{{.Escape}}    * Neither the name of Google Inc. nor the names of its
{{.Escape}} contributors may be used to endorse or promote products derived from
{{.Escape}} this software without specific prior written permission.
{{.Escape}}
{{.Escape}} THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
{{.Escape}} "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
{{.Escape}} LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
{{.Escape}} A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
{{.Escape}} OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
{{.Escape}} SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
{{.Escape}} LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
{{.Escape}} DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
{{.Escape}} THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
{{.Escape}} (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
{{.Escape}} OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
{{.Escape}}
`))

	slushEscape = "//"
	sharpEscape = "#"
)

type Data struct {
	Escape     string
	Maintainer string
	Year       int
}

const (
	minimumArgumentLength = 2
	defaultMaintainer     = "vdaas.org vald team <vald@vdaas.org>"
	maintainerKey         = "MAINTAINER"
	yearKey               = "YEAR"
)

func main() {
	log.Init()
	if len(os.Args) < minimumArgumentLength {
		// skipcq: RVV-A0003
		log.Fatal(errors.New("invalid argument"))
	}
	for _, path := range dirwalk(os.Args[1]) {
		fmt.Println(path)
		err := readAndRewrite(path)
		if err != nil {
			// skipcq: RVV-A0003
			log.Fatal(err)
		}
	}
}

func dirwalk(dir string) []string {
	files, err := file.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	var paths []string
	for _, f := range files {
		if f.IsDir() {
			if !strings.Contains(f.Name(), "vendor") &&
				!strings.Contains(f.Name(), "versions") &&
				!strings.Contains(f.Name(), ".git") &&
				!strings.Contains(f.Name(), "target") ||
				strings.HasPrefix(f.Name(), ".github") {
				paths = append(paths, dirwalk(file.Join(dir, f.Name()))...)
			}
			continue
		}
		switch filepath.Ext(f.Name()) {
		case
			".ai",
			".all-contributorsrc",
			".cfg",
			".crt",
			".default",
			".drawio",
			".git",
			".gitignore",
			".gitkeep",
			".gitmodules",
			".gotmpl",
			".hdf5",
			".helmignore",
			".html",
			".json",
			".key",
			".kvsdb",
			".lock",
			".md",
			".md5",
			".mod",
			".pdf",
			".pem",
			".png",
			".ssv",
			".sum",
			".svg",
			".tmpl",
			".tpl",
			".txt",
			".webp",
			".whitesource",
			"LICENSE",
			"Pipefile":
		default:
			switch f.Name() {
			case
				"AUTHORS",
				"CONTRIBUTORS",
				"FAISS_VERSION",
				"GO_VERSION",
				"NGT_VERSION",
				"Pipefile",
				"VALD_VERSION",
				"grp",
				"obj",
				"prf",
				"rust-toolchain",
				"src",
				"tre":
			default:
				path := file.Join(dir, f.Name())
				log.Info(path)
				paths = append(paths, path)
			}
		}
	}
	return paths
}

func isSymlink(path string) (bool, error) {
	lst, err := os.Lstat(path)
	if err != nil {
		return false, err
	}
	return lst.Mode()&os.ModeSymlink == os.ModeSymlink, nil
}

func readAndRewrite(path string) error {
	// return if it is a symlink
	isSym, err := isSymlink(path)
	if err != nil {
		return err
	}
	if isSym {
		return nil
	}

	f, err := os.OpenFile(path, os.O_RDWR|os.O_SYNC, fs.ModePerm)
	if err != nil {
		return errors.Errorf("filepath %s, could not open", path)
	}
	fi, err := f.Stat()
	if err != nil {
		err = f.Close()
		if err != nil {
			// skipcq: RVV-A0003
			log.Fatal(err)
		}
		return errors.Errorf("filepath %s, could not open", path)
	}
	buf := bytes.NewBuffer(make([]byte, 0, fi.Size()))
	maintainer := os.Getenv(maintainerKey)
	if maintainer == "" {
		maintainer = defaultMaintainer
	}
	var year int
	if yearString := os.Getenv(yearKey); yearString == "" {
		year = time.Now().Year()
	} else {
		y, err := time.Parse("2006", yearString)
		if err != nil {
			// skipcq: RVV-A0003
			log.Fatal(err)
		}
		year = y.Year()
	}
	d := Data{
		Maintainer: maintainer,
		Year:       year,
		Escape:     sharpEscape,
	}
	if fi.Name() == "LICENSE" {
		err = license.Execute(buf, d)
		if err != nil {
			// skipcq: RVV-A0003
			log.Fatal(err)
		}
	} else {
		tmpl := apache
		switch filepath.Ext(path) {
		case ".go", ".c", ".h", ".hpp", ".cpp":
			d.Escape = slushEscape
			switch fi.Name() {
			case "errgroup_test.go",
				"singleflight.go",
				"semaphore.go",
				"semaphore_bench_test.go",
				"semaphore_example_test.go",
				"semaphore_test.go":
				tmpl = goStandard
			case "error_details.pb.go",
				"error_details.pb.json.go",
				"error_details_vtproto.pb.go":
				tmpl = googleProtoApache
			default:
			}
		case ".proto":
			if fi.Name() == "error_details.proto" {
				tmpl = googleProtoApache
			}
			d.Escape = slushEscape
		case ".rs":
			d.Escape = slushEscape
		default:
			if fi.Name() == "Dockerfile" {
				tmpl = docker
			}
		}
		lf := true
		bf := false
		sc := bufio.NewScanner(f)
		once := sync.Once{}
		for sc.Scan() {
			line := sc.Text()
			if filepath.Ext(path) == ".go" && strings.HasPrefix(line, "//go:") ||
				filepath.Ext(path) == ".py" && strings.HasPrefix(line, "# -*-") ||
				filepath.Ext(path) == ".sh" && strings.HasPrefix(line, "#!") ||
				filepath.Ext(path) == ".yaml" && strings.HasPrefix(line, "# !") ||
				filepath.Ext(path) == ".yml" && strings.HasPrefix(line, "# !") {
				bf = true
				_, err = buf.WriteString(line)
				if err != nil {
					// skipcq: RVV-A0003
					log.Fatal(err)
				}
				_, err = buf.WriteString("\n")
				if err != nil {
					// skipcq: RVV-A0003
					log.Fatal(err)
				}
				_, err = buf.WriteString("\n")
				if err != nil {
					// skipcq: RVV-A0003
					log.Fatal(err)
				}
				continue
			}
			if (filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml") && strings.HasPrefix(line, "---") {
				_, err = buf.WriteString(line)
				if err != nil {
					// skipcq: RVV-A0003
					log.Fatal(err)
				}
				_, err = buf.WriteString("\n")
				if err != nil {
					// skipcq: RVV-A0003
					log.Fatal(err)
				}
				continue
			}
			if lf && strings.HasPrefix(line, d.Escape) {
				continue
			} else if !bf {
				once.Do(func() {
					err = tmpl.Execute(buf, d)
					if err != nil {
						// skipcq: RVV-A0003
						log.Fatal(err)
					}
				})
				lf = false
			}
			if !lf {
				_, err = buf.WriteString(line)
				if err != nil {
					// skipcq: RVV-A0003
					log.Fatal(err)
				}
				_, err = buf.WriteString("\n")
				if err != nil {
					// skipcq: RVV-A0003
					log.Fatal(err)
				}
			}
			bf = false
		}
	}
	err = f.Close()
	if err != nil {
		return errors.Errorf("filepath %s, could not close", path)
	}
	err = os.RemoveAll(path)
	if err != nil {
		return errors.Errorf("filepath %s, could not delete", path)
	}
	f, err = os.Create(path)
	if err != nil {
		err = f.Close()
		if err != nil {
			// skipcq: RVV-A0003
			log.Fatal(err)
		}
		return errors.Errorf("filepath %s, could not open", path)
	}
	_, err = f.WriteString(strings.ReplaceAll(buf.String(), d.Escape+"\n\n\n", d.Escape+"\n\n"))
	if err != nil {
		// skipcq: RVV-A0003
		log.Fatal(err)
	}
	err = f.Close()
	if err != nil {
		// skipcq: RVV-A0003
		log.Fatal(err)
	}
	return nil
}

var license = template.Must(template.New("LICENSE").Parse(
	"                                 Apache License\n                           Version 2.0, January 2004\n                        https://www.apache.org/licenses/\n\n   TERMS AND CONDITIONS FOR USE, REPRODUCTION, AND DISTRIBUTION\n\n   1. Definitions.\n\n      \"License\" shall mean the terms and conditions for use, reproduction,\n      and distribution as defined by Sections 1 through 9 of this document.\n\n      \"Licensor\" shall mean the copyright owner or entity authorized by\n      the copyright owner that is granting the License.\n\n      \"Legal Entity\" shall mean the union of the acting entity and all\n      other entities that control, are controlled by, or are under common\n      control with that entity. For the purposes of this definition,\n      \"control\" means (i) the power, direct or indirect, to cause the\n      direction or management of such entity, whether by contract or\n      otherwise, or (ii) ownership of fifty percent (50%) or more of the\n      outstanding shares, or (iii) beneficial ownership of such entity.\n\n      \"You\" (or \"Your\") shall mean an individual or Legal Entity\n      exercising permissions granted by this License.\n\n      \"Source\" form shall mean the preferred form for making modifications,\n      including but not limited to software source code, documentation\n      source, and configuration files.\n\n      \"Object\" form shall mean any form resulting from mechanical\n      transformation or translation of a Source form, including but\n      not limited to compiled object code, generated documentation,\n      and conversions to other media types.\n\n      \"Work\" shall mean the work of authorship, whether in Source or\n      Object form, made available under the License, as indicated by a\n      copyright notice that is included in or attached to the work\n      (an example is provided in the Appendix below).\n\n      \"Derivative Works\" shall mean any work, whether in Source or Object\n      form, that is based on (or derived from) the Work and for which the\n      editorial revisions, annotations, elaborations, or other modifications\n      represent, as a whole, an original work of authorship. For the purposes\n      of this License, Derivative Works shall not include works that remain\n      separable from, or merely link (or bind by name) to the interfaces of,\n      the Work and Derivative Works thereof.\n\n      \"Contribution\" shall mean any work of authorship, including\n      the original version of the Work and any modifications or additions\n      to that Work or Derivative Works thereof, that is intentionally\n      submitted to Licensor for inclusion in the Work by the copyright owner\n      or by an individual or Legal Entity authorized to submit on behalf of\n      the copyright owner. For the purposes of this definition, \"submitted\"\n      means any form of electronic, verbal, or written communication sent\n      to the Licensor or its representatives, including but not limited to\n      communication on electronic mailing lists, source code control systems,\n      and issue tracking systems that are managed by, or on behalf of, the\n      Licensor for the purpose of discussing and improving the Work, but\n      excluding communication that is conspicuously marked or otherwise\n      designated in writing by the copyright owner as \"Not a Contribution.\"\n\n      \"Contributor\" shall mean Licensor and any individual or Legal Entity\n      on behalf of whom a Contribution has been received by Licensor and\n      subsequently incorporated within the Work.\n\n   2. Grant of Copyright License. Subject to the terms and conditions of\n      this License, each Contributor hereby grants to You a perpetual,\n      worldwide, non-exclusive, no-charge, royalty-free, irrevocable\n      copyright license to reproduce, prepare Derivative Works of,\n      publicly display, publicly perform, sublicense, and distribute the\n      Work and such Derivative Works in Source or Object form.\n\n   3. Grant of Patent License. Subject to the terms and conditions of\n      this License, each Contributor hereby grants to You a perpetual,\n      worldwide, non-exclusive, no-charge, royalty-free, irrevocable\n      (except as stated in this section) patent license to make, have made,\n      use, offer to sell, import, and otherwise transfer the Work,\n      where such license applies only to those patent claims licensable\n      by such Contributor that are necessarily infringed by their\n      Contribution(s) alone or by combination of their Contribution(s)\n      with the Work to which such Contribution(s) was submitted. If You\n      institute patent litigation against any entity (including a\n      cross-claim or counterclaim in a lawsuit) alleging that the Work\n      or a Contribution incorporated within the Work constitutes direct\n      or contributory patent infringement, then any patent licenses\n      granted to You under this License for that Work shall terminate\n      as of the date such litigation is filed.\n\n   4. Redistribution. You may reproduce and distribute copies of the\n      Work or Derivative Works thereof in any medium, with or without\n      modifications, and in Source or Object form, provided that You\n      meet the following conditions:\n\n      (a) You must give any other recipients of the Work or\n          Derivative Works a copy of this License; and\n\n      (b) You must cause any modified files to carry prominent notices\n          stating that You changed the files; and\n\n      (c) You must retain, in the Source form of any Derivative Works\n          that You distribute, all copyright, patent, trademark, and\n          attribution notices from the Source form of the Work,\n          excluding those notices that do not pertain to any part of\n          the Derivative Works; and\n\n      (d) If the Work includes a \"NOTICE\" text file as part of its\n          distribution, then any Derivative Works that You distribute must\n          include a readable copy of the attribution notices contained\n          within such NOTICE file, excluding those notices that do not\n          pertain to any part of the Derivative Works, in at least one\n          of the following places: within a NOTICE text file distributed\n          as part of the Derivative Works; within the Source form or\n          documentation, if provided along with the Derivative Works; or,\n          within a display generated by the Derivative Works, if and\n          wherever such third-party notices normally appear. The contents\n          of the NOTICE file are for informational purposes only and\n          do not modify the License. You may add Your own attribution\n          notices within Derivative Works that You distribute, alongside\n          or as an addendum to the NOTICE text from the Work, provided\n          that such additional attribution notices cannot be construed\n          as modifying the License.\n\n      You may add Your own copyright statement to Your modifications and\n      may provide additional or different license terms and conditions\n      for use, reproduction, or distribution of Your modifications, or\n      for any such Derivative Works as a whole, provided Your use,\n      reproduction, and distribution of the Work otherwise complies with\n      the conditions stated in this License.\n\n   5. Submission of Contributions. Unless You explicitly state otherwise,\n      any Contribution intentionally submitted for inclusion in the Work\n      by You to the Licensor shall be under the terms and conditions of\n      this License, without any additional terms or conditions.\n      Notwithstanding the above, nothing herein shall supersede or modify\n      the terms of any separate license agreement you may have executed\n      with Licensor regarding such Contributions.\n\n   6. Trademarks. This License does not grant permission to use the trade\n      names, trademarks, service marks, or product names of the Licensor,\n      except as required for reasonable and customary use in describing the\n      origin of the Work and reproducing the content of the NOTICE file.\n\n   7. Disclaimer of Warranty. Unless required by applicable law or\n      agreed to in writing, Licensor provides the Work (and each\n      Contributor provides its Contributions) on an \"AS IS\" BASIS,\n      WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or\n      implied, including, without limitation, any warranties or conditions\n      of TITLE, NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A\n      PARTICULAR PURPOSE. You are solely responsible for determining the\n      appropriateness of using or redistributing the Work and assume any\n      risks associated with Your exercise of permissions under this License.\n\n   8. Limitation of Liability. In no event and under no legal theory,\n      whether in tort (including negligence), contract, or otherwise,\n      unless required by applicable law (such as deliberate and grossly\n      negligent acts) or agreed to in writing, shall any Contributor be\n      liable to You for damages, including any direct, indirect, special,\n      incidental, or consequential damages of any character arising as a\n      result of this License or out of the use or inability to use the\n      Work (including but not limited to damages for loss of goodwill,\n      work stoppage, computer failure or malfunction, or any and all\n      other commercial damages or losses), even if such Contributor\n      has been advised of the possibility of such damages.\n\n   9. Accepting Warranty or Additional Liability. While redistributing\n      the Work or Derivative Works thereof, You may choose to offer,\n      and charge a fee for, acceptance of support, warranty, indemnity,\n      or other liability obligations and/or rights consistent with this\n      License. However, in accepting such obligations, You may act only\n      on Your own behalf and on Your sole responsibility, not on behalf\n      of any other Contributor, and only if You agree to indemnify,\n      defend, and hold each Contributor harmless for any liability\n      incurred by, or claims asserted against, such Contributor by reason\n      of your accepting any such warranty or additional liability.\n\n   END OF TERMS AND CONDITIONS\n\n   APPENDIX: How to apply the Apache License to your work.\n\n      To apply the Apache License to your work, attach the following\n      boilerplate notice, with the fields enclosed by brackets \"[]\"\n      replaced with your own identifying information. (Don't include\n      the brackets!)  The text should be enclosed in the appropriate\n      comment syntax for the file format. We also recommend that a\n      file or class name and description of purpose be included on the\n      same \"printed page\" as the copyright notice for easier\n      identification within third-party archives.\n\n   Copyright (C) 2019-{{.Year}} {{.Maintainer}}\n\n   Licensed under the Apache License, Version 2.0 (the \"License\");\n   You may not use this file except in compliance with the License.\n   You may obtain a copy of the License at\n\n       https://www.apache.org/licenses/LICENSE-2.0\n\n   Unless required by applicable law or agreed to in writing, software\n   distributed under the License is distributed on an \"AS IS\" BASIS,\n   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.\n   See the License for the specific language governing permissions and\n   limitations under the License."),
)
