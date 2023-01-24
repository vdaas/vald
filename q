#
# Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

[1;4;33mCHANGELOG.md[0m

[36m───[0m[36m┐[0m
[34m1[0m: [36m│[0m
[36m───[0m[36m┴[0m[36m──────────────────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[38;5;208m# CHANGELOG[0m                                          [34m│ [38;5;208m# CHANGELOG[0m
[34m[0m                                                     [34m│ [0m
[34m[48;5;52;38;5;208m## v1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                                            [0m[34m│ [48;5;22;38;5;208m## v1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[0m                                                     [34m│ [0m
[34m[38;5;208m### Docker images[0m                                    [34m│ [38;5;208m### Docker images[0m
[34m[0m                                                     [34m│ [0m

[36m────[0m[36m┐[0m
[34m12[0m: [36m│[0m
[36m────[0m[36m┴[0m[36m─────────────────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[38;5;231m  <[38;5;203mtr[38;5;231m>[0m                                               [34m│ [38;5;231m  <[38;5;203mtr[38;5;231m>[0m
[34m[38;5;231m    <[38;5;203mtd[38;5;231m>Agent NGT</[38;5;203mtd[38;5;231m>[0m                               [34m│ [38;5;231m    <[38;5;203mtd[38;5;231m>Agent NGT</[38;5;203mtd[38;5;231m>[0m
[34m[38;5;231m    <[38;5;203mtd[38;5;231m>[0m                                             [34m│ [38;5;231m    <[38;5;203mtd[38;5;231m>[0m
[34m[48;5;52;38;5;231m      <[38;5;203mcode[38;5;231m>docker pull vdaas/vald-agent-ngt:v1.[48;5;124m7[48;5;52m.[48;5;124m0[48;5;52m<[34m↴[0m[34m│ [48;5;22;38;5;231m      <[38;5;203mcode[38;5;231m>docker pull vdaas/vald-agent-ngt:v1.[48;5;28m6[48;5;22m.[34m↴[0m[48;5;22m[0K[0m
[34m[0m[48;5;52m                                         [34m…[38;5;231m/[38;5;203mcode[38;5;231m><[38;5;203mbr[38;5;231m/>[0m[34m│ [0m[48;5;22m                                     [34m…[48;5;28;38;5;231m3[48;5;22m</[38;5;203mcode[38;5;231m><[38;5;203mbr[38;5;231m/>[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231m      <[38;5;203mcode[38;5;231m>docker pull ghcr.io/vdaas/vald/vald-agen[34m↴[0m[34m│ [48;5;22;38;5;231m      <[38;5;203mcode[38;5;231m>docker pull ghcr.io/vdaas/vald/vald-ag[34m↵[0m[48;5;22m[0K[0m
[34m[0m[48;5;52m                                 [34m…[38;5;231mt-ngt:v1.[48;5;124m7[48;5;52m.[48;5;124m0[48;5;52m</[38;5;203mcode[38;5;231m>[0m[34m│ [48;5;22;38;5;231ment-ngt:v1.[48;5;28m6[48;5;22m.[48;5;28m3[48;5;22m</[38;5;203mcode[38;5;231m>[0m[48;5;22m[0K[0m
[34m[38;5;231m    </[38;5;203mtd[38;5;231m>[0m                                            [34m│ [38;5;231m    </[38;5;203mtd[38;5;231m>[0m
[34m[38;5;231m  </[38;5;203mtr[38;5;231m>[0m                                              [34m│ [38;5;231m  </[38;5;203mtr[38;5;231m>[0m
[34m[38;5;231m  <[38;5;203mtr[38;5;231m>[0m                                               [34m│ [38;5;231m  <[38;5;203mtr[38;5;231m>[0m
[34m[38;5;231m    <[38;5;203mtd[38;5;231m>Agent sidecar</[38;5;203mtd[38;5;231m>[0m                           [34m│ [38;5;231m    <[38;5;203mtd[38;5;231m>Agent sidecar</[38;5;203mtd[38;5;231m>[0m
[34m[38;5;231m    <[38;5;203mtd[38;5;231m>[0m                                             [34m│ [38;5;231m    <[38;5;203mtd[38;5;231m>[0m
[34m[48;5;52;38;5;231m      <[38;5;203mcode[38;5;231m>docker pull vdaas/vald-agent-sidecar:v1.[34m↴[0m[34m│ [48;5;22;38;5;231m      <[38;5;203mcode[38;5;231m>docker pull vdaas/vald-agent-sidecar:v[34m↴[0m[48;5;22m[0K[0m
[34m[0m[48;5;52m                                     [34m…[48;5;124;38;5;231m7[48;5;52m.[48;5;124m0[48;5;52m</[38;5;203mcode[38;5;231m><[38;5;203mbr[38;5;231m/>[0m[34m│ [0m[48;5;22m                                 [34m…[38;5;231m1.[48;5;28m6[48;5;22m.[48;5;28m3[48;5;22m</[38;5;203mcode[38;5;231m><[38;5;203mbr[38;5;231m/>[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231m      <[38;5;203mcode[38;5;231m>docker pull ghcr.io/vdaas/vald/vald-agen[34m↵[0m[34m│ [48;5;22;38;5;231m      <[38;5;203mcode[38;5;231m>docker pull ghcr.io/vdaas/vald/vald-ag[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mt-sidecar:v1.[48;5;124m7[48;5;52m.[48;5;124m0[48;5;52m</[38;5;203mcode[38;5;231m>[0m[48;5;52m                              [0m[34m│ [48;5;22;38;5;231ment-sidecar:v1.[48;5;28m6[48;5;22m.[48;5;28m3[48;5;22m</[38;5;203mcode[38;5;231m>[0m[48;5;22m[0K[0m
[34m[38;5;231m    </[38;5;203mtd[38;5;231m>[0m                                            [34m│ [38;5;231m    </[38;5;203mtd[38;5;231m>[0m
[34m[38;5;231m  </[38;5;203mtr[38;5;231m>[0m                                              [34m│ [38;5;231m  </[38;5;203mtr[38;5;231m>[0m
[34m[38;5;231m  <[38;5;203mtr[38;5;231m>[0m                                               [34m│ [38;5;231m  <[38;5;203mtr[38;5;231m>[0m
[34m[38;5;231m    <[38;5;203mtd[38;5;231m>Discoverers</[38;5;203mtd[38;5;231m>[0m                             [34m│ [38;5;231m    <[38;5;203mtd[38;5;231m>Discoverers</[38;5;203mtd[38;5;231m>[0m
[34m[38;5;231m    <[38;5;203mtd[38;5;231m>[0m                                             [34m│ [38;5;231m    <[38;5;203mtd[38;5;231m>[0m
[34m[48;5;52;38;5;231m      <[38;5;203mcode[38;5;231m>docker pull vdaas/vald-discoverer-k8s:v1[34m↴[0m[34m│ [48;5;22;38;5;231m      <[38;5;203mcode[38;5;231m>docker pull vdaas/vald-discoverer-k8s:[34m↴[0m[48;5;22m[0K[0m
[34m[0m[48;5;52m                                    [34m…[38;5;231m.[48;5;124m7[48;5;52m.[48;5;124m0[48;5;52m</[38;5;203mcode[38;5;231m><[38;5;203mbr[38;5;231m/>[0m[34m│ [0m[48;5;22m                                [34m…[38;5;231mv1.[48;5;28m6[48;5;22m.[48;5;28m3[48;5;22m</[38;5;203mcode[38;5;231m><[38;5;203mbr[38;5;231m/>[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231m      <[38;5;203mcode[38;5;231m>docker pull ghcr.io/vdaas/vald/vald-disc[34m↵[0m[34m│ [48;5;22;38;5;231m      <[38;5;203mcode[38;5;231m>docker pull ghcr.io/vdaas/vald/vald-di[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231moverer-k8s:v1.[48;5;124m7[48;5;52m.[48;5;124m0[48;5;52m</[38;5;203mcode[38;5;231m>[0m[48;5;52m                             [0m[34m│ [48;5;22;38;5;231mscoverer-k8s:v1.[48;5;28m6[48;5;22m.[48;5;28m3[48;5;22m</[38;5;203mcode[38;5;231m>[0m[48;5;22m[0K[0m
[34m[38;5;231m    </[38;5;203mtd[38;5;231m>[0m                                            [34m│ [38;5;231m    </[38;5;203mtd[38;5;231m>[0m
[34m[38;5;231m  </[38;5;203mtr[38;5;231m>[0m                                              [34m│ [38;5;231m  </[38;5;203mtr[38;5;231m>[0m
[34m[38;5;231m  <[38;5;203mtr[38;5;231m>[0m                                               [34m│ [38;5;231m  <[38;5;203mtr[38;5;231m>[0m
[34m[38;5;231m    <[38;5;203mtd[38;5;231m>Gateways</[38;5;203mtd[38;5;231m>[0m                                [34m│ [38;5;231m    <[38;5;203mtd[38;5;231m>Gateways</[38;5;203mtd[38;5;231m>[0m
[34m[38;5;231m    <[38;5;203mtd[38;5;231m>[0m                                             [34m│ [38;5;231m    <[38;5;203mtd[38;5;231m>[0m
[34m[48;5;52;38;5;231m      <[38;5;203mcode[38;5;231m>docker pull vdaas/vald-lb-gateway:v1.[48;5;124m7[48;5;52m.[48;5;124m0[48;5;52;34m↴[0m[34m│ [48;5;22;38;5;231m      <[38;5;203mcode[38;5;231m>docker pull vdaas/vald-lb-gateway:v1.[48;5;28m6[48;5;22;34m↴[0m[48;5;22m[0K[0m
[34m[0m[48;5;52m                                        [34m…[38;5;231m</[38;5;203mcode[38;5;231m><[38;5;203mbr[38;5;231m/>[0m[34m│ [0m[48;5;22m                                    [34m…[38;5;231m.[48;5;28m3[48;5;22m</[38;5;203mcode[38;5;231m><[38;5;203mbr[38;5;231m/>[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231m      <[38;5;203mcode[38;5;231m>docker pull ghcr.io/vdaas/vald/vald-lb-g[34m↵[0m[34m│ [48;5;22;38;5;231m      <[38;5;203mcode[38;5;231m>docker pull ghcr.io/vdaas/vald/vald-lb[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mateway:v1.[48;5;124m7[48;5;52m.[48;5;124m0[48;5;52m</[38;5;203mcode[38;5;231m><[38;5;203mbr[38;5;231m/>[0m[48;5;52m                            [0m[34m│ [48;5;22;38;5;231m-gateway:v1.[48;5;28m6[48;5;22m.[48;5;28m3[48;5;22m</[38;5;203mcode[38;5;231m><[38;5;203mbr[38;5;231m/>[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231m      <[38;5;203mcode[38;5;231m>docker pull vdaas/vald-filter-gateway:v1[34m↴[0m[34m│ [48;5;22;38;5;231m      <[38;5;203mcode[38;5;231m>docker pull vdaas/vald-filter-gateway:[34m↴[0m[48;5;22m[0K[0m
[34m[0m[48;5;52m                                    [34m…[38;5;231m.[48;5;124m7[48;5;52m.[48;5;124m0[48;5;52m</[38;5;203mcode[38;5;231m><[38;5;203mbr[38;5;231m/>[0m[34m│ [0m[48;5;22m                                [34m…[38;5;231mv1.[48;5;28m6[48;5;22m.[48;5;28m3[48;5;22m</[38;5;203mcode[38;5;231m><[38;5;203mbr[38;5;231m/>[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231m      <[38;5;203mcode[38;5;231m>docker pull ghcr.io/vdaas/vald/vald-filt[34m↵[0m[34m│ [48;5;22;38;5;231m      <[38;5;203mcode[38;5;231m>docker pull ghcr.io/vdaas/vald/vald-fi[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mer-gateway:v1.[48;5;124m7[48;5;52m.[48;5;124m0[48;5;52m</[38;5;203mcode[38;5;231m>[0m[48;5;52m                             [0m[34m│ [48;5;22;38;5;231mlter-gateway:v1.[48;5;28m6[48;5;22m.[48;5;28m3[48;5;22m</[38;5;203mcode[38;5;231m>[0m[48;5;22m[0K[0m
[34m[38;5;231m    </[38;5;203mtd[38;5;231m>[0m                                            [34m│ [38;5;231m    </[38;5;203mtd[38;5;231m>[0m
[34m[38;5;231m  </[38;5;203mtr[38;5;231m>[0m                                              [34m│ [38;5;231m  </[38;5;203mtr[38;5;231m>[0m
[34m[38;5;231m  <[38;5;203mtr[38;5;231m>[0m                                               [34m│ [38;5;231m  <[38;5;203mtr[38;5;231m>[0m
[34m[38;5;231m    <[38;5;203mtd[38;5;231m>Index Manager</[38;5;203mtd[38;5;231m>[0m                           [34m│ [38;5;231m    <[38;5;203mtd[38;5;231m>Index Manager</[38;5;203mtd[38;5;231m>[0m
[34m[38;5;231m    <[38;5;203mtd[38;5;231m>[0m                                             [34m│ [38;5;231m    <[38;5;203mtd[38;5;231m>[0m
[34m[48;5;52;38;5;231m      <[38;5;203mcode[38;5;231m>docker pull vdaas/vald-manager-index:v1.[34m↴[0m[34m│ [48;5;22;38;5;231m      <[38;5;203mcode[38;5;231m>docker pull vdaas/vald-manager-index:v[34m↴[0m[48;5;22m[0K[0m
[34m[0m[48;5;52m                                     [34m…[48;5;124;38;5;231m7[48;5;52m.[48;5;124m0[48;5;52m</[38;5;203mcode[38;5;231m><[38;5;203mbr[38;5;231m/>[0m[34m│ [0m[48;5;22m                                 [34m…[38;5;231m1.[48;5;28m6[48;5;22m.[48;5;28m3[48;5;22m</[38;5;203mcode[38;5;231m><[38;5;203mbr[38;5;231m/>[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231m      <[38;5;203mcode[38;5;231m>docker pull ghcr.io/vdaas/vald/vald-mana[34m↵[0m[34m│ [48;5;22;38;5;231m      <[38;5;203mcode[38;5;231m>docker pull ghcr.io/vdaas/vald/vald-ma[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mger-index:v1.[48;5;124m7[48;5;52m.[48;5;124m0[48;5;52m</[38;5;203mcode[38;5;231m>[0m[48;5;52m                              [0m[34m│ [48;5;22;38;5;231mnager-index:v1.[48;5;28m6[48;5;22m.[48;5;28m3[48;5;22m</[38;5;203mcode[38;5;231m>[0m[48;5;22m[0K[0m
[34m[38;5;231m    </[38;5;203mtd[38;5;231m>[0m                                            [34m│ [38;5;231m    </[38;5;203mtd[38;5;231m>[0m
[34m[38;5;231m  </[38;5;203mtr[38;5;231m>[0m                                              [34m│ [38;5;231m  </[38;5;203mtr[38;5;231m>[0m
[34m[38;5;231m  <[38;5;203mtr[38;5;231m>[0m                                               [34m│ [38;5;231m  <[38;5;203mtr[38;5;231m>[0m
[34m[38;5;231m    <[38;5;203mtd[38;5;231m>Helm Operator</[38;5;203mtd[38;5;231m>[0m                           [34m│ [38;5;231m    <[38;5;203mtd[38;5;231m>Helm Operator</[38;5;203mtd[38;5;231m>[0m
[34m[38;5;231m    <[38;5;203mtd[38;5;231m>[0m                                             [34m│ [38;5;231m    <[38;5;203mtd[38;5;231m>[0m
[34m[48;5;52;38;5;231m      <[38;5;203mcode[38;5;231m>docker pull vdaas/vald-helm-operator:v1.[34m↴[0m[34m│ [48;5;22;38;5;231m      <[38;5;203mcode[38;5;231m>docker pull vdaas/vald-helm-operator:v[34m↴[0m[48;5;22m[0K[0m
[34m[0m[48;5;52m                                     [34m…[48;5;124;38;5;231m7[48;5;52m.[48;5;124m0[48;5;52m</[38;5;203mcode[38;5;231m><[38;5;203mbr[38;5;231m/>[0m[34m│ [0m[48;5;22m                                 [34m…[38;5;231m1.[48;5;28m6[48;5;22m.[48;5;28m3[48;5;22m</[38;5;203mcode[38;5;231m><[38;5;203mbr[38;5;231m/>[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231m      <[38;5;203mcode[38;5;231m>docker pull ghcr.io/vdaas/vald/vald-helm[34m↵[0m[34m│ [48;5;22;38;5;231m      <[38;5;203mcode[38;5;231m>docker pull ghcr.io/vdaas/vald/vald-he[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231m-operator:v1.[48;5;124m7[48;5;52m.[48;5;124m0[48;5;52m</[38;5;203mcode[38;5;231m>[0m[48;5;52m                              [0m[34m│ [48;5;22;38;5;231mlm-operator:v1.[48;5;28m6[48;5;22m.[48;5;28m3[48;5;22m</[38;5;203mcode[38;5;231m>[0m[48;5;22m[0K[0m
[34m[38;5;231m    </[38;5;203mtd[38;5;231m>[0m                                            [34m│ [38;5;231m    </[38;5;203mtd[38;5;231m>[0m
[34m[38;5;231m  </[38;5;203mtr[38;5;231m>[0m                                              [34m│ [38;5;231m  </[38;5;203mtr[38;5;231m>[0m
[34m[38;5;231m</[38;5;203mtable[38;5;231m>[0m                                             [34m│ [38;5;231m</[38;5;203mtable[38;5;231m>[0m
[34m[0m                                                     [34m│ [0m
[34m[38;5;208m### Documents[0m                                        [34m│ [38;5;208m### Documents[0m
[34m[0m                                                     [34m│ [0m
[34m[48;5;52;38;5;243m-[38;5;231m [GoDoc]([38;5;149mhttps://pkg.go.dev/github.com/vdaas/vald@v[34m↴[0m[34m│ [48;5;22;38;5;243m-[38;5;231m [GoDoc]([38;5;149mhttps://pkg.go.dev/github.com/vdaas/vald[34m↴[0m[48;5;22m[0K[0m
[34m[0m[48;5;52m                                              [34m…[38;5;149m1.[48;5;124m7[48;5;52m.[48;5;124m0[48;5;52;38;5;231m)[0m[34m│ [0m[48;5;22m                                          [34m…[38;5;149m@v1.[48;5;28m6[48;5;22m.[48;5;28m3[48;5;22;38;5;231m)[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;243m-[38;5;231m [Helm Chart Reference]([38;5;149mhttps://github.com/vdaas/va[34m↵[0m[34m│ [48;5;22;38;5;243m-[38;5;231m [Helm Chart Reference]([38;5;149mhttps://github.com/vdaas/[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;149mld/blob/v1.[48;5;124m7[48;5;52m.[48;5;124m0[48;5;52m/charts/vald/README.md[38;5;231m)[0m[48;5;52m                [0m[34m│ [48;5;22;38;5;149mvald/blob/v1.[48;5;28m6[48;5;22m.[48;5;28m3[48;5;22m/charts/vald/README.md[38;5;231m)[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;243m-[38;5;231m [Helm Operator Chart Reference]([38;5;149mhttps://github.com[34m↵[0m[34m│ [48;5;22;38;5;243m-[38;5;231m [Helm Operator Chart Reference]([38;5;149mhttps://github.c[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;149m/vdaas/vald/blob/v1.[48;5;124m7[48;5;52m.[48;5;124m0[48;5;52m/charts/vald-helm-operator/RE[34m↵[0m[34m│ [48;5;22;38;5;149mom/vdaas/vald/blob/v1.[48;5;28m6[48;5;22m.[48;5;28m3[48;5;22m/charts/vald-helm-operato[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;149mADME.md[38;5;231m)[0m[48;5;52m                                             [0m[34m│ [48;5;22;38;5;149mr/README.md[38;5;231m)[0m[48;5;22m[0K[0m
[34m[0m                                                     [34m│ [0m
[34m[38;5;208m### Changes[0m                                          [34m│ [38;5;208m### Changes[0m
[34m[0m                                                     [34m│ [0m

[1;4;33mcharts/vald-helm-operator/Chart.yaml[0m

[36m────[0m[36m┐[0m
[34m16[0m: [36m│[0m
[36m────[0m[36m┴[0m[36m─────────────────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[0m                                                     [34m│ [0m
[34m[38;5;203mapiVersion[38;5;231m: [38;5;186mv2[0m                                       [34m│ [38;5;203mapiVersion[38;5;231m: [38;5;186mv2[0m
[34m[38;5;203mname[38;5;231m: [38;5;186mvald-helm-operator[0m                             [34m│ [38;5;203mname[38;5;231m: [38;5;186mvald-helm-operator[0m
[34m[48;5;52;38;5;203mversion[38;5;231m: [38;5;186mv1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                                      [0m[34m│ [48;5;22;38;5;203mversion[38;5;231m: [38;5;186mv1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;203mdescription[38;5;231m: [38;5;186mA Helm chart for vald-helm-operator[0m     [34m│ [38;5;203mdescription[38;5;231m: [38;5;186mA Helm chart for vald-helm-operator[0m
[34m[38;5;203mtype[38;5;231m: [38;5;186mapplication[0m                                    [34m│ [38;5;203mtype[38;5;231m: [38;5;186mapplication[0m
[34m[38;5;203mkeywords[38;5;231m:[0m                                            [34m│ [38;5;203mkeywords[38;5;231m:[0m

[1;4;33mcharts/vald-helm-operator/README.md[0m

[36m───[0m[36m┐[0m
[34m2[0m: [36m│[0m
[36m───[0m[36m┴[0m[36m──────────────────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[0m                                                     [34m│ [0m
[34m[38;5;231mThis is a Helm chart to install vald-helm-operator.[0m  [34m│ [38;5;231mThis is a Helm chart to install vald-helm-operator.[0m
[34m[0m                                                     [34m│ [0m
[34m[48;5;52;38;5;231mCurrent chart version is [38;5;203m`v1.[48;5;124m7[48;5;52m.[48;5;124m0[48;5;52m`[0m[48;5;52m                    [0m[34m│ [48;5;22;38;5;231mCurrent chart version is [38;5;203m`v1.[48;5;28m6[48;5;22m.[48;5;28m3[48;5;22m`[0m[48;5;22m[0K[0m
[34m[0m                                                     [34m│ [0m
[34m[38;5;208m## Table of Contents[0m                                 [34m│ [38;5;208m## Table of Contents[0m
[34m[0m                                                     [34m│ [0m

[36m────────────────────────────────────────────────────[0m[36m┐[0m
[34m26[0m:[38;5;231m Run the following command to install the chart, [0m[36m│[0m
[36m────────────────────────────────────────────────────[0m[36m┴[0m[36m─────────────────────────────────────────────────────[0m
[34m[0m                                                     [34m│ [0m
[34m[38;5;231mPlease upgrade the CRDs first because Helm doesn't[34m↵[0m  [34m│ [38;5;231mPlease upgrade the CRDs first because Helm doesn't[34m↵[0m
[34m[38;5;231m have a support to upgrade CRDs.[0m                     [34m│ [38;5;231m have a support to upgrade CRDs.[0m
[34m[0m                                                     [34m│ [0m
[34m[48;5;52;38;5;141m    $ kubectl replace -f https://raw.githubuserconte[34m↵[0m[34m│ [48;5;22;38;5;141m    $ kubectl replace -f https://raw.githubusercon[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;141mnt.com/vdaas/vald/v1.[48;5;124m7[48;5;52m.[48;5;124m0[48;5;52m/charts/vald-helm-operator/c[34m↵[0m[34m│ [48;5;22;38;5;141mtent.com/vdaas/vald/v1.[48;5;28m6[48;5;22m.[48;5;28m3[48;5;22m/charts/vald-helm-operat[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;141mrds/valdrelease.yaml[0m[48;5;52m                                 [0m[34m│ [48;5;22;38;5;141mor/crds/valdrelease.yaml[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;141m    $ kubectl replace -f https://raw.githubuserconte[34m↵[0m[34m│ [48;5;22;38;5;141m    $ kubectl replace -f https://raw.githubusercon[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;141mnt.com/vdaas/vald/v1.[48;5;124m7[48;5;52m.[48;5;124m0[48;5;52m/charts/vald-helm-operator/c[34m↵[0m[34m│ [48;5;22;38;5;141mtent.com/vdaas/vald/v1.[48;5;28m6[48;5;22m.[48;5;28m3[48;5;22m/charts/vald-helm-operat[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;141mrds/valdhelmoperatorrelease.yaml[0m[48;5;52m                     [0m[34m│ [48;5;22;38;5;141mor/crds/valdhelmoperatorrelease.yaml[0m[48;5;22m[0K[0m
[34m[0m                                                     [34m│ [0m
[34m[38;5;231mAfter upgrading CRDs, you can upgrade the operator.[0m  [34m│ [38;5;231mAfter upgrading CRDs, you can upgrade the operator.[0m
[34m[38;5;231mIf you're using [38;5;203m`valdhelmoperatorrelease`[38;5;231m (or [38;5;203m`vho[34m↵[0m  [34m│ [38;5;231mIf you're using [38;5;203m`valdhelmoperatorrelease`[38;5;231m (or [38;5;203m`vho[34m↵[0m
[34m[38;5;203mr`[38;5;231m) resource, please update the [38;5;203m`spec.image.tag`[38;5;231m f[34m↵[0m  [34m│ [38;5;203mr`[38;5;231m) resource, please update the [38;5;203m`spec.image.tag`[38;5;231m f[34m↵[0m
[34m[38;5;231mield of it.[0m                                          [34m│ [38;5;231mield of it.[0m
[34m[0m                                                     [34m│ [0m
[34m[48;5;52;38;5;141m    $ kubectl patch vhor vhor-release -p '{"spec":{"[34m↵[0m[34m│ [48;5;22;38;5;141m    $ kubectl patch vhor vhor-release -p '{"spec":[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;141mimage":{"tag":"v1.[48;5;124m7[48;5;52m.[48;5;124m0[48;5;52m"}}}'[0m[48;5;52m                           [0m[34m│ [48;5;22;38;5;141m{"image":{"tag":"v1.[48;5;28m6[48;5;22m.[48;5;28m3[48;5;22m"}}}'[0m[48;5;22m[0K[0m
[34m[0m                                                     [34m│ [0m
[34m[38;5;231mOn the other hand, please update the operator's de[34m↴[0m  [34m│ [38;5;231mOn the other hand, please update the operator's de[34m↴[0m
[34m[0m                                [34m…[38;5;231mployment manually.[0m  [34m│ [0m                                [34m…[38;5;231mployment manually.[0m
[34m[0m                                                     [34m│ [0m

[36m─────────────[0m[36m┐[0m
[34m79[0m:[38;5;231m spec: {} [0m[36m│[0m
[36m─────────────[0m[36m┴[0m[36m────────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[38;5;231m| healthPort                         | int    | [38;5;203m`8[34m↵[0m  [34m│ [38;5;231m| healthPort                         | int    | [38;5;203m`8[34m↵[0m
[34m[38;5;203m081`[38;5;231m                                              [34m↵[0m  [34m│ [38;5;203m081`[38;5;231m                                              [34m↵[0m
[34m[38;5;231m                                                    [0m[7m→[0m[34m│ [38;5;231m                                                  [0m[7m→[0m
[34m[38;5;231m| image.pullPolicy                   | string | `"[34m↵[0m  [34m│ [38;5;231m| image.pullPolicy                   | string | `"[34m↵[0m
[34m[38;5;231mAlways"`                                          [34m↵[0m  [34m│ [38;5;231mAlways"`                                          [34m↵[0m
[34m[38;5;231m                                                    [0m[7m→[0m[34m│ [38;5;231m                                                  [0m[7m→[0m
[34m[38;5;231m| image.repository                   | string | [38;5;203m`"[34m↵[0m  [34m│ [38;5;231m| image.repository                   | string | [38;5;203m`"[34m↵[0m
[34m[38;5;203mvdaas/vald-helm-operator"`[38;5;231m                        [34m↵[0m  [34m│ [38;5;203mvdaas/vald-helm-operator"`[38;5;231m                        [34m↵[0m
[34m[38;5;231m                                                    [0m[7m→[0m[34m│ [38;5;231m                                                  [0m[7m→[0m
[34m[48;5;52;38;5;231m| image.tag                          | string | [38;5;203m`"v1[34m↵[0m[34m│ [48;5;22;38;5;231m| image.tag                          | string | [38;5;203m`"[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;203m.[48;5;124m7[48;5;52m.[48;5;124m0[48;5;52m"`[38;5;231m                                              [34m↵[0m[34m│ [48;5;22;38;5;203mv1.[48;5;28m6[48;5;22m.[48;5;28m3[48;5;22m"`[38;5;231m                                          [34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231m                                                    [0m[7m→[0m[34m│ [48;5;22;38;5;231m                                                  [0m[7m→[0m[48;5;22m[0K[0m
[34m[38;5;231m| leaderElectionID                   | string | [38;5;203m`"[34m↵[0m  [34m│ [38;5;231m| leaderElectionID                   | string | [38;5;203m`"[34m↵[0m
[34m[38;5;203mvald-helm-operator"`[38;5;231m                              [34m↵[0m  [34m│ [38;5;203mvald-helm-operator"`[38;5;231m                              [34m↵[0m
[34m[38;5;231m                                                    [0m[7m→[0m[34m│ [38;5;231m                                                  [0m[7m→[0m
[34m[38;5;231m| livenessProbe.enabled              | bool   | [38;5;203m`t[34m↵[0m  [34m│ [38;5;231m| livenessProbe.enabled              | bool   | [38;5;203m`t[34m↵[0m
[34m[38;5;203mrue`[38;5;231m                                              [34m↵[0m  [34m│ [38;5;203mrue`[38;5;231m                                              [34m↵[0m
[34m[38;5;231m                                                    [0m[7m→[0m[34m│ [38;5;231m                                                  [0m[7m→[0m
[34m[38;5;231m| livenessProbe.failureThreshold     | int    | [38;5;203m`2[34m↵[0m  [34m│ [38;5;231m| livenessProbe.failureThreshold     | int    | [38;5;203m`2[34m↵[0m
[34m[38;5;203m`[38;5;231m                                                 [34m↵[0m  [34m│ [38;5;203m`[38;5;231m                                                 [34m↵[0m
[34m[38;5;231m                                                    [0m[7m→[0m[34m│ [38;5;231m                                                  [0m[7m→[0m

[1;4;33mcharts/vald-helm-operator/values.yaml[0m

[36m───────────[0m[36m┐[0m
[34m29[0m:[38;5;231m [38;5;203mimage[38;5;231m: [0m[36m│[0m
[36m───────────[0m[36m┴[0m[36m──────────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[38;5;231m  [38;5;203mrepository[38;5;231m: [38;5;186mvdaas/vald-helm-operator[0m               [34m│ [38;5;231m  [38;5;203mrepository[38;5;231m: [38;5;186mvdaas/vald-helm-operator[0m
[34m[38;5;231m  [38;5;242m# @schema {"name": "image.tag", "type": "string"}[0m  [34m│ [38;5;231m  [38;5;242m# @schema {"name": "image.tag", "type": "string"}[0m
[34m[38;5;231m  [38;5;242m# image.tag -- image tag[0m                           [34m│ [38;5;231m  [38;5;242m# image.tag -- image tag[0m
[34m[48;5;52;38;5;231m  [38;5;203mtag[38;5;231m: [38;5;186mv1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                                        [0m[34m│ [48;5;22;38;5;231m  [38;5;203mtag[38;5;231m: [38;5;186mv1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;231m  [38;5;242m# @schema {"name": "image.pullPolicy", "type": "[34m↵[0m  [34m│ [38;5;231m  [38;5;242m# @schema {"name": "image.pullPolicy", "type": "[34m↵[0m
[34m[38;5;242mstring", "enum": ["Always", "Never", "IfNotPresent[34m↵[0m  [34m│ [38;5;242mstring", "enum": ["Always", "Never", "IfNotPresent[34m↵[0m
[34m[38;5;242m"]}[0m                                                  [34m│ [38;5;242m"]}[0m
[34m[38;5;231m  [38;5;242m# image.pullPolicy -- image pull policy[0m            [34m│ [38;5;231m  [38;5;242m# image.pullPolicy -- image pull policy[0m
[34m[38;5;231m  [38;5;203mpullPolicy[38;5;231m: [38;5;186mAlways[0m                                 [34m│ [38;5;231m  [38;5;203mpullPolicy[38;5;231m: [38;5;186mAlways[0m

[1;4;33mcharts/vald/Chart.yaml[0m

[36m────[0m[36m┐[0m
[34m16[0m: [36m│[0m
[36m────[0m[36m┴[0m[36m─────────────────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[0m                                                     [34m│ [0m
[34m[38;5;203mapiVersion[38;5;231m: [38;5;186mv2[0m                                       [34m│ [38;5;203mapiVersion[38;5;231m: [38;5;186mv2[0m
[34m[38;5;203mname[38;5;231m: [38;5;186mvald[0m                                           [34m│ [38;5;203mname[38;5;231m: [38;5;186mvald[0m
[34m[48;5;52;38;5;203mversion[38;5;231m: [38;5;186mv1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                                      [0m[34m│ [48;5;22;38;5;203mversion[38;5;231m: [38;5;186mv1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;203mdescription[38;5;231m: [38;5;186mA distributed high scalable & high-sp[34m↵[0m  [34m│ [38;5;203mdescription[38;5;231m: [38;5;186mA distributed high scalable & high-sp[34m↵[0m
[34m[38;5;186meed approximate nearest neighbor search engine[0m       [34m│ [38;5;186meed approximate nearest neighbor search engine[0m
[34m[38;5;203mtype[38;5;231m: [38;5;186mapplication[0m                                    [34m│ [38;5;203mtype[38;5;231m: [38;5;186mapplication[0m
[34m[38;5;203mkeywords[38;5;231m:[0m                                            [34m│ [38;5;203mkeywords[38;5;231m:[0m

[1;4;33mcharts/vald/README.md[0m

[36m───[0m[36m┐[0m
[34m2[0m: [36m│[0m
[36m───[0m[36m┴[0m[36m──────────────────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[0m                                                     [34m│ [0m
[34m[38;5;231mThis is a Helm chart to install Vald components.[0m     [34m│ [38;5;231mThis is a Helm chart to install Vald components.[0m
[34m[0m                                                     [34m│ [0m
[34m[48;5;52;38;5;231mCurrent chart version is [38;5;203m`v1.[48;5;124m7[48;5;52m.[48;5;124m0[48;5;52m`[0m[48;5;52m                    [0m[34m│ [48;5;22;38;5;231mCurrent chart version is [38;5;203m`v1.[48;5;28m6[48;5;22m.[48;5;28m3[48;5;22m`[0m[48;5;22m[0K[0m
[34m[0m                                                     [34m│ [0m
[34m[38;5;208m## Table of Contents[0m                                 [34m│ [38;5;208m## Table of Contents[0m
[34m[0m                                                     [34m│ [0m

[36m─────────────────────────────────────────────────────[0m[36m┐[0m
[34m290[0m:[38;5;231m Run the following command to install the chart, [0m[36m│[0m
[36m─────────────────────────────────────────────────────[0m[36m┴[0m[36m────────────────────────────────────────────────────[0m
[34m[38;5;231m| defaults.grpc.client.tls.enabled                [34m↵[0m  [34m│ [38;5;231m| defaults.grpc.client.tls.enabled                [34m↵[0m
[34m[38;5;231m                                                  [34m↵[0m  [34m│ [38;5;231m                                                  [34m↵[0m
[34m[38;5;231m    | bool   | [38;5;203m`false`[38;5;231m                              [0m[7m→[0m[34m│ [38;5;231m    | bool   | [38;5;203m`false`[38;5;231m                            [0m[7m→[0m
[34m[38;5;231m| defaults.grpc.client.tls.insecure_skip_verify   [34m↵[0m  [34m│ [38;5;231m| defaults.grpc.client.tls.insecure_skip_verify   [34m↵[0m
[34m[38;5;231m                                                  [34m↵[0m  [34m│ [38;5;231m                                                  [34m↵[0m
[34m[38;5;231m    | bool   | `false`                              [0m[7m→[0m[34m│ [38;5;231m    | bool   | `false`                            [0m[7m→[0m
[34m[38;5;231m| defaults.grpc.client.tls.key                    [34m↵[0m  [34m│ [38;5;231m| defaults.grpc.client.tls.key                    [34m↵[0m
[34m[38;5;231m                                                  [34m↵[0m  [34m│ [38;5;231m                                                  [34m↵[0m
[34m[38;5;231m    | string | [38;5;203m`"/path/to/key"`[38;5;231m                     [0m[7m→[0m[34m│ [38;5;231m    | string | [38;5;203m`"/path/to/key"`[38;5;231m                   [0m[7m→[0m
[34m[48;5;52;38;5;231m| defaults.image.tag                                [34m↵[0m[34m│ [48;5;22;38;5;231m| defaults.image.tag                              [34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231m                                                    [34m↵[0m[34m│ [48;5;22;38;5;231m                                                  [34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231m| string | [38;5;203m`"v1.[48;5;124m7[48;5;52m.[48;5;124m0[48;5;52m"`[38;5;231m                               [0m[7m→[0m[34m│ [48;5;22;38;5;231m    | string | [38;5;203m`"v1.[48;5;28m6[48;5;22m.[48;5;28m3[48;5;22m"`[38;5;231m                         [0m[7m→[0m[48;5;22m[0K[0m
[34m[38;5;231m| defaults.logging.format                         [34m↵[0m  [34m│ [38;5;231m| defaults.logging.format                         [34m↵[0m
[34m[38;5;231m                                                  [34m↵[0m  [34m│ [38;5;231m                                                  [34m↵[0m
[34m[38;5;231m    | string | [38;5;203m`"raw"`[38;5;231m                              [0m[7m→[0m[34m│ [38;5;231m    | string | [38;5;203m`"raw"`[38;5;231m                            [0m[7m→[0m
[34m[38;5;231m| defaults.logging.level                          [34m↵[0m  [34m│ [38;5;231m| defaults.logging.level                          [34m↵[0m
[34m[38;5;231m                                                  [34m↵[0m  [34m│ [38;5;231m                                                  [34m↵[0m
[34m[38;5;231m    | string | [38;5;203m`"debug"`[38;5;231m                            [0m[7m→[0m[34m│ [38;5;231m    | string | [38;5;203m`"debug"`[38;5;231m                          [0m[7m→[0m
[34m[38;5;231m| defaults.logging.logger                         [34m↵[0m  [34m│ [38;5;231m| defaults.logging.logger                         [34m↵[0m
[34m[38;5;231m                                                  [34m↵[0m  [34m│ [38;5;231m                                                  [34m↵[0m
[34m[38;5;231m    | string | [38;5;203m`"glg"`[38;5;231m                              [0m[7m→[0m[34m│ [38;5;231m    | string | [38;5;203m`"glg"`[38;5;231m                            [0m[7m→[0m

[1;4;33mcharts/vald/values.yaml[0m

[36m──────────────[0m[36m┐[0m
[34m38[0m:[38;5;231m [38;5;203mdefaults[38;5;231m: [0m[36m│[0m
[36m──────────────[0m[36m┴[0m[36m───────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[38;5;231m  [38;5;203mimage[38;5;231m:[0m                                             [34m│ [38;5;231m  [38;5;203mimage[38;5;231m:[0m
[34m[38;5;231m    [38;5;242m# @schema {"name": "defaults.image.tag", "type[34m↴[0m  [34m│ [38;5;231m    [38;5;242m# @schema {"name": "defaults.image.tag", "type[34m↴[0m
[34m[0m                                      [34m…[38;5;242m": "string"}[0m  [34m│ [0m                                      [34m…[38;5;242m": "string"}[0m
[34m[38;5;231m    [38;5;242m# defaults.image.tag -- docker image tag[0m         [34m│ [38;5;231m    [38;5;242m# defaults.image.tag -- docker image tag[0m
[34m[48;5;52;38;5;231m    [38;5;203mtag[38;5;231m: [38;5;186mv1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                                      [0m[34m│ [48;5;22;38;5;231m    [38;5;203mtag[38;5;231m: [38;5;186mv1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;231m  [38;5;242m# @schema {"name": "defaults.server_config", "ty[34m↵[0m  [34m│ [38;5;231m  [38;5;242m# @schema {"name": "defaults.server_config", "ty[34m↵[0m
[34m[38;5;242mpe": "object", "anchor": "server_config"}[0m            [34m│ [38;5;242mpe": "object", "anchor": "server_config"}[0m
[34m[38;5;231m  [38;5;203mserver_config[38;5;231m:[0m                                     [34m│ [38;5;231m  [38;5;203mserver_config[38;5;231m:[0m
[34m[38;5;231m    [38;5;242m# @schema {"name": "defaults.server_config.ser[34m↵[0m  [34m│ [38;5;231m    [38;5;242m# @schema {"name": "defaults.server_config.ser[34m↵[0m
[34m[38;5;242mvers", "type": "object"}[0m                             [34m│ [38;5;242mvers", "type": "object"}[0m

[1;4;33mexample/client/go.mod[0m

[36m──────────────[0m[36m┐[0m
[34m11[0m:[38;5;231m replace ( [0m[36m│[0m
[36m──────────────[0m[36m┴[0m[36m───────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[38;5;231m    golang.org/x/crypto => golang.org/x/crypto v0.[34m↴[0m  [34m│ [38;5;231m    golang.org/x/crypto => golang.org/x/crypto v0.[34m↴[0m
[34m[0m                                               [34m…[38;5;231m5.0[0m  [34m│ [0m                                               [34m…[38;5;231m5.0[0m
[34m[38;5;231m    golang.org/x/net => golang.org/x/net v0.5.0[0m      [34m│ [38;5;231m    golang.org/x/net => golang.org/x/net v0.5.0[0m
[34m[38;5;231m    golang.org/x/text => golang.org/x/text v0.6.0[0m    [34m│ [38;5;231m    golang.org/x/text => golang.org/x/text v0.6.0[0m
[34m[48;5;52;38;5;231m    google.golang.org/genproto => google.golang.org/[34m↵[0m[34m│ [48;5;22;38;5;231m    google.golang.org/genproto => google.golang.or[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mgenproto v0.0.0-[48;5;124m20230117162540[48;5;52m-[48;5;124m28d6b9783ac4[0m[48;5;52m          [0m[34m│ [48;5;22;38;5;231mg/genproto v0.0.0-[48;5;28m20230123190316[48;5;22m-[48;5;28m2c411cf9d197[0m[48;5;22m[0K[0m
[34m[38;5;231m    google.golang.org/grpc => google.golang.org/gr[34m↴[0m  [34m│ [38;5;231m    google.golang.org/grpc => google.golang.org/gr[34m↴[0m
[34m[0m                                        [34m…[38;5;231mpc v1.52.0[0m  [34m│ [0m                                        [34m…[38;5;231mpc v1.52.0[0m
[34m[38;5;231m    google.golang.org/protobuf => google.golang.or[34m↴[0m  [34m│ [38;5;231m    google.golang.org/protobuf => google.golang.or[34m↴[0m
[34m[0m                                [34m…[38;5;231mg/protobuf v1.28.1[0m  [34m│ [0m                                [34m…[38;5;231mg/protobuf v1.28.1[0m
[34m[38;5;231m    gopkg.in/yaml.v2 => gopkg.in/yaml.v2 v2.4.0[0m      [34m│ [38;5;231m    gopkg.in/yaml.v2 => gopkg.in/yaml.v2 v2.4.0[0m

[36m──────────────[0m[36m┐[0m
[34m21[0m:[38;5;231m replace ( [0m[36m│[0m
[36m──────────────[0m[36m┴[0m[36m───────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[38;5;231mrequire ([0m                                            [34m│ [38;5;231mrequire ([0m
[34m[38;5;231m    github.com/kpango/fuid v0.0.0-20221203053508-5[34m↴[0m  [34m│ [38;5;231m    github.com/kpango/fuid v0.0.0-20221203053508-5[34m↴[0m
[34m[0m                                       [34m…[38;5;231m03b5ad89aa1[0m  [34m│ [0m                                       [34m…[38;5;231m03b5ad89aa1[0m
[34m[38;5;231m    github.com/kpango/glg v1.6.14[0m                    [34m│ [38;5;231m    github.com/kpango/glg v1.6.14[0m
[34m[48;5;52;38;5;231m    github.com/vdaas/vald-client-go v1.[48;5;124m6[48;5;52m.[48;5;124m3[0m[48;5;52m           [0m[34m│ [48;5;22;38;5;231m    github.com/vdaas/vald-client-go v1.[48;5;28m7[48;5;22m.[48;5;28m0[0m[48;5;22m[0K[0m
[34m[38;5;231m    gonum.org/v1/hdf5 v0.0.0-20210714002203-8c5d23[34m↴[0m  [34m│ [38;5;231m    gonum.org/v1/hdf5 v0.0.0-20210714002203-8c5d23[34m↴[0m
[34m[0m                                            [34m…[38;5;231mbc6946[0m  [34m│ [0m                                            [34m…[38;5;231mbc6946[0m
[34m[38;5;231m    google.golang.org/grpc v1.51.0[0m                   [34m│ [38;5;231m    google.golang.org/grpc v1.51.0[0m
[34m[38;5;231m)[0m                                                    [34m│ [38;5;231m)[0m

[1;4;33mexample/client/go.sum[0m

[36m─────────────────────────────────────────────────────────────────────────────────────[0m[36m┐[0m
[34m13[0m:[38;5;231m github.com/kpango/fuid v0.0.0-20221203053508-503b5ad89aa1/go.mod h1:CAYeq6us9Nfn [0m[36m│[0m
[36m─────────────────────────────────────────────────────────────────────────────────────[0m[36m┴[0m[36m────────────────────[0m
[34m[38;5;231mgithub.com/kpango/glg v1.6.14 h1:Ss3ZvTQ23blUCDYiz[34m↵[0m  [34m│ [38;5;231mgithub.com/kpango/glg v1.6.14 h1:Ss3ZvTQ23blUCDYiz[34m↵[0m
[34m[38;5;231mSAijiFTZsgGeYr/lanUGgQ10rY=[0m                          [34m│ [38;5;231mSAijiFTZsgGeYr/lanUGgQ10rY=[0m
[34m[38;5;231mgithub.com/kpango/glg v1.6.14/go.mod h1:2djk7Zr4zK[34m↵[0m  [34m│ [38;5;231mgithub.com/kpango/glg v1.6.14/go.mod h1:2djk7Zr4zK[34m↵[0m
[34m[38;5;231mIYPHlORH8tJVlhCEh+XXW8W4K3qJyNXMI=[0m                   [34m│ [38;5;231mIYPHlORH8tJVlhCEh+XXW8W4K3qJyNXMI=[0m
[34m[38;5;231mgithub.com/sirupsen/logrus v1.9.0 h1:trlNQbNUG3OdD[34m↵[0m  [34m│ [38;5;231mgithub.com/sirupsen/logrus v1.9.0 h1:trlNQbNUG3OdD[34m↵[0m
[34m[38;5;231mrDil03MCb1H2o9nJ1x4/5LYw7byDE0=[0m                      [34m│ [38;5;231mrDil03MCb1H2o9nJ1x4/5LYw7byDE0=[0m
[34m[48;5;52;38;5;231mgithub.com/vdaas/vald-client-go v1.[48;5;124m6[48;5;52m.[48;5;124m3[48;5;52m h1:[48;5;124m93aY9jOWlt[48;5;52;34m↵[0m[34m│ [48;5;22;38;5;231mgithub.com/vdaas/vald-client-go v1.[48;5;28m7[48;5;22m.[48;5;28m0[48;5;22m h1:[48;5;28mmFkd0/E7[48;5;22;34m↵[0m[48;5;22m[0K[0m
[34m[48;5;124;38;5;231mrO2b8hnIrz5P3NcxWAHMbqY8AnhNpud7w[48;5;52m=[0m[48;5;52m                   [0m[34m│ [48;5;28;38;5;231mOHsAq6of04mUuYTzlbsW+j+daPjRfaau0SA[48;5;22m=[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mgithub.com/vdaas/vald-client-go v1.[48;5;124m6[48;5;52m.[48;5;124m3[48;5;52m/go.mod h1:[48;5;124mWiE[48;5;52;34m↵[0m[34m│ [48;5;22;38;5;231mgithub.com/vdaas/vald-client-go v1.[48;5;28m7[48;5;22m.[48;5;28m0[48;5;22m/go.mod h1:[48;5;28m1[48;5;22;34m↵[0m[48;5;22m[0K[0m
[34m[48;5;124;38;5;231m3uVM1gjAEi4wbQi3S7lwfASR4BMiUvdOsM34XGqw[48;5;52m=[0m[48;5;52m            [0m[34m│ [48;5;28;38;5;231m5SmvrXHrbBmKQYwG3n/uiNN3kj5r56m3kZ5630msPA[48;5;22m=[0m[48;5;22m[0K[0m
[34m[38;5;231mgo.uber.org/atomic v1.7.0 h1:ADUqmZGgLDDfbSL9ZmPxK[34m↵[0m  [34m│ [38;5;231mgo.uber.org/atomic v1.7.0 h1:ADUqmZGgLDDfbSL9ZmPxK[34m↵[0m
[34m[38;5;231mTybcoEYHgpYfELNoN+7hsw=[0m                              [34m│ [38;5;231mTybcoEYHgpYfELNoN+7hsw=[0m
[34m[38;5;231mgo.uber.org/multierr v1.6.0 h1:y6IPFStTAIT5Ytl7/XY[34m↵[0m  [34m│ [38;5;231mgo.uber.org/multierr v1.6.0 h1:y6IPFStTAIT5Ytl7/XY[34m↵[0m
[34m[38;5;231mmHvzXQ7S3g/IeZW9hyZ5thw4=[0m                            [34m│ [38;5;231mmHvzXQ7S3g/IeZW9hyZ5thw4=[0m
[34m[38;5;231mgo.uber.org/zap v1.23.0 h1:OjGQ5KQDEUawVHxNwQgPpiy[34m↵[0m  [34m│ [38;5;231mgo.uber.org/zap v1.23.0 h1:OjGQ5KQDEUawVHxNwQgPpiy[34m↵[0m
[34m[38;5;231mpGHOxo2mNZsOqTak4fFY=[0m                                [34m│ [38;5;231mpGHOxo2mNZsOqTak4fFY=[0m

[36m────────────────────────────────────────────────────────────────────────────────────[0m[36m┐[0m
[34m27[0m:[38;5;231m golang.org/x/text v0.6.0/go.mod h1:mrYo+phRRbMaCq/xk9113O4dZlRixOauAjOtrjsXDZ8= [0m[36m│[0m
[36m────────────────────────────────────────────────────────────────────────────────────[0m[36m┴[0m[36m─────────────────────[0m
[34m[38;5;231mgolang.org/x/xerrors v0.0.0-20191204190536-9bdfabe[34m↵[0m  [34m│ [38;5;231mgolang.org/x/xerrors v0.0.0-20191204190536-9bdfabe[34m↵[0m
[34m[38;5;231m68543/go.mod h1:I/5z698sn9Ka8TeJc9MKroUUfqBBauWjQq[34m↵[0m  [34m│ [38;5;231m68543/go.mod h1:I/5z698sn9Ka8TeJc9MKroUUfqBBauWjQq[34m↵[0m
[34m[38;5;231mLJ2OPfmY0=[0m                                           [34m│ [38;5;231mLJ2OPfmY0=[0m
[34m[38;5;231mgonum.org/v1/hdf5 v0.0.0-20210714002203-8c5d23bc69[34m↵[0m  [34m│ [38;5;231mgonum.org/v1/hdf5 v0.0.0-20210714002203-8c5d23bc69[34m↵[0m
[34m[38;5;231m46 h1:vJpL69PeUullhJyKtTjHjENEmZU3BkO4e+fod7nKzgM=[0m   [34m│ [38;5;231m46 h1:vJpL69PeUullhJyKtTjHjENEmZU3BkO4e+fod7nKzgM=[0m
[34m[38;5;231mgonum.org/v1/hdf5 v0.0.0-20210714002203-8c5d23bc69[34m↵[0m  [34m│ [38;5;231mgonum.org/v1/hdf5 v0.0.0-20210714002203-8c5d23bc69[34m↵[0m
[34m[38;5;231m46/go.mod h1:BQUWDHIAygjdt1HnUPQ0eWqLN2n5FwJycrpYU[34m↵[0m  [34m│ [38;5;231m46/go.mod h1:BQUWDHIAygjdt1HnUPQ0eWqLN2n5FwJycrpYU[34m↵[0m
[34m[38;5;231mVUOx2I=[0m                                              [34m│ [38;5;231mVUOx2I=[0m
[34m[48;5;52;38;5;231mgoogle.golang.org/genproto v0.0.0-20230117162540-28d[34m↵[0m[34m│ [0m
[34m[48;5;52;38;5;231m6b9783ac4 h1:yF0uHwqqYt2tIL2F4hxRWA1ZFX43SEunWAK8MnQ[34m↵[0m[34m│ [0m
[34m[48;5;52;38;5;231miclk=[0m[48;5;52m                                                [0m[34m│ [0m
[34m[0m                                                     [34m│ [48;5;22;38;5;231mgoogle.golang.org/genproto v0.0.0-20230123190316-2[34m↵[0m[48;5;22m[0K[0m
[34m[0m                                                     [34m│ [48;5;22;38;5;231mc411cf9d197 h1:BwjeHhu4HS48EZmu1nS7flldBIDPC3qn+Hq[34m↵[0m[48;5;22m[0K[0m
[34m[0m                                                     [34m│ [48;5;22;38;5;231maSQ1K4x8=[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mgoogle.golang.org/genproto v0.0.0-[48;5;124m20230117162540[48;5;52m-[48;5;124m28d[48;5;52;34m↵[0m[34m│ [48;5;22;38;5;231mgoogle.golang.org/genproto v0.0.0-[48;5;28m20230123190316[48;5;22m-[48;5;28m2[48;5;22;34m↵[0m[48;5;22m[0K[0m
[34m[48;5;124;38;5;231m6b9783ac4[48;5;52m/go.mod h1:RGgjbofJ8xD9Sq1VVhDM1Vok1vRONV+r[34m↵[0m[34m│ [48;5;28;38;5;231mc411cf9d197[48;5;22m/go.mod h1:RGgjbofJ8xD9Sq1VVhDM1Vok1vRO[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mg+CjzG4SZKM=[0m[48;5;52m                                         [0m[34m│ [48;5;22;38;5;231mNV+rg+CjzG4SZKM=[0m[48;5;22m[0K[0m
[34m[38;5;231mgoogle.golang.org/grpc v1.52.0 h1:kd48UiU7EHsV4rnL[34m↵[0m  [34m│ [38;5;231mgoogle.golang.org/grpc v1.52.0 h1:kd48UiU7EHsV4rnL[34m↵[0m
[34m[38;5;231myOJRuP/Il/UHE7gdDAQ+SZI7nZk=[0m                         [34m│ [38;5;231myOJRuP/Il/UHE7gdDAQ+SZI7nZk=[0m
[34m[38;5;231mgoogle.golang.org/grpc v1.52.0/go.mod h1:pu6fVzoFb[34m↵[0m  [34m│ [38;5;231mgoogle.golang.org/grpc v1.52.0/go.mod h1:pu6fVzoFb[34m↵[0m
[34m[38;5;231m+NBYNAvQL08ic+lvB2IojljRYuun5vorUY=[0m                  [34m│ [38;5;231m+NBYNAvQL08ic+lvB2IojljRYuun5vorUY=[0m
[34m[38;5;231mgoogle.golang.org/protobuf v1.28.1 h1:d0NfwRgPtno5[34m↵[0m  [34m│ [38;5;231mgoogle.golang.org/protobuf v1.28.1 h1:d0NfwRgPtno5[34m↵[0m
[34m[38;5;231mB1Wa6L2DAG+KivqkdutMf1UhdNx175w=[0m                     [34m│ [38;5;231mB1Wa6L2DAG+KivqkdutMf1UhdNx175w=[0m

[1;4;33mgo.mod[0m

[36m────────────────────────────────[0m[36m┐[0m
[34m3[0m:[38;5;231m module github.com/vdaas/vald [0m[36m│[0m
[36m────────────────────────────────[0m[36m┴[0m[36m─────────────────────────────────────────────────────────────────────────[0m
[34m[38;5;231mgo 1.19[0m                                              [34m│ [38;5;231mgo 1.19[0m
[34m[0m                                                     [34m│ [0m
[34m[38;5;231mreplace ([0m                                            [34m│ [38;5;231mreplace ([0m
[34m[48;5;52;38;5;231m    cloud.google.com/go => cloud.google.com/go v0.[48;5;124m10[48;5;52;34m↴[0m[34m│ [48;5;22;38;5;231m    cloud.google.com/go => cloud.google.com/go v0.[34m↴[0m[48;5;22m[0K[0m
[34m[0m[48;5;52m                                                 [34m…[48;5;124;38;5;231m8[48;5;52m.0[0m[34m│ [0m[48;5;22m                                             [34m…[48;5;28;38;5;231m109[48;5;22m.0[0m[48;5;22m[0K[0m
[34m[38;5;231m    cloud.google.com/go/bigquery => cloud.google.c[34m↵[0m  [34m│ [38;5;231m    cloud.google.com/go/bigquery => cloud.google.c[34m↵[0m
[34m[38;5;231mom/go/bigquery v1.45.0[0m                               [34m│ [38;5;231mom/go/bigquery v1.45.0[0m
[34m[38;5;231m    cloud.google.com/go/compute => cloud.google.co[34m↵[0m  [34m│ [38;5;231m    cloud.google.com/go/compute => cloud.google.co[34m↵[0m
[34m[38;5;231mm/go/compute v1.15.1[0m                                 [34m│ [38;5;231mm/go/compute v1.15.1[0m
[34m[38;5;231m    cloud.google.com/go/datastore => cloud.google.[34m↵[0m  [34m│ [38;5;231m    cloud.google.com/go/datastore => cloud.google.[34m↵[0m
[34m[38;5;231mcom/go/datastore v1.10.0[0m                             [34m│ [38;5;231mcom/go/datastore v1.10.0[0m
[34m[38;5;231m    cloud.google.com/go/firestore => cloud.google.[34m↵[0m  [34m│ [38;5;231m    cloud.google.com/go/firestore => cloud.google.[34m↵[0m
[34m[38;5;231mcom/go/firestore v1.9.0[0m                              [34m│ [38;5;231mcom/go/firestore v1.9.0[0m
[34m[38;5;231m    cloud.google.com/go/iam => cloud.google.com/go[34m↴[0m  [34m│ [38;5;231m    cloud.google.com/go/iam => cloud.google.com/go[34m↴[0m
[34m[0m                                      [34m…[38;5;231m/iam v0.10.0[0m  [34m│ [0m                                      [34m…[38;5;231m/iam v0.10.0[0m
[34m[38;5;231m    cloud.google.com/go/kms => cloud.google.com/go[34m↴[0m  [34m│ [38;5;231m    cloud.google.com/go/kms => cloud.google.com/go[34m↴[0m
[34m[0m                                       [34m…[38;5;231m/kms v1.8.0[0m  [34m│ [0m                                       [34m…[38;5;231m/kms v1.8.0[0m
[34m[48;5;52;38;5;231m    cloud.google.com/go/monitoring => cloud.google.c[34m↵[0m[34m│ [48;5;22;38;5;231m    cloud.google.com/go/monitoring => cloud.google[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mom/go/monitoring v1.[48;5;124m11[48;5;52m.0[0m[48;5;52m                             [0m[34m│ [48;5;22;38;5;231m.com/go/monitoring v1.[48;5;28m12[48;5;22m.0[0m[48;5;22m[0K[0m
[34m[38;5;231m    cloud.google.com/go/pubsub => cloud.google.com[34m↴[0m  [34m│ [38;5;231m    cloud.google.com/go/pubsub => cloud.google.com[34m↴[0m
[34m[0m                                [34m…[38;5;231m/go/pubsub v1.28.0[0m  [34m│ [0m                                [34m…[38;5;231m/go/pubsub v1.28.0[0m
[34m[38;5;231m    cloud.google.com/go/secretmanager => cloud.goo[34m↵[0m  [34m│ [38;5;231m    cloud.google.com/go/secretmanager => cloud.goo[34m↵[0m
[34m[38;5;231mgle.com/go/secretmanager v1.10.0[0m                     [34m│ [38;5;231mgle.com/go/secretmanager v1.10.0[0m
[34m[48;5;52;38;5;231m    cloud.google.com/go/storage => cloud.google.com/[34m↴[0m[34m│ [48;5;22;38;5;231m    cloud.google.com/go/storage => cloud.google.co[34m↵[0m[48;5;22m[0K[0m
[34m[0m[48;5;52m                                  [34m…[38;5;231mgo/storage v1.[48;5;124m28[48;5;52m.[48;5;124m1[0m[34m│ [48;5;22;38;5;231mm/go/storage v1.[48;5;28m29[48;5;22m.[48;5;28m0[0m[48;5;22m[0K[0m
[34m[38;5;231m    cloud.google.com/go/trace => cloud.google.com/[34m↴[0m  [34m│ [38;5;231m    cloud.google.com/go/trace => cloud.google.com/[34m↴[0m
[34m[0m                                   [34m…[38;5;231mgo/trace v1.5.0[0m  [34m│ [0m                                   [34m…[38;5;231mgo/trace v1.5.0[0m
[34m[38;5;231m    code.cloudfoundry.org/bytefmt => code.cloudfou[34m↵[0m  [34m│ [38;5;231m    code.cloudfoundry.org/bytefmt => code.cloudfou[34m↵[0m
[34m[38;5;231mndry.org/bytefmt v0.0.0-20211005130812-5bb3c17173e5[0m  [34m│ [38;5;231mndry.org/bytefmt v0.0.0-20211005130812-5bb3c17173e5[0m
[34m[38;5;231m    contrib.go.opencensus.io/exporter/aws => contr[34m↵[0m  [34m│ [38;5;231m    contrib.go.opencensus.io/exporter/aws => contr[34m↵[0m
[34m[38;5;231mib.go.opencensus.io/exporter/aws v0.0.0-2020061720[34m↵[0m  [34m│ [38;5;231mib.go.opencensus.io/exporter/aws v0.0.0-2020061720[34m↵[0m
[34m[38;5;231m4711-c478e41e60e9[0m                                    [34m│ [38;5;231m4711-c478e41e60e9[0m

[36m──────────────[0m[36m┐[0m
[34m22[0m:[38;5;231m replace ( [0m[36m│[0m
[36m──────────────[0m[36m┴[0m[36m───────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[38;5;231m    git.sr.ht/~sbinet/gg => git.sr.ht/~sbinet/gg v[34m↴[0m  [34m│ [38;5;231m    git.sr.ht/~sbinet/gg => git.sr.ht/~sbinet/gg v[34m↴[0m
[34m[0m                                             [34m…[38;5;231m0.3.1[0m  [34m│ [0m                                             [34m…[38;5;231m0.3.1[0m
[34m[38;5;231m    github.com/AdaLogics/go-fuzz-headers => github[34m↵[0m  [34m│ [38;5;231m    github.com/AdaLogics/go-fuzz-headers => github[34m↵[0m
[34m[38;5;231m.com/AdaLogics/go-fuzz-headers v0.0.0-202301062348[34m↵[0m  [34m│ [38;5;231m.com/AdaLogics/go-fuzz-headers v0.0.0-202301062348[34m↵[0m
[34m[38;5;231m47-43070de90fa1[0m                                      [34m│ [38;5;231m47-43070de90fa1[0m
[34m[38;5;231m    github.com/Azure/azure-amqp-common-go/v3 => gi[34m↵[0m  [34m│ [38;5;231m    github.com/Azure/azure-amqp-common-go/v3 => gi[34m↵[0m
[34m[38;5;231mthub.com/Azure/azure-amqp-common-go/v3 v3.2.3[0m        [34m│ [38;5;231mthub.com/Azure/azure-amqp-common-go/v3 v3.2.3[0m
[34m[48;5;52;38;5;231m    github.com/Azure/azure-sdk-for-go => github.com/[34m↵[0m[34m│ [48;5;22;38;5;231m    github.com/Azure/azure-sdk-for-go => github.co[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mAzure/azure-sdk-for-go [48;5;124mv67[48;5;52m.[48;5;124m3[48;5;52m.0+incompatible[0m[48;5;52m          [0m[34m│ [48;5;22;38;5;231mm/Azure/azure-sdk-for-go [48;5;28mv68[48;5;22m.[48;5;28m0[48;5;22m.0+incompatible[0m[48;5;22m[0K[0m
[34m[38;5;231m    github.com/Azure/azure-sdk-for-go/sdk/azcore =[34m↵[0m  [34m│ [38;5;231m    github.com/Azure/azure-sdk-for-go/sdk/azcore =[34m↵[0m
[34m[38;5;231m> github.com/Azure/azure-sdk-for-go/sdk/azcore v1.[34m↵[0m  [34m│ [38;5;231m> github.com/Azure/azure-sdk-for-go/sdk/azcore v1.[34m↵[0m
[34m[38;5;231m3.0[0m                                                  [34m│ [38;5;231m3.0[0m
[34m[38;5;231m    github.com/Azure/azure-sdk-for-go/sdk/azidenti[34m↵[0m  [34m│ [38;5;231m    github.com/Azure/azure-sdk-for-go/sdk/azidenti[34m↵[0m
[34m[38;5;231mty => github.com/Azure/azure-sdk-for-go/sdk/aziden[34m↵[0m  [34m│ [38;5;231mty => github.com/Azure/azure-sdk-for-go/sdk/aziden[34m↵[0m
[34m[38;5;231mtity v1.2.0[0m                                          [34m│ [38;5;231mtity v1.2.0[0m
[34m[38;5;231m    github.com/Azure/azure-sdk-for-go/sdk/internal[34m↵[0m  [34m│ [38;5;231m    github.com/Azure/azure-sdk-for-go/sdk/internal[34m↵[0m
[34m[38;5;231m => github.com/Azure/azure-sdk-for-go/sdk/internal[34m↵[0m  [34m│ [38;5;231m => github.com/Azure/azure-sdk-for-go/sdk/internal[34m↵[0m
[34m[38;5;231m v1.1.2[0m                                              [34m│ [38;5;231m v1.1.2[0m

[36m──────────────[0m[36m┐[0m
[34m60[0m:[38;5;231m replace ( [0m[36m│[0m
[36m──────────────[0m[36m┴[0m[36m───────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[38;5;231m    github.com/armon/go-radix => github.com/armon/[34m↴[0m  [34m│ [38;5;231m    github.com/armon/go-radix => github.com/armon/[34m↴[0m
[34m[0m                                   [34m…[38;5;231mgo-radix v1.0.0[0m  [34m│ [0m                                   [34m…[38;5;231mgo-radix v1.0.0[0m
[34m[38;5;231m    github.com/armon/go-socks5 => github.com/armon[34m↵[0m  [34m│ [38;5;231m    github.com/armon/go-socks5 => github.com/armon[34m↵[0m
[34m[38;5;231m/go-socks5 v0.0.0-20160902184237-e75332964ef5[0m        [34m│ [38;5;231m/go-socks5 v0.0.0-20160902184237-e75332964ef5[0m
[34m[38;5;231m    github.com/asaskevich/govalidator => github.co[34m↵[0m  [34m│ [38;5;231m    github.com/asaskevich/govalidator => github.co[34m↵[0m
[34m[38;5;231mm/asaskevich/govalidator v0.0.0-20210307081110-f21[34m↵[0m  [34m│ [38;5;231mm/asaskevich/govalidator v0.0.0-20210307081110-f21[34m↵[0m
[34m[38;5;231m760c49a8d[0m                                            [34m│ [38;5;231m760c49a8d[0m
[34m[48;5;52;38;5;231m    github.com/aws/aws-sdk-go => github.com/aws/aws-[34m↴[0m[34m│ [48;5;22;38;5;231m    github.com/aws/aws-sdk-go => github.com/aws/aw[34m↴[0m[48;5;22m[0K[0m
[34m[0m[48;5;52m                                    [34m…[38;5;231msdk-go v1.44.[48;5;124m181[0m[34m│ [0m[48;5;22m                                [34m…[38;5;231ms-sdk-go v1.44.[48;5;28m185[0m[48;5;22m[0K[0m
[34m[38;5;231m    github.com/aws/aws-sdk-go-v2 => github.com/aws[34m↵[0m  [34m│ [38;5;231m    github.com/aws/aws-sdk-go-v2 => github.com/aws[34m↵[0m
[34m[38;5;231m/aws-sdk-go-v2 v1.17.3[0m                               [34m│ [38;5;231m/aws-sdk-go-v2 v1.17.3[0m
[34m[38;5;231m    github.com/aws/aws-sdk-go-v2/aws/protocol/even[34m↵[0m  [34m│ [38;5;231m    github.com/aws/aws-sdk-go-v2/aws/protocol/even[34m↵[0m
[34m[38;5;231mtstream => github.com/aws/aws-sdk-go-v2/aws/protoc[34m↵[0m  [34m│ [38;5;231mtstream => github.com/aws/aws-sdk-go-v2/aws/protoc[34m↵[0m
[34m[38;5;231mol/eventstream v1.4.10[0m                               [34m│ [38;5;231mol/eventstream v1.4.10[0m
[34m[48;5;52;38;5;231m    github.com/aws/aws-sdk-go-v2/config => github.co[34m↵[0m[34m│ [48;5;22;38;5;231m    github.com/aws/aws-sdk-go-v2/config => github.[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mm/aws/aws-sdk-go-v2/config v1.18.[48;5;124m8[0m[48;5;52m                   [0m[34m│ [48;5;22;38;5;231mcom/aws/aws-sdk-go-v2/config v1.18.[48;5;28m9[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231m    github.com/aws/aws-sdk-go-v2/credentials => gith[34m↵[0m[34m│ [48;5;22;38;5;231m    github.com/aws/aws-sdk-go-v2/credentials => gi[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mub.com/aws/aws-sdk-go-v2/credentials v1.13.[48;5;124m8[0m[48;5;52m         [0m[34m│ [48;5;22;38;5;231mthub.com/aws/aws-sdk-go-v2/credentials v1.13.[48;5;28m9[0m[48;5;22m[0K[0m
[34m[38;5;231m    github.com/aws/aws-sdk-go-v2/feature/ec2/imds [34m↵[0m  [34m│ [38;5;231m    github.com/aws/aws-sdk-go-v2/feature/ec2/imds [34m↵[0m
[34m[38;5;231m=> github.com/aws/aws-sdk-go-v2/feature/ec2/imds v[34m↵[0m  [34m│ [38;5;231m=> github.com/aws/aws-sdk-go-v2/feature/ec2/imds v[34m↵[0m
[34m[38;5;231m1.12.21[0m                                              [34m│ [38;5;231m1.12.21[0m
[34m[48;5;52;38;5;231m    github.com/aws/aws-sdk-go-v2/feature/s3/manager [34m↵[0m[34m│ [48;5;22;38;5;231m    github.com/aws/aws-sdk-go-v2/feature/s3/manage[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231m=> github.com/aws/aws-sdk-go-v2/feature/s3/manager v[34m↵[0m[34m│ [48;5;22;38;5;231mr => github.com/aws/aws-sdk-go-v2/feature/s3/manag[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231m1.11.[48;5;124m47[0m[48;5;52m                                              [0m[34m│ [48;5;22;38;5;231mer v1.11.[48;5;28m48[0m[48;5;22m[0K[0m
[34m[38;5;231m    github.com/aws/aws-sdk-go-v2/internal/configso[34m↵[0m  [34m│ [38;5;231m    github.com/aws/aws-sdk-go-v2/internal/configso[34m↵[0m
[34m[38;5;231murces => github.com/aws/aws-sdk-go-v2/internal/con[34m↵[0m  [34m│ [38;5;231murces => github.com/aws/aws-sdk-go-v2/internal/con[34m↵[0m
[34m[38;5;231mfigsources v1.1.27[0m                                   [34m│ [38;5;231mfigsources v1.1.27[0m
[34m[38;5;231m    github.com/aws/aws-sdk-go-v2/internal/endpoint[34m↵[0m  [34m│ [38;5;231m    github.com/aws/aws-sdk-go-v2/internal/endpoint[34m↵[0m
[34m[38;5;231ms/v2 => github.com/aws/aws-sdk-go-v2/internal/endp[34m↵[0m  [34m│ [38;5;231ms/v2 => github.com/aws/aws-sdk-go-v2/internal/endp[34m↵[0m
[34m[38;5;231moints/v2 v2.4.21[0m                                     [34m│ [38;5;231moints/v2 v2.4.21[0m
[34m[38;5;231m    github.com/aws/aws-sdk-go-v2/internal/ini => g[34m↵[0m  [34m│ [38;5;231m    github.com/aws/aws-sdk-go-v2/internal/ini => g[34m↵[0m
[34m[38;5;231mithub.com/aws/aws-sdk-go-v2/internal/ini v1.3.28[0m     [34m│ [38;5;231mithub.com/aws/aws-sdk-go-v2/internal/ini v1.3.28[0m

[36m──────────────[0m[36m┐[0m
[34m74[0m:[38;5;231m replace ( [0m[36m│[0m
[36m──────────────[0m[36m┴[0m[36m───────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[38;5;231m    github.com/aws/aws-sdk-go-v2/service/internal/[34m↵[0m  [34m│ [38;5;231m    github.com/aws/aws-sdk-go-v2/service/internal/[34m↵[0m
[34m[38;5;231mchecksum => github.com/aws/aws-sdk-go-v2/service/i[34m↵[0m  [34m│ [38;5;231mchecksum => github.com/aws/aws-sdk-go-v2/service/i[34m↵[0m
[34m[38;5;231mnternal/checksum v1.1.22[0m                             [34m│ [38;5;231mnternal/checksum v1.1.22[0m
[34m[38;5;231m    github.com/aws/aws-sdk-go-v2/service/internal/[34m↵[0m  [34m│ [38;5;231m    github.com/aws/aws-sdk-go-v2/service/internal/[34m↵[0m
[34m[38;5;231mpresigned-url => github.com/aws/aws-sdk-go-v2/serv[34m↵[0m  [34m│ [38;5;231mpresigned-url => github.com/aws/aws-sdk-go-v2/serv[34m↵[0m
[34m[38;5;231mice/internal/presigned-url v1.9.21[0m                   [34m│ [38;5;231mice/internal/presigned-url v1.9.21[0m
[34m[38;5;231m    github.com/aws/aws-sdk-go-v2/service/internal/[34m↵[0m  [34m│ [38;5;231m    github.com/aws/aws-sdk-go-v2/service/internal/[34m↵[0m
[34m[38;5;231ms3shared => github.com/aws/aws-sdk-go-v2/service/i[34m↵[0m  [34m│ [38;5;231ms3shared => github.com/aws/aws-sdk-go-v2/service/i[34m↵[0m
[34m[38;5;231mnternal/s3shared v1.13.21[0m                            [34m│ [38;5;231mnternal/s3shared v1.13.21[0m
[34m[48;5;52;38;5;231m    github.com/aws/aws-sdk-go-v2/service/kms => gith[34m↵[0m[34m│ [48;5;22;38;5;231m    github.com/aws/aws-sdk-go-v2/service/kms => gi[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mub.com/aws/aws-sdk-go-v2/service/kms v1.20.[48;5;124m0[0m[48;5;52m         [0m[34m│ [48;5;22;38;5;231mthub.com/aws/aws-sdk-go-v2/service/kms v1.20.[48;5;28m1[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231m    github.com/aws/aws-sdk-go-v2/service/s3 => githu[34m↵[0m[34m│ [48;5;22;38;5;231m    github.com/aws/aws-sdk-go-v2/service/s3 => git[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mb.com/aws/aws-sdk-go-v2/service/s3 v1.30.[48;5;124m0[0m[48;5;52m           [0m[34m│ [48;5;22;38;5;231mhub.com/aws/aws-sdk-go-v2/service/s3 v1.30.[48;5;28m1[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231m    github.com/aws/aws-sdk-go-v2/service/secretsmana[34m↵[0m[34m│ [48;5;22;38;5;231m    github.com/aws/aws-sdk-go-v2/service/secretsma[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mger => github.com/aws/aws-sdk-go-v2/service/secretsm[34m↵[0m[34m│ [48;5;22;38;5;231mnager => github.com/aws/aws-sdk-go-v2/service/secr[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231manager v1.18.[48;5;124m1[0m[48;5;52m                                       [0m[34m│ [48;5;22;38;5;231metsmanager v1.18.[48;5;28m2[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231m    github.com/aws/aws-sdk-go-v2/service/sns => gith[34m↵[0m[34m│ [48;5;22;38;5;231m    github.com/aws/aws-sdk-go-v2/service/sns => gi[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mub.com/aws/aws-sdk-go-v2/service/sns v1.19.[48;5;124m0[0m[48;5;52m         [0m[34m│ [48;5;22;38;5;231mthub.com/aws/aws-sdk-go-v2/service/sns v1.19.[48;5;28m1[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231m    github.com/aws/aws-sdk-go-v2/service/sqs => gith[34m↵[0m[34m│ [48;5;22;38;5;231m    github.com/aws/aws-sdk-go-v2/service/sqs => gi[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mub.com/aws/aws-sdk-go-v2/service/sqs v1.20.[48;5;124m0[0m[48;5;52m         [0m[34m│ [48;5;22;38;5;231mthub.com/aws/aws-sdk-go-v2/service/sqs v1.20.[48;5;28m1[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231m    github.com/aws/aws-sdk-go-v2/service/ssm => gith[34m↵[0m[34m│ [48;5;22;38;5;231m    github.com/aws/aws-sdk-go-v2/service/ssm => gi[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mub.com/aws/aws-sdk-go-v2/service/ssm v1.35.[48;5;124m0[0m[48;5;52m         [0m[34m│ [48;5;22;38;5;231mthub.com/aws/aws-sdk-go-v2/service/ssm v1.35.[48;5;28m1[0m[48;5;22m[0K[0m
[34m[38;5;231m    github.com/aws/aws-sdk-go-v2/service/sso => gi[34m↵[0m  [34m│ [38;5;231m    github.com/aws/aws-sdk-go-v2/service/sso => gi[34m↵[0m
[34m[38;5;231mthub.com/aws/aws-sdk-go-v2/service/sso v1.12.0[0m       [34m│ [38;5;231mthub.com/aws/aws-sdk-go-v2/service/sso v1.12.0[0m
[34m[48;5;52;38;5;231m    github.com/aws/aws-sdk-go-v2/service/sts => gith[34m↵[0m[34m│ [48;5;22;38;5;231m    github.com/aws/aws-sdk-go-v2/service/sts => gi[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mub.com/aws/aws-sdk-go-v2/service/sts v1.18.[48;5;124m0[0m[48;5;52m         [0m[34m│ [48;5;22;38;5;231mthub.com/aws/aws-sdk-go-v2/service/sts v1.18.[48;5;28m1[0m[48;5;22m[0K[0m
[34m[38;5;231m    github.com/aws/smithy-go => github.com/aws/smi[34m↴[0m  [34m│ [38;5;231m    github.com/aws/smithy-go => github.com/aws/smi[34m↴[0m
[34m[0m                                    [34m…[38;5;231mthy-go v1.13.5[0m  [34m│ [0m                                    [34m…[38;5;231mthy-go v1.13.5[0m
[34m[38;5;231m    github.com/benbjohnson/clock => github.com/ben[34m↵[0m  [34m│ [38;5;231m    github.com/benbjohnson/clock => github.com/ben[34m↵[0m
[34m[38;5;231mbjohnson/clock v1.3.0[0m                                [34m│ [38;5;231mbjohnson/clock v1.3.0[0m
[34m[38;5;231m    github.com/beorn7/perks => github.com/beorn7/p[34m↴[0m  [34m│ [38;5;231m    github.com/beorn7/perks => github.com/beorn7/p[34m↴[0m
[34m[0m                                       [34m…[38;5;231merks v1.0.1[0m  [34m│ [0m                                       [34m…[38;5;231merks v1.0.1[0m

[36m───────────────[0m[36m┐[0m
[34m123[0m:[38;5;231m replace ( [0m[36m│[0m
[36m───────────────[0m[36m┴[0m[36m──────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[38;5;231m    github.com/containerd/zfs => github.com/contai[34m↴[0m  [34m│ [38;5;231m    github.com/containerd/zfs => github.com/contai[34m↴[0m
[34m[0m                                   [34m…[38;5;231mnerd/zfs v1.0.0[0m  [34m│ [0m                                   [34m…[38;5;231mnerd/zfs v1.0.0[0m
[34m[38;5;231m    github.com/containernetworking/cni => github.c[34m↵[0m  [34m│ [38;5;231m    github.com/containernetworking/cni => github.c[34m↵[0m
[34m[38;5;231mom/containernetworking/cni v1.1.2[0m                    [34m│ [38;5;231mom/containernetworking/cni v1.1.2[0m
[34m[38;5;231m    github.com/containernetworking/plugins => gith[34m↵[0m  [34m│ [38;5;231m    github.com/containernetworking/plugins => gith[34m↵[0m
[34m[38;5;231mub.com/containernetworking/plugins v1.2.0[0m            [34m│ [38;5;231mub.com/containernetworking/plugins v1.2.0[0m
[34m[48;5;52;38;5;231m    github.com/containers/ocicrypt => github.com/con[34m↵[0m[34m│ [48;5;22;38;5;231m    github.com/containers/ocicrypt => github.com/c[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mtainers/ocicrypt v1.1.[48;5;124m6[0m[48;5;52m                              [0m[34m│ [48;5;22;38;5;231montainers/ocicrypt v1.1.[48;5;28m7[0m[48;5;22m[0K[0m
[34m[38;5;231m    github.com/coreos/go-iptables => github.com/co[34m↵[0m  [34m│ [38;5;231m    github.com/coreos/go-iptables => github.com/co[34m↵[0m
[34m[38;5;231mreos/go-iptables v0.6.0[0m                              [34m│ [38;5;231mreos/go-iptables v0.6.0[0m
[34m[38;5;231m    github.com/coreos/go-oidc => github.com/coreos[34m↵[0m  [34m│ [38;5;231m    github.com/coreos/go-oidc => github.com/coreos[34m↵[0m
[34m[38;5;231m/go-oidc v2.2.1+incompatible[0m                         [34m│ [38;5;231m/go-oidc v2.2.1+incompatible[0m
[34m[48;5;52;38;5;231m    github.com/coreos/go-semver => github.com/coreos[34m↴[0m[34m│ [48;5;22;38;5;231m    github.com/coreos/go-semver => github.com/core[34m↵[0m[48;5;22m[0K[0m
[34m[0m[48;5;52m                                   [34m…[38;5;231m/go-semver v0.3.[48;5;124m0[0m[34m│ [48;5;22;38;5;231mos/go-semver v0.3.[48;5;28m1[0m[48;5;22m[0K[0m
[34m[38;5;231m    github.com/coreos/go-systemd/v22 => github.com[34m↵[0m  [34m│ [38;5;231m    github.com/coreos/go-systemd/v22 => github.com[34m↵[0m
[34m[38;5;231m/coreos/go-systemd/v22 v22.5.0[0m                       [34m│ [38;5;231m/coreos/go-systemd/v22 v22.5.0[0m
[34m[38;5;231m    github.com/cpuguy83/go-md2man/v2 => github.com[34m↵[0m  [34m│ [38;5;231m    github.com/cpuguy83/go-md2man/v2 => github.com[34m↵[0m
[34m[38;5;231m/cpuguy83/go-md2man/v2 v2.0.2[0m                        [34m│ [38;5;231m/cpuguy83/go-md2man/v2 v2.0.2[0m
[34m[38;5;231m    github.com/creack/pty => github.com/creack/pty[34m↴[0m  [34m│ [38;5;231m    github.com/creack/pty => github.com/creack/pty[34m↴[0m
[34m[0m                                          [34m…[38;5;231m v1.1.18[0m  [34m│ [0m                                          [34m…[38;5;231m v1.1.18[0m

[36m───────────────[0m[36m┐[0m
[34m140[0m:[38;5;231m replace ( [0m[36m│[0m
[36m───────────────[0m[36m┴[0m[36m──────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[38;5;231m    github.com/devigned/tab => github.com/devigned[34m↴[0m  [34m│ [38;5;231m    github.com/devigned/tab => github.com/devigned[34m↴[0m
[34m[0m                                       [34m…[38;5;231m/tab v0.1.1[0m  [34m│ [0m                                       [34m…[38;5;231m/tab v0.1.1[0m
[34m[38;5;231m    github.com/dgryski/go-rendezvous => github.com[34m↵[0m  [34m│ [38;5;231m    github.com/dgryski/go-rendezvous => github.com[34m↵[0m
[34m[38;5;231m/dgryski/go-rendezvous v0.0.0-20200823014737-9f700[34m↵[0m  [34m│ [38;5;231m/dgryski/go-rendezvous v0.0.0-20200823014737-9f700[34m↵[0m
[34m[38;5;231m1d12a5f[0m                                              [34m│ [38;5;231m1d12a5f[0m
[34m[38;5;231m    github.com/dgryski/go-sip13 => github.com/dgry[34m↵[0m  [34m│ [38;5;231m    github.com/dgryski/go-sip13 => github.com/dgry[34m↵[0m
[34m[38;5;231mski/go-sip13 v0.0.0-20200911182023-62edffca9245[0m      [34m│ [38;5;231mski/go-sip13 v0.0.0-20200911182023-62edffca9245[0m
[34m[48;5;52;38;5;231m    github.com/digitalocean/godo => github.com/digit[34m↵[0m[34m│ [48;5;22;38;5;231m    github.com/digitalocean/godo => github.com/dig[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231malocean/godo v1.[48;5;124m93[48;5;52m.0[0m[48;5;52m                                 [0m[34m│ [48;5;22;38;5;231mitalocean/godo v1.[48;5;28m94[48;5;22m.0[0m[48;5;22m[0K[0m
[34m[38;5;231m    github.com/dimchansky/utfbom => github.com/dim[34m↵[0m  [34m│ [38;5;231m    github.com/dimchansky/utfbom => github.com/dim[34m↵[0m
[34m[38;5;231mchansky/utfbom v1.1.1[0m                                [34m│ [38;5;231mchansky/utfbom v1.1.1[0m
[34m[38;5;231m    github.com/dnaeon/go-vcr => github.com/dnaeon/[34m↴[0m  [34m│ [38;5;231m    github.com/dnaeon/go-vcr => github.com/dnaeon/[34m↴[0m
[34m[0m                                     [34m…[38;5;231mgo-vcr v1.2.0[0m  [34m│ [0m                                     [34m…[38;5;231mgo-vcr v1.2.0[0m
[34m[48;5;52;38;5;231m    github.com/docker/cli => github.com/docker/cli v[34m↵[0m[34m│ [48;5;22;38;5;231m    github.com/docker/cli => github.com/docker/cli[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231m20.10.[48;5;124m22[48;5;52m+incompatible[0m[48;5;52m                                [0m[34m│ [48;5;22;38;5;231m v20.10.[48;5;28m23[48;5;22m+incompatible[0m[48;5;22m[0K[0m
[34m[38;5;231m    github.com/docker/distribution => github.com/d[34m↵[0m  [34m│ [38;5;231m    github.com/docker/distribution => github.com/d[34m↵[0m
[34m[38;5;231mocker/distribution v2.8.1+incompatible[0m               [34m│ [38;5;231mocker/distribution v2.8.1+incompatible[0m
[34m[48;5;52;38;5;231m    github.com/docker/docker => github.com/docker/do[34m↵[0m[34m│ [48;5;22;38;5;231m    github.com/docker/docker => github.com/docker/[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mcker v20.10.[48;5;124m22[48;5;52m+incompatible[0m[48;5;52m                          [0m[34m│ [48;5;22;38;5;231mdocker v20.10.[48;5;28m23[48;5;22m+incompatible[0m[48;5;22m[0K[0m
[34m[38;5;231m    github.com/docker/docker-credential-helpers =>[34m↵[0m  [34m│ [38;5;231m    github.com/docker/docker-credential-helpers =>[34m↵[0m
[34m[38;5;231m github.com/docker/docker-credential-helpers v0.7.0[0m  [34m│ [38;5;231m github.com/docker/docker-credential-helpers v0.7.0[0m
[34m[38;5;231m    github.com/docker/go-connections => github.com[34m↵[0m  [34m│ [38;5;231m    github.com/docker/go-connections => github.com[34m↵[0m
[34m[38;5;231m/docker/go-connections v0.4.0[0m                        [34m│ [38;5;231m/docker/go-connections v0.4.0[0m
[34m[38;5;231m    github.com/docker/go-events => github.com/dock[34m↵[0m  [34m│ [38;5;231m    github.com/docker/go-events => github.com/dock[34m↵[0m
[34m[38;5;231mer/go-events v0.0.0-20190806004212-e31b211e4f1c[0m      [34m│ [38;5;231mer/go-events v0.0.0-20190806004212-e31b211e4f1c[0m

[36m───────────────[0m[36m┐[0m
[34m161[0m:[38;5;231m replace ( [0m[36m│[0m
[36m───────────────[0m[36m┴[0m[36m──────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[38;5;231m    github.com/envoyproxy/go-control-plane => gith[34m↵[0m  [34m│ [38;5;231m    github.com/envoyproxy/go-control-plane => gith[34m↵[0m
[34m[38;5;231mub.com/envoyproxy/go-control-plane v0.10.3[0m           [34m│ [38;5;231mub.com/envoyproxy/go-control-plane v0.10.3[0m
[34m[38;5;231m    github.com/envoyproxy/protoc-gen-validate => g[34m↵[0m  [34m│ [38;5;231m    github.com/envoyproxy/protoc-gen-validate => g[34m↵[0m
[34m[38;5;231mithub.com/envoyproxy/protoc-gen-validate v0.9.1[0m      [34m│ [38;5;231mithub.com/envoyproxy/protoc-gen-validate v0.9.1[0m
[34m[38;5;231m    github.com/evanphx/json-patch => github.com/ev[34m↵[0m  [34m│ [38;5;231m    github.com/evanphx/json-patch => github.com/ev[34m↵[0m
[34m[38;5;231manphx/json-patch v0.5.2[0m                              [34m│ [38;5;231manphx/json-patch v0.5.2[0m
[34m[48;5;52;38;5;231m    github.com/fatih/color => github.com/fatih/color[34m↴[0m[34m│ [48;5;22;38;5;231m    github.com/fatih/color => github.com/fatih/col[34m↴[0m[48;5;22m[0K[0m
[34m[0m[48;5;52m                                            [34m…[38;5;231m v1.[48;5;124m13[48;5;52m.[48;5;124m0[0m[34m│ [0m[48;5;22m                                        [34m…[38;5;231mor v1.[48;5;28m14[48;5;22m.[48;5;28m1[0m[48;5;22m[0K[0m
[34m[38;5;231m    github.com/felixge/httpsnoop => github.com/fel[34m↵[0m  [34m│ [38;5;231m    github.com/felixge/httpsnoop => github.com/fel[34m↵[0m
[34m[38;5;231mixge/httpsnoop v1.0.3[0m                                [34m│ [38;5;231mixge/httpsnoop v1.0.3[0m
[34m[38;5;231m    github.com/fogleman/gg => github.com/fogleman/[34m↴[0m  [34m│ [38;5;231m    github.com/fogleman/gg => github.com/fogleman/[34m↴[0m
[34m[0m                                         [34m…[38;5;231mgg v1.3.0[0m  [34m│ [0m                                         [34m…[38;5;231mgg v1.3.0[0m
[34m[38;5;231m    github.com/form3tech-oss/jwt-go => github.com/[34m↵[0m  [34m│ [38;5;231m    github.com/form3tech-oss/jwt-go => github.com/[34m↵[0m
[34m[38;5;231mform3tech-oss/jwt-go v3.2.5+incompatible[0m             [34m│ [38;5;231mform3tech-oss/jwt-go v3.2.5+incompatible[0m

[36m───────────────[0m[36m┐[0m
[34m308[0m:[38;5;231m replace ( [0m[36m│[0m
[36m───────────────[0m[36m┴[0m[36m──────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[38;5;231m    github.com/julienschmidt/httprouter => github.[34m↵[0m  [34m│ [38;5;231m    github.com/julienschmidt/httprouter => github.[34m↵[0m
[34m[38;5;231mcom/julienschmidt/httprouter v1.3.0[0m                  [34m│ [38;5;231mcom/julienschmidt/httprouter v1.3.0[0m
[34m[38;5;231m    github.com/kisielk/errcheck => github.com/kisi[34m↵[0m  [34m│ [38;5;231m    github.com/kisielk/errcheck => github.com/kisi[34m↵[0m
[34m[38;5;231melk/errcheck v1.6.3[0m                                  [34m│ [38;5;231melk/errcheck v1.6.3[0m
[34m[38;5;231m    github.com/kisielk/gotool => github.com/kisiel[34m↴[0m  [34m│ [38;5;231m    github.com/kisielk/gotool => github.com/kisiel[34m↴[0m
[34m[0m                                   [34m…[38;5;231mk/gotool v1.0.0[0m  [34m│ [0m                                   [34m…[38;5;231mk/gotool v1.0.0[0m
[34m[48;5;52;38;5;231m    github.com/klauspost/compress => github.com/klau[34m↵[0m[34m│ [48;5;22;38;5;231m    github.com/klauspost/compress => github.com/kl[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mspost/compress v1.15.[48;5;124m15[48;5;52m-0.[48;5;124m20230116143836[48;5;52m-[48;5;124mfbae784ff625[0m[34m│ [48;5;22;38;5;231mauspost/compress v1.15.[48;5;28m16[48;5;22m-0.[48;5;28m20230121171712[48;5;22m-[48;5;28mfe37dc6[48;5;22;34m↵[0m[48;5;22m[0K[0m
[34m[0m                                                     [34m│ [48;5;28;38;5;231m783c8[0m[48;5;22m[0K[0m
[34m[38;5;231m    github.com/klauspost/cpuid/v2 => github.com/kl[34m↵[0m  [34m│ [38;5;231m    github.com/klauspost/cpuid/v2 => github.com/kl[34m↵[0m
[34m[38;5;231mauspost/cpuid/v2 v2.2.3[0m                              [34m│ [38;5;231mauspost/cpuid/v2 v2.2.3[0m
[34m[38;5;231m    github.com/kolo/xmlrpc => github.com/kolo/xmlr[34m↵[0m  [34m│ [38;5;231m    github.com/kolo/xmlrpc => github.com/kolo/xmlr[34m↵[0m
[34m[38;5;231mpc v0.0.0-20220921171641-a4b6fa1dd06b[0m                [34m│ [38;5;231mpc v0.0.0-20220921171641-a4b6fa1dd06b[0m
[34m[38;5;231m    github.com/kpango/fastime => github.com/kpango[34m↴[0m  [34m│ [38;5;231m    github.com/kpango/fastime => github.com/kpango[34m↴[0m
[34m[0m                                   [34m…[38;5;231m/fastime v1.1.6[0m  [34m│ [0m                                   [34m…[38;5;231m/fastime v1.1.6[0m

[36m───────────────[0m[36m┐[0m
[34m323[0m:[38;5;231m replace ( [0m[36m│[0m
[36m───────────────[0m[36m┴[0m[36m──────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[38;5;231m    github.com/leodido/go-urn => github.com/leodid[34m↴[0m  [34m│ [38;5;231m    github.com/leodido/go-urn => github.com/leodid[34m↴[0m
[34m[0m                                   [34m…[38;5;231mo/go-urn v1.2.1[0m  [34m│ [0m                                   [34m…[38;5;231mo/go-urn v1.2.1[0m
[34m[38;5;231m    github.com/lib/pq => github.com/lib/pq v1.10.7[0m   [34m│ [38;5;231m    github.com/lib/pq => github.com/lib/pq v1.10.7[0m
[34m[38;5;231m    github.com/liggitt/tabwriter => github.com/lig[34m↵[0m  [34m│ [38;5;231m    github.com/liggitt/tabwriter => github.com/lig[34m↵[0m
[34m[38;5;231mgitt/tabwriter v0.0.0-20181228230101-89fcab3d43de[0m    [34m│ [38;5;231mgitt/tabwriter v0.0.0-20181228230101-89fcab3d43de[0m
[34m[48;5;52;38;5;231m    github.com/linode/linodego => github.com/linode/[34m↴[0m[34m│ [48;5;22;38;5;231m    github.com/linode/linodego => github.com/linod[34m↴[0m[48;5;22m[0K[0m
[34m[0m[48;5;52m                                    [34m…[38;5;231mlinodego v1.[48;5;124m11[48;5;52m.0[0m[34m│ [0m[48;5;22m                                [34m…[38;5;231me/linodego v1.[48;5;28m12[48;5;22m.0[0m[48;5;22m[0K[0m
[34m[38;5;231m    github.com/linuxkit/virtsock => github.com/lin[34m↵[0m  [34m│ [38;5;231m    github.com/linuxkit/virtsock => github.com/lin[34m↵[0m
[34m[38;5;231muxkit/virtsock v0.0.0-20220523201153-1a23e78aa7a2[0m    [34m│ [38;5;231muxkit/virtsock v0.0.0-20220523201153-1a23e78aa7a2[0m
[34m[38;5;231m    github.com/lucasb-eyer/go-colorful => github.c[34m↵[0m  [34m│ [38;5;231m    github.com/lucasb-eyer/go-colorful => github.c[34m↵[0m
[34m[38;5;231mom/lucasb-eyer/go-colorful v1.2.0[0m                    [34m│ [38;5;231mom/lucasb-eyer/go-colorful v1.2.0[0m
[34m[38;5;231m    github.com/lyft/protoc-gen-star => github.com/[34m↵[0m  [34m│ [38;5;231m    github.com/lyft/protoc-gen-star => github.com/[34m↵[0m
[34m[38;5;231mlyft/protoc-gen-star v0.6.2[0m                          [34m│ [38;5;231mlyft/protoc-gen-star v0.6.2[0m

[36m───────────────[0m[36m┐[0m
[34m362[0m:[38;5;231m replace ( [0m[36m│[0m
[36m───────────────[0m[36m┴[0m[36m──────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[38;5;231m    github.com/oklog/run => github.com/oklog/run v[34m↴[0m  [34m│ [38;5;231m    github.com/oklog/run => github.com/oklog/run v[34m↴[0m
[34m[0m                                             [34m…[38;5;231m1.1.0[0m  [34m│ [0m                                             [34m…[38;5;231m1.1.0[0m
[34m[38;5;231m    github.com/oklog/ulid => github.com/oklog/ulid[34m↴[0m  [34m│ [38;5;231m    github.com/oklog/ulid => github.com/oklog/ulid[34m↴[0m
[34m[0m                                           [34m…[38;5;231m v1.3.1[0m  [34m│ [0m                                           [34m…[38;5;231m v1.3.1[0m
[34m[38;5;231m    github.com/onsi/ginkgo => github.com/onsi/gink[34m↴[0m  [34m│ [38;5;231m    github.com/onsi/ginkgo => github.com/onsi/gink[34m↴[0m
[34m[0m                                        [34m…[38;5;231mgo v1.16.5[0m  [34m│ [0m                                        [34m…[38;5;231mgo v1.16.5[0m
[34m[48;5;52;38;5;231m    github.com/onsi/gomega => github.com/onsi/gomega[34m↴[0m[34m│ [48;5;22;38;5;231m    github.com/onsi/gomega => github.com/onsi/gome[34m↴[0m[48;5;22m[0K[0m
[34m[0m[48;5;52m                                            [34m…[38;5;231m v1.[48;5;124m25[48;5;52m.0[0m[34m│ [0m[48;5;22m                                        [34m…[38;5;231mga v1.[48;5;28m26[48;5;22m.0[0m[48;5;22m[0K[0m
[34m[38;5;231m    github.com/opencontainers/go-digest => github.[34m↵[0m  [34m│ [38;5;231m    github.com/opencontainers/go-digest => github.[34m↵[0m
[34m[38;5;231mcom/opencontainers/go-digest v1.0.0[0m                  [34m│ [38;5;231mcom/opencontainers/go-digest v1.0.0[0m
[34m[38;5;231m    github.com/opencontainers/image-spec => github[34m↵[0m  [34m│ [38;5;231m    github.com/opencontainers/image-spec => github[34m↵[0m
[34m[38;5;231m.com/opencontainers/image-spec v1.0.2[0m                [34m│ [38;5;231m.com/opencontainers/image-spec v1.0.2[0m
[34m[38;5;231m    github.com/opencontainers/runc => github.com/o[34m↵[0m  [34m│ [38;5;231m    github.com/opencontainers/runc => github.com/o[34m↵[0m
[34m[38;5;231mpencontainers/runc v1.1.4[0m                            [34m│ [38;5;231mpencontainers/runc v1.1.4[0m

[36m───────────────[0m[36m┐[0m
[34m437[0m:[38;5;231m replace ( [0m[36m│[0m
[36m───────────────[0m[36m┴[0m[36m──────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[38;5;231m    github.com/tv42/httpunix => github.com/tv42/ht[34m↵[0m  [34m│ [38;5;231m    github.com/tv42/httpunix => github.com/tv42/ht[34m↵[0m
[34m[38;5;231mtpunix v0.0.0-20191220191345-2ba4b9c3382c[0m            [34m│ [38;5;231mtpunix v0.0.0-20191220191345-2ba4b9c3382c[0m
[34m[38;5;231m    github.com/ugorji/go => github.com/ugorji/go v[34m↴[0m  [34m│ [38;5;231m    github.com/ugorji/go => github.com/ugorji/go v[34m↴[0m
[34m[0m                                             [34m…[38;5;231m1.2.8[0m  [34m│ [0m                                             [34m…[38;5;231m1.2.8[0m
[34m[38;5;231m    github.com/ugorji/go/codec => github.com/ugorj[34m↴[0m  [34m│ [38;5;231m    github.com/ugorji/go/codec => github.com/ugorj[34m↴[0m
[34m[0m                                 [34m…[38;5;231mi/go/codec v1.2.8[0m  [34m│ [0m                                 [34m…[38;5;231mi/go/codec v1.2.8[0m
[34m[48;5;52;38;5;231m    github.com/urfave/cli => github.com/urfave/cli v[34m↴[0m[34m│ [48;5;22;38;5;231m    github.com/urfave/cli => github.com/urfave/cli[34m↴[0m[48;5;22m[0K[0m
[34m[0m[48;5;52m                                             [34m…[38;5;231m1.22.[48;5;124m11[0m[34m│ [0m[48;5;22m                                         [34m…[38;5;231m v1.22.[48;5;28m12[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231m    github.com/vdaas/vald-client-go => github.com/vd[34m↵[0m[34m│ [48;5;22;38;5;231m    github.com/vdaas/vald-client-go => github.com/[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231maas/vald-client-go v1.[48;5;124m6[48;5;52m.[48;5;124m3[0m[48;5;52m                            [0m[34m│ [48;5;22;38;5;231mvdaas/vald-client-go v1.[48;5;28m7[48;5;22m.[48;5;28m0[0m[48;5;22m[0K[0m
[34m[38;5;231m    github.com/vishvananda/netlink => github.com/v[34m↵[0m  [34m│ [38;5;231m    github.com/vishvananda/netlink => github.com/v[34m↵[0m
[34m[38;5;231mishvananda/netlink v1.1.0[0m                            [34m│ [38;5;231mishvananda/netlink v1.1.0[0m
[34m[48;5;52;38;5;231m    github.com/vishvananda/netns => github.com/vishv[34m↴[0m[34m│ [48;5;22;38;5;231m    github.com/vishvananda/netns => github.com/vis[34m↵[0m[48;5;22m[0K[0m
[34m[0m[48;5;52m                                 [34m…[38;5;231mananda/netns v0.0.[48;5;124m2[0m[34m│ [48;5;22;38;5;231mhvananda/netns v0.0.[48;5;28m4[0m[48;5;22m[0K[0m
[34m[38;5;231m    github.com/xdg-go/pbkdf2 => github.com/xdg-go/[34m↴[0m  [34m│ [38;5;231m    github.com/xdg-go/pbkdf2 => github.com/xdg-go/[34m↴[0m
[34m[0m                                     [34m…[38;5;231mpbkdf2 v1.0.0[0m  [34m│ [0m                                     [34m…[38;5;231mpbkdf2 v1.0.0[0m
[34m[38;5;231m    github.com/xdg-go/scram => github.com/xdg-go/s[34m↴[0m  [34m│ [38;5;231m    github.com/xdg-go/scram => github.com/xdg-go/s[34m↴[0m
[34m[0m                                       [34m…[38;5;231mcram v1.1.2[0m  [34m│ [0m                                       [34m…[38;5;231mcram v1.1.2[0m
[34m[38;5;231m    github.com/xdg-go/stringprep => github.com/xdg[34m↵[0m  [34m│ [38;5;231m    github.com/xdg-go/stringprep => github.com/xdg[34m↵[0m
[34m[38;5;231m-go/stringprep v1.0.4[0m                                [34m│ [38;5;231m-go/stringprep v1.0.4[0m

[36m───────────────[0m[36m┐[0m
[34m454[0m:[38;5;231m replace ( [0m[36m│[0m
[36m───────────────[0m[36m┴[0m[36m──────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[38;5;231m    github.com/zeebo/assert => github.com/zeebo/as[34m↴[0m  [34m│ [38;5;231m    github.com/zeebo/assert => github.com/zeebo/as[34m↴[0m
[34m[0m                                       [34m…[38;5;231msert v1.3.1[0m  [34m│ [0m                                       [34m…[38;5;231msert v1.3.1[0m
[34m[38;5;231m    github.com/zeebo/xxh3 => github.com/zeebo/xxh3[34m↴[0m  [34m│ [38;5;231m    github.com/zeebo/xxh3 => github.com/zeebo/xxh3[34m↴[0m
[34m[0m                                           [34m…[38;5;231m v1.0.2[0m  [34m│ [0m                                           [34m…[38;5;231m v1.0.2[0m
[34m[38;5;231m    go.etcd.io/bbolt => go.etcd.io/bbolt v1.3.6[0m      [34m│ [38;5;231m    go.etcd.io/bbolt => go.etcd.io/bbolt v1.3.6[0m
[34m[48;5;52;38;5;231m    go.etcd.io/etcd/api/v3 => go.etcd.io/etcd/api/v3[34m↴[0m[34m│ [48;5;22;38;5;231m    go.etcd.io/etcd/api/v3 => go.etcd.io/etcd/api/[34m↴[0m[48;5;22m[0K[0m
[34m[0m[48;5;52m                                             [34m…[38;5;231m v3.5.[48;5;124m6[0m[34m│ [0m[48;5;22m                                         [34m…[38;5;231mv3 v3.5.[48;5;28m7[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231m    go.etcd.io/etcd/client/pkg/v3 => go.etcd.io/etcd[34m↵[0m[34m│ [48;5;22;38;5;231m    go.etcd.io/etcd/client/pkg/v3 => go.etcd.io/et[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231m/client/pkg/v3 v3.5.[48;5;124m6[0m[48;5;52m                                [0m[34m│ [48;5;22;38;5;231mcd/client/pkg/v3 v3.5.[48;5;28m7[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231m    go.etcd.io/etcd/client/v2 => go.etcd.io/etcd/cli[34m↴[0m[34m│ [48;5;22;38;5;231m    go.etcd.io/etcd/client/v2 => go.etcd.io/etcd/c[34m↴[0m[48;5;22m[0K[0m
[34m[0m[48;5;52m                                     [34m…[38;5;231ment/v2 v2.305.[48;5;124m6[0m[34m│ [0m[48;5;22m                                 [34m…[38;5;231mlient/v2 v2.305.[48;5;28m7[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231m    go.etcd.io/etcd/client/v3 => go.etcd.io/etcd/cli[34m↴[0m[34m│ [48;5;22;38;5;231m    go.etcd.io/etcd/client/v3 => go.etcd.io/etcd/c[34m↴[0m[48;5;22m[0K[0m
[34m[0m[48;5;52m                                       [34m…[38;5;231ment/v3 v3.5.[48;5;124m6[0m[34m│ [0m[48;5;22m                                   [34m…[38;5;231mlient/v3 v3.5.[48;5;28m7[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231m    go.etcd.io/etcd/pkg/v3 => go.etcd.io/etcd/pkg/v3[34m↴[0m[34m│ [48;5;22;38;5;231m    go.etcd.io/etcd/pkg/v3 => go.etcd.io/etcd/pkg/[34m↴[0m[48;5;22m[0K[0m
[34m[0m[48;5;52m                                             [34m…[38;5;231m v3.5.[48;5;124m6[0m[34m│ [0m[48;5;22m                                         [34m…[38;5;231mv3 v3.5.[48;5;28m7[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231m    go.etcd.io/etcd/raft/v3 => go.etcd.io/etcd/raft/[34m↴[0m[34m│ [48;5;22;38;5;231m    go.etcd.io/etcd/raft/v3 => go.etcd.io/etcd/raf[34m↴[0m[48;5;22m[0K[0m
[34m[0m[48;5;52m                                           [34m…[38;5;231mv3 v3.5.[48;5;124m6[0m[34m│ [0m[48;5;22m                                       [34m…[38;5;231mt/v3 v3.5.[48;5;28m7[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231m    go.etcd.io/etcd/server/v3 => go.etcd.io/etcd/ser[34m↴[0m[34m│ [48;5;22;38;5;231m    go.etcd.io/etcd/server/v3 => go.etcd.io/etcd/s[34m↴[0m[48;5;22m[0K[0m
[34m[0m[48;5;52m                                       [34m…[38;5;231mver/v3 v3.5.[48;5;124m6[0m[34m│ [0m[48;5;22m                                   [34m…[38;5;231merver/v3 v3.5.[48;5;28m7[0m[48;5;22m[0K[0m
[34m[38;5;231m    go.mongodb.org/mongo-driver => go.mongodb.org/[34m↵[0m  [34m│ [38;5;231m    go.mongodb.org/mongo-driver => go.mongodb.org/[34m↵[0m
[34m[38;5;231mmongo-driver v1.11.1[0m                                 [34m│ [38;5;231mmongo-driver v1.11.1[0m
[34m[38;5;231m    go.mozilla.org/pkcs7 => go.mozilla.org/pkcs7 v[34m↵[0m  [34m│ [38;5;231m    go.mozilla.org/pkcs7 => go.mozilla.org/pkcs7 v[34m↵[0m
[34m[38;5;231m0.0.0-20210826202110-33d05740a352[0m                    [34m│ [38;5;231m0.0.0-20210826202110-33d05740a352[0m
[34m[38;5;231m    go.opencensus.io => go.opencensus.io v0.24.0[0m     [34m│ [38;5;231m    go.opencensus.io => go.opencensus.io v0.24.0[0m

[36m───────────────[0m[36m┐[0m
[34m476[0m:[38;5;231m replace ( [0m[36m│[0m
[36m───────────────[0m[36m┴[0m[36m──────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[38;5;231m    go.opentelemetry.io/otel/sdk/metric => go.open[34m↵[0m  [34m│ [38;5;231m    go.opentelemetry.io/otel/sdk/metric => go.open[34m↵[0m
[34m[38;5;231mtelemetry.io/otel/sdk/metric v0.33.0[0m                 [34m│ [38;5;231mtelemetry.io/otel/sdk/metric v0.33.0[0m
[34m[38;5;231m    go.opentelemetry.io/otel/trace => go.opentelem[34m↵[0m  [34m│ [38;5;231m    go.opentelemetry.io/otel/trace => go.opentelem[34m↵[0m
[34m[38;5;231metry.io/otel/trace v1.11.1[0m                           [34m│ [38;5;231metry.io/otel/trace v1.11.1[0m
[34m[38;5;231m    go.opentelemetry.io/proto/otlp => go.opentelem[34m↵[0m  [34m│ [38;5;231m    go.opentelemetry.io/proto/otlp => go.opentelem[34m↵[0m
[34m[38;5;231metry.io/proto/otlp v0.19.0[0m                           [34m│ [38;5;231metry.io/proto/otlp v0.19.0[0m
[34m[48;5;52;38;5;231m    go.starlark.net => go.starlark.net v0.0.0-[48;5;124m202301[48;5;52;34m↵[0m[34m│ [48;5;22;38;5;231m    go.starlark.net => go.starlark.net v0.0.0-[48;5;28m2023[48;5;22;34m↵[0m[48;5;22m[0K[0m
[34m[48;5;124;38;5;231m12144946[48;5;52m-[48;5;124mfae38c8a6d89[0m[48;5;52m                                [0m[34m│ [48;5;28;38;5;231m0122040757[48;5;22m-[48;5;28m066229b0515d[0m[48;5;22m[0K[0m
[34m[38;5;231m    go.uber.org/atomic => go.uber.org/atomic v1.10[34m↴[0m  [34m│ [38;5;231m    go.uber.org/atomic => go.uber.org/atomic v1.10[34m↴[0m
[34m[0m                                                [34m…[38;5;231m.0[0m  [34m│ [0m                                                [34m…[38;5;231m.0[0m
[34m[38;5;231m    go.uber.org/automaxprocs => go.uber.org/automa[34m↴[0m  [34m│ [38;5;231m    go.uber.org/automaxprocs => go.uber.org/automa[34m↴[0m
[34m[0m                                     [34m…[38;5;231mxprocs v1.5.1[0m  [34m│ [0m                                     [34m…[38;5;231mxprocs v1.5.1[0m
[34m[38;5;231m    go.uber.org/goleak => go.uber.org/goleak v1.2.0[0m  [34m│ [38;5;231m    go.uber.org/goleak => go.uber.org/goleak v1.2.0[0m

[36m───────────────[0m[36m┐[0m
[34m486[0m:[38;5;231m replace ( [0m[36m│[0m
[36m───────────────[0m[36m┴[0m[36m──────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[38;5;231m    go4.org/unsafe/assume-no-moving-gc => go4.org/[34m↵[0m  [34m│ [38;5;231m    go4.org/unsafe/assume-no-moving-gc => go4.org/[34m↵[0m
[34m[38;5;231munsafe/assume-no-moving-gc v0.0.0-20220617031537-9[34m↵[0m  [34m│ [38;5;231munsafe/assume-no-moving-gc v0.0.0-20220617031537-9[34m↵[0m
[34m[38;5;231m28513b29760[0m                                          [34m│ [38;5;231m28513b29760[0m
[34m[38;5;231m    gocloud.dev => gocloud.dev v0.28.0[0m               [34m│ [38;5;231m    gocloud.dev => gocloud.dev v0.28.0[0m
[34m[38;5;231m    golang.org/x/crypto => golang.org/x/crypto v0.[34m↴[0m  [34m│ [38;5;231m    golang.org/x/crypto => golang.org/x/crypto v0.[34m↴[0m
[34m[0m                                               [34m…[38;5;231m5.0[0m  [34m│ [0m                                               [34m…[38;5;231m5.0[0m
[34m[48;5;52;38;5;231m    golang.org/x/exp => golang.org/x/exp v0.0.0-[48;5;124m2023[48;5;52;34m↵[0m[34m│ [48;5;22;38;5;231m    golang.org/x/exp => golang.org/x/exp v0.0.0-[48;5;28m20[48;5;22;34m↵[0m[48;5;22m[0K[0m
[34m[48;5;124;38;5;231m0116083435[48;5;52m-[48;5;124m1de6713980de[0m[48;5;52m                              [0m[34m│ [48;5;28;38;5;231m230118134722[48;5;22m-[48;5;28ma68e582fa157[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231m    golang.org/x/exp/typeparams => golang.org/x/exp/[34m↵[0m[34m│ [48;5;22;38;5;231m    golang.org/x/exp/typeparams => golang.org/x/ex[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mtypeparams v0.0.0-[48;5;124m20230116083435[48;5;52m-[48;5;124m1de6713980de[0m[48;5;52m        [0m[34m│ [48;5;22;38;5;231mp/typeparams v0.0.0-[48;5;28m20230118134722[48;5;22m-[48;5;28ma68e582fa157[0m[48;5;22m[0K[0m
[34m[38;5;231m    golang.org/x/image => golang.org/x/image v0.3.0[0m  [34m│ [38;5;231m    golang.org/x/image => golang.org/x/image v0.3.0[0m
[34m[38;5;231m    golang.org/x/lint => golang.org/x/lint v0.0.0-[34m↵[0m  [34m│ [38;5;231m    golang.org/x/lint => golang.org/x/lint v0.0.0-[34m↵[0m
[34m[38;5;231m20210508222113-6edffad5e616[0m                          [34m│ [38;5;231m20210508222113-6edffad5e616[0m
[34m[38;5;231m    golang.org/x/mobile => golang.org/x/mobile v0.[34m↵[0m  [34m│ [38;5;231m    golang.org/x/mobile => golang.org/x/mobile v0.[34m↵[0m
[34m[38;5;231m0.0-20221110043201-43a038452099[0m                      [34m│ [38;5;231m0.0-20221110043201-43a038452099[0m

[36m───────────────[0m[36m┐[0m
[34m505[0m:[38;5;231m replace ( [0m[36m│[0m
[36m───────────────[0m[36m┴[0m[36m──────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[38;5;231m    gonum.org/v1/gonum => gonum.org/v1/gonum v0.12[34m↴[0m  [34m│ [38;5;231m    gonum.org/v1/gonum => gonum.org/v1/gonum v0.12[34m↴[0m
[34m[0m                                                [34m…[38;5;231m.0[0m  [34m│ [0m                                                [34m…[38;5;231m.0[0m
[34m[38;5;231m    gonum.org/v1/hdf5 => gonum.org/v1/hdf5 v0.0.0-[34m↵[0m  [34m│ [38;5;231m    gonum.org/v1/hdf5 => gonum.org/v1/hdf5 v0.0.0-[34m↵[0m
[34m[38;5;231m20210714002203-8c5d23bc6946[0m                          [34m│ [38;5;231m20210714002203-8c5d23bc6946[0m
[34m[38;5;231m    gonum.org/v1/plot => gonum.org/v1/plot v0.12.0[0m   [34m│ [38;5;231m    gonum.org/v1/plot => gonum.org/v1/plot v0.12.0[0m
[34m[48;5;52;38;5;231m    google.golang.org/api => google.golang.org/api v[34m↴[0m[34m│ [48;5;22;38;5;231m    google.golang.org/api => google.golang.org/api[34m↴[0m[48;5;22m[0K[0m
[34m[0m[48;5;52m                                             [34m…[38;5;231m0.[48;5;124m107[48;5;52m.0[0m[34m│ [0m[48;5;22m                                         [34m…[38;5;231m v0.[48;5;28m108[48;5;22m.0[0m[48;5;22m[0K[0m
[34m[38;5;231m    google.golang.org/appengine => google.golang.o[34m↵[0m  [34m│ [38;5;231m    google.golang.org/appengine => google.golang.o[34m↵[0m
[34m[38;5;231mrg/appengine v1.6.7[0m                                  [34m│ [38;5;231mrg/appengine v1.6.7[0m
[34m[48;5;52;38;5;231m    google.golang.org/genproto => google.golang.org/[34m↵[0m[34m│ [48;5;22;38;5;231m    google.golang.org/genproto => google.golang.or[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mgenproto v0.0.0-[48;5;124m20230117162540[48;5;52m-[48;5;124m28d6b9783ac4[0m[48;5;52m          [0m[34m│ [48;5;22;38;5;231mg/genproto v0.0.0-[48;5;28m20230123190316[48;5;22m-[48;5;28m2c411cf9d197[0m[48;5;22m[0K[0m
[34m[38;5;231m    google.golang.org/grpc => google.golang.org/gr[34m↴[0m  [34m│ [38;5;231m    google.golang.org/grpc => google.golang.org/gr[34m↴[0m
[34m[0m                                        [34m…[38;5;231mpc v1.52.0[0m  [34m│ [0m                                        [34m…[38;5;231mpc v1.52.0[0m
[34m[38;5;231m    google.golang.org/grpc/cmd/protoc-gen-go-grpc [34m↵[0m  [34m│ [38;5;231m    google.golang.org/grpc/cmd/protoc-gen-go-grpc [34m↵[0m
[34m[38;5;231m=> google.golang.org/grpc/cmd/protoc-gen-go-grpc v[34m↵[0m  [34m│ [38;5;231m=> google.golang.org/grpc/cmd/protoc-gen-go-grpc v[34m↵[0m
[34m[38;5;231m1.2.0[0m                                                [34m│ [38;5;231m1.2.0[0m
[34m[38;5;231m    google.golang.org/protobuf => google.golang.or[34m↴[0m  [34m│ [38;5;231m    google.golang.org/protobuf => google.golang.or[34m↴[0m
[34m[0m                                [34m…[38;5;231mg/protobuf v1.28.1[0m  [34m│ [0m                                [34m…[38;5;231mg/protobuf v1.28.1[0m
[34m[38;5;231m    gopkg.in/alecthomas/kingpin.v2 => gopkg.in/ale[34m↵[0m  [34m│ [38;5;231m    gopkg.in/alecthomas/kingpin.v2 => gopkg.in/ale[34m↵[0m
[34m[38;5;231mcthomas/kingpin.v2 v2.2.6[0m                            [34m│ [38;5;231mcthomas/kingpin.v2 v2.2.6[0m
[34m[38;5;231m    gopkg.in/check.v1 => gopkg.in/check.v1 v1.0.0-[34m↵[0m  [34m│ [38;5;231m    gopkg.in/check.v1 => gopkg.in/check.v1 v1.0.0-[34m↵[0m
[34m[38;5;231m20201130134442-10cb98267c6c[0m                          [34m│ [38;5;231m20201130134442-10cb98267c6c[0m
[34m[38;5;231m    gopkg.in/gcfg.v1 => gopkg.in/gcfg.v1 v1.2.3[0m      [34m│ [38;5;231m    gopkg.in/gcfg.v1 => gopkg.in/gcfg.v1 v1.2.3[0m
[34m[48;5;52;38;5;231m    gopkg.in/inconshreveable/log15.v2 => gopkg.in/in[34m↵[0m[34m│ [48;5;22;38;5;231m    gopkg.in/inconshreveable/log15.v2 => gopkg.in/[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mconshreveable/log15.v2 v2.[48;5;124m0[48;5;52m.0[48;5;124m-20221122034931-5555550[48;5;52;34m↵[0m[34m│ [48;5;22;38;5;231minconshreveable/log15.v2 v2.[48;5;28m16[48;5;22m.0[0m[48;5;22m[0K[0m
[34m[48;5;124;38;5;231m54819[0m[48;5;52m                                                [0m[34m│ [0m
[34m[38;5;231m    gopkg.in/inf.v0 => gopkg.in/inf.v0 v0.9.1[0m        [34m│ [38;5;231m    gopkg.in/inf.v0 => gopkg.in/inf.v0 v0.9.1[0m
[34m[38;5;231m    gopkg.in/ini.v1 => gopkg.in/ini.v1 v1.67.0[0m       [34m│ [38;5;231m    gopkg.in/ini.v1 => gopkg.in/ini.v1 v1.67.0[0m
[34m[38;5;231m    gopkg.in/natefinch/lumberjack.v2 => gopkg.in/n[34m↵[0m  [34m│ [38;5;231m    gopkg.in/natefinch/lumberjack.v2 => gopkg.in/n[34m↵[0m
[34m[38;5;231matefinch/lumberjack.v2 v2.0.0[0m                        [34m│ [38;5;231matefinch/lumberjack.v2 v2.0.0[0m

[36m───────────────[0m[36m┐[0m
[34m537[0m:[38;5;231m replace ( [0m[36m│[0m
[36m───────────────[0m[36m┴[0m[36m──────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[38;5;231m    k8s.io/cri-api => k8s.io/cri-api v0.26.0[0m         [34m│ [38;5;231m    k8s.io/cri-api => k8s.io/cri-api v0.26.0[0m
[34m[38;5;231m    k8s.io/gengo => k8s.io/gengo v0.0.0-2022101119[34m↴[0m  [34m│ [38;5;231m    k8s.io/gengo => k8s.io/gengo v0.0.0-2022101119[34m↴[0m
[34m[0m                                 [34m…[38;5;231m3443-fad74ee6edd9[0m  [34m│ [0m                                 [34m…[38;5;231m3443-fad74ee6edd9[0m
[34m[38;5;231m    k8s.io/klog => k8s.io/klog v1.0.0[0m                [34m│ [38;5;231m    k8s.io/klog => k8s.io/klog v1.0.0[0m
[34m[48;5;52;38;5;231m    k8s.io/klog/v2 => k8s.io/klog/v2 v2.[48;5;124m80[48;5;52m.[48;5;124m1[0m[48;5;52m         [0m[34m│ [48;5;22;38;5;231m    k8s.io/klog/v2 => k8s.io/klog/v2 v2.[48;5;28m90[48;5;22m.[48;5;28m0[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231m    k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.[34m↵[0m[34m│ [48;5;22;38;5;231m    k8s.io/kube-openapi => k8s.io/kube-openapi v0.[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231m0-[48;5;124m20230117224833[48;5;52m-[48;5;124m444ee995c120[0m[48;5;52m                        [0m[34m│ [48;5;22;38;5;231m0.0-[48;5;28m20230123231816[48;5;22m-[48;5;28m1cb3ae25d79a[0m[48;5;22m[0K[0m
[34m[38;5;231m    k8s.io/kubernetes => k8s.io/kubernetes v0.26.0[0m   [34m│ [38;5;231m    k8s.io/kubernetes => k8s.io/kubernetes v0.26.0[0m
[34m[38;5;231m    k8s.io/metrics => k8s.io/metrics v0.26.0[0m         [34m│ [38;5;231m    k8s.io/metrics => k8s.io/metrics v0.26.0[0m
[34m[38;5;231m    nhooyr.io/websocket => nhooyr.io/websocket v1.[34m↴[0m  [34m│ [38;5;231m    nhooyr.io/websocket => nhooyr.io/websocket v1.[34m↴[0m
[34m[0m                                               [34m…[38;5;231m8.7[0m  [34m│ [0m                                               [34m…[38;5;231m8.7[0m
[34m[38;5;231m    rsc.io/pdf => rsc.io/pdf v0.1.1[0m                  [34m│ [38;5;231m    rsc.io/pdf => rsc.io/pdf v0.1.1[0m
[34m[48;5;52;38;5;231m    sigs.k8s.io/apiserver-network-proxy/konnectivity[34m↵[0m[34m│ [48;5;22;38;5;231m    sigs.k8s.io/apiserver-network-proxy/konnectivi[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231m-client => sigs.k8s.io/apiserver-network-proxy/konne[34m↵[0m[34m│ [48;5;22;38;5;231mty-client => sigs.k8s.io/apiserver-network-proxy/k[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mctivity-client v0.1.[48;5;124m0[0m[48;5;52m                                [0m[34m│ [48;5;22;38;5;231monnectivity-client v0.1.[48;5;28m1[0m[48;5;22m[0K[0m
[34m[38;5;231m    sigs.k8s.io/controller-runtime => sigs.k8s.io/[34m↵[0m  [34m│ [38;5;231m    sigs.k8s.io/controller-runtime => sigs.k8s.io/[34m↵[0m
[34m[38;5;231mcontroller-runtime v0.14.1[0m                           [34m│ [38;5;231mcontroller-runtime v0.14.1[0m
[34m[38;5;231m    sigs.k8s.io/json => sigs.k8s.io/json v0.0.0-20[34m↵[0m  [34m│ [38;5;231m    sigs.k8s.io/json => sigs.k8s.io/json v0.0.0-20[34m↵[0m
[34m[38;5;231m221116044647-bc3834ca7abd[0m                            [34m│ [38;5;231m221116044647-bc3834ca7abd[0m
[34m[38;5;231m    sigs.k8s.io/kustomize => sigs.k8s.io/kustomize[34m↵[0m  [34m│ [38;5;231m    sigs.k8s.io/kustomize => sigs.k8s.io/kustomize[34m↵[0m
[34m[38;5;231m v2.0.3+incompatible[0m                                 [34m│ [38;5;231m v2.0.3+incompatible[0m

[36m───────────────[0m[36m┐[0m
[34m600[0m:[38;5;231m require ( [0m[36m│[0m
[36m───────────────[0m[36m┴[0m[36m──────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[38;5;231m    golang.org/x/tools v0.4.0[0m                        [34m│ [38;5;231m    golang.org/x/tools v0.4.0[0m
[34m[38;5;231m    gonum.org/v1/hdf5 v0.0.0-00010101000000-000000[34m↴[0m  [34m│ [38;5;231m    gonum.org/v1/hdf5 v0.0.0-00010101000000-000000[34m↴[0m
[34m[0m                                            [34m…[38;5;231m000000[0m  [34m│ [0m                                            [34m…[38;5;231m000000[0m
[34m[38;5;231m    gonum.org/v1/plot v0.0.0-00010101000000-000000[34m↴[0m  [34m│ [38;5;231m    gonum.org/v1/plot v0.0.0-00010101000000-000000[34m↴[0m
[34m[0m                                            [34m…[38;5;231m000000[0m  [34m│ [0m                                            [34m…[38;5;231m000000[0m
[34m[48;5;52;38;5;231m    google.golang.org/genproto v0.0.0-[48;5;124m20230109162033[48;5;52;34m↴[0m[34m│ [48;5;22;38;5;231m    google.golang.org/genproto v0.0.0-[48;5;28m202301131545[48;5;22;34m↴[0m[48;5;22m[0K[0m
[34m[0m[48;5;52m                                       [34m…[38;5;231m-[48;5;124m3c3c17ce83e6[0m[34m│ [0m[48;5;22m                                   [34m…[48;5;28;38;5;231m10[48;5;22m-[48;5;28mdbe35b8444a5[0m[48;5;22m[0K[0m
[34m[38;5;231m    google.golang.org/grpc v1.51.0[0m                   [34m│ [38;5;231m    google.golang.org/grpc v1.51.0[0m
[34m[38;5;231m    google.golang.org/protobuf v1.28.1[0m               [34m│ [38;5;231m    google.golang.org/protobuf v1.28.1[0m
[34m[38;5;231m    gopkg.in/yaml.v2 v2.4.0[0m                          [34m│ [38;5;231m    gopkg.in/yaml.v2 v2.4.0[0m

[36m───────────────[0m[36m┐[0m
[34m696[0m:[38;5;231m require ( [0m[36m│[0m
[36m───────────────[0m[36m┴[0m[36m──────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[38;5;231m    golang.org/x/time v0.3.0 // indirect[0m             [34m│ [38;5;231m    golang.org/x/time v0.3.0 // indirect[0m
[34m[38;5;231m    golang.org/x/xerrors v0.0.0-20220907171357-04b[34m↵[0m  [34m│ [38;5;231m    golang.org/x/xerrors v0.0.0-20220907171357-04b[34m↵[0m
[34m[38;5;231me3eba64a2 // indirect[0m                                [34m│ [38;5;231me3eba64a2 // indirect[0m
[34m[38;5;231m    gomodules.xyz/jsonpatch/v2 v2.2.0 // indirect[0m    [34m│ [38;5;231m    gomodules.xyz/jsonpatch/v2 v2.2.0 // indirect[0m
[34m[48;5;52;38;5;231m    google.golang.org/api v0.[48;5;124m106[48;5;52m.0 // indirect[0m[48;5;52m       [0m[34m│ [48;5;22;38;5;231m    google.golang.org/api v0.[48;5;28m107[48;5;22m.0 // indirect[0m[48;5;22m[0K[0m
[34m[38;5;231m    google.golang.org/appengine v1.6.7 // indirect[0m   [34m│ [38;5;231m    google.golang.org/appengine v1.6.7 // indirect[0m
[34m[38;5;231m    gopkg.in/inf.v0 v0.9.1 // indirect[0m               [34m│ [38;5;231m    gopkg.in/inf.v0 v0.9.1 // indirect[0m
[34m[38;5;231m    gopkg.in/yaml.v3 v3.0.1 // indirect[0m              [34m│ [38;5;231m    gopkg.in/yaml.v3 v3.0.1 // indirect[0m

[1;4;33mgo.sum[0m

[36m───[0m[36m┐[0m
[34m1[0m: [36m│[0m
[36m───[0m[36m┴[0m[36m──────────────────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[48;5;52;38;5;231mcloud.google.com/go v0.[48;5;124m108[48;5;52m.0 h1:[48;5;124mxntQwnfn8oHGX0crLVin[48;5;52;34m↵[0m[34m│ [48;5;22;38;5;231mcloud.google.com/go v0.[48;5;28m109[48;5;22m.0 h1:[48;5;28m38CZoKGlCnPZjGdyj0[48;5;22;34m↵[0m[48;5;22m[0K[0m
[34m[48;5;124;38;5;231mvHM+AhXvi3QHQIEcX[48;5;52m/[48;5;124m2hiWk[48;5;52m=[0m[48;5;52m                             [0m[34m│ [48;5;28;38;5;231mZfpoGae0[48;5;22m/[48;5;28mwgNfy5F0byyxg0Gk[48;5;22m=[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mcloud.google.com/go v0.[48;5;124m108[48;5;52m.0/go.mod h1:[48;5;124mlNUfQqusBJp0b[48;5;52;34m↵[0m[34m│ [48;5;22;38;5;231mcloud.google.com/go v0.[48;5;28m109[48;5;22m.0/go.mod h1:[48;5;28m2sYycXt75t/[48;5;22;34m↵[0m[48;5;22m[0K[0m
[34m[48;5;124;38;5;231mgAg6qrHgYFYbTB+dOiob1itwnlD33Q[48;5;52m=[0m[48;5;52m                      [0m[34m│ [48;5;28;38;5;231mCSB5R9M2wPU1tJmire7AQZTPtITcGBVE[48;5;22m=[0m[48;5;22m[0K[0m
[34m[38;5;231mcloud.google.com/go/accessapproval v1.5.0/go.mod h[34m↵[0m  [34m│ [38;5;231mcloud.google.com/go/accessapproval v1.5.0/go.mod h[34m↵[0m
[34m[38;5;231m1:HFy3tuiGvMdcd/u+Cu5b9NkO1pEICJ46IR82PoUdplw=[0m       [34m│ [38;5;231m1:HFy3tuiGvMdcd/u+Cu5b9NkO1pEICJ46IR82PoUdplw=[0m
[34m[38;5;231mcloud.google.com/go/accesscontextmanager v1.3.0/go[34m↵[0m  [34m│ [38;5;231mcloud.google.com/go/accesscontextmanager v1.3.0/go[34m↵[0m
[34m[38;5;231m.mod h1:TgCBehyr5gNMz7ZaH9xubp+CE8dkrszb4oK9CWyvD4[34m↵[0m  [34m│ [38;5;231m.mod h1:TgCBehyr5gNMz7ZaH9xubp+CE8dkrszb4oK9CWyvD4[34m↵[0m
[34m[38;5;231mo=[0m                                                   [34m│ [38;5;231mo=[0m
[34m[38;5;231mcloud.google.com/go/accesscontextmanager v1.4.0/go[34m↵[0m  [34m│ [38;5;231mcloud.google.com/go/accesscontextmanager v1.4.0/go[34m↵[0m
[34m[38;5;231m.mod h1:/Kjh7BBu/Gh83sv+K60vN9QE5NJcd80sU33vIe2IFP[34m↵[0m  [34m│ [38;5;231m.mod h1:/Kjh7BBu/Gh83sv+K60vN9QE5NJcd80sU33vIe2IFP[34m↵[0m
[34m[38;5;231mE=[0m                                                   [34m│ [38;5;231mE=[0m

[36m─────────────────────────────────────────────────────────────────────────────────────[0m[36m┐[0m
[34m81[0m:[38;5;231m cloud.google.com/go/maps v0.1.0/go.mod h1:BQM97WGyfw9FWEmQMpZ5T6cpovXXSd1cGmFma9 [0m[36m│[0m
[36m─────────────────────────────────────────────────────────────────────────────────────[0m[36m┴[0m[36m────────────────────[0m
[34m[38;5;231mcloud.google.com/go/mediatranslation v0.6.0/go.mod[34m↵[0m  [34m│ [38;5;231mcloud.google.com/go/mediatranslation v0.6.0/go.mod[34m↵[0m
[34m[38;5;231m h1:hHdBCTYNigsBxshbznuIMFNe5QXEowAuNmmC7h8pu5w=[0m     [34m│ [38;5;231m h1:hHdBCTYNigsBxshbznuIMFNe5QXEowAuNmmC7h8pu5w=[0m
[34m[38;5;231mcloud.google.com/go/memcache v1.7.0/go.mod h1:ywMK[34m↵[0m  [34m│ [38;5;231mcloud.google.com/go/memcache v1.7.0/go.mod h1:ywMK[34m↵[0m
[34m[38;5;231mfjWhNtkQTxrWxCkCFkoPjLHPW6A7WOTVI8xy3LY=[0m             [34m│ [38;5;231mfjWhNtkQTxrWxCkCFkoPjLHPW6A7WOTVI8xy3LY=[0m
[34m[38;5;231mcloud.google.com/go/metastore v1.8.0/go.mod h1:zHi[34m↵[0m  [34m│ [38;5;231mcloud.google.com/go/metastore v1.8.0/go.mod h1:zHi[34m↵[0m
[34m[38;5;231mMc4ZUpBiM7twCIFQmJ9JMEkDSyZS9U12uf7wHqSI=[0m            [34m│ [38;5;231mMc4ZUpBiM7twCIFQmJ9JMEkDSyZS9U12uf7wHqSI=[0m
[34m[48;5;52;38;5;231mcloud.google.com/go/monitoring v1.[48;5;124m11[48;5;52m.0/go.mod h1:[48;5;124mR40[48;5;52;34m↵[0m[34m│ [48;5;22;38;5;231mcloud.google.com/go/monitoring v1.[48;5;28m12[48;5;22m.0/go.mod h1:[48;5;28my[48;5;22;34m↵[0m[48;5;22m[0K[0m
[34m[48;5;124;38;5;231mQ1vMzEVF76pD2O6s9[48;5;52m/[48;5;124mUtvQfrJrLSMXWLByUdEWwU[48;5;52m=[0m[48;5;52m            [0m[34m│ [48;5;28;38;5;231mx8Jj2fZNEkL[48;5;22m/[48;5;28mGYZyTLS4ZtZEZN8WtDEiEqG4kLK50w[48;5;22m=[0m[48;5;22m[0K[0m
[34m[38;5;231mcloud.google.com/go/networkconnectivity v1.7.0/go.[34m↵[0m  [34m│ [38;5;231mcloud.google.com/go/networkconnectivity v1.7.0/go.[34m↵[0m
[34m[38;5;231mmod h1:RMuSbkdbPwNMQjB5HBWD5MpTBnNm39iAVpC3TmsExt8=[0m  [34m│ [38;5;231mmod h1:RMuSbkdbPwNMQjB5HBWD5MpTBnNm39iAVpC3TmsExt8=[0m
[34m[38;5;231mcloud.google.com/go/networkmanagement v1.5.0/go.mo[34m↵[0m  [34m│ [38;5;231mcloud.google.com/go/networkmanagement v1.5.0/go.mo[34m↵[0m
[34m[38;5;231md h1:ZnOeZ/evzUdUsnvRt792H0uYEnHQEMaz+REhhzJRcf4=[0m    [34m│ [38;5;231md h1:ZnOeZ/evzUdUsnvRt792H0uYEnHQEMaz+REhhzJRcf4=[0m
[34m[38;5;231mcloud.google.com/go/networksecurity v0.6.0/go.mod [34m↵[0m  [34m│ [38;5;231mcloud.google.com/go/networksecurity v0.6.0/go.mod [34m↵[0m
[34m[38;5;231mh1:Q5fjhTr9WMI5mbpRYEbiexTzROf7ZbDzvzCrNl14nyU=[0m      [34m│ [38;5;231mh1:Q5fjhTr9WMI5mbpRYEbiexTzROf7ZbDzvzCrNl14nyU=[0m

[36m──────────────────────────────────────────────────────────────────────────────────────[0m[36m┐[0m
[34m117[0m:[38;5;231m cloud.google.com/go/serviceusage v1.4.0/go.mod h1:SB4yxXSaYVuUBYUml6qklyONXNLt83 [0m[36m│[0m
[36m──────────────────────────────────────────────────────────────────────────────────────[0m[36m┴[0m[36m───────────────────[0m
[34m[38;5;231mcloud.google.com/go/shell v1.4.0/go.mod h1:HDxPzZf[34m↵[0m  [34m│ [38;5;231mcloud.google.com/go/shell v1.4.0/go.mod h1:HDxPzZf[34m↵[0m
[34m[38;5;231m3GkDdhExzD/gs8Grqk+dmYcEjGShZgYa9URw=[0m                [34m│ [38;5;231m3GkDdhExzD/gs8Grqk+dmYcEjGShZgYa9URw=[0m
[34m[38;5;231mcloud.google.com/go/spanner v1.41.0/go.mod h1:MLYD[34m↵[0m  [34m│ [38;5;231mcloud.google.com/go/spanner v1.41.0/go.mod h1:MLYD[34m↵[0m
[34m[38;5;231mBJR/dY4Wt7ZaMIQ7rXOTLjYrmxLE/5ve9vFfWos=[0m             [34m│ [38;5;231mBJR/dY4Wt7ZaMIQ7rXOTLjYrmxLE/5ve9vFfWos=[0m
[34m[38;5;231mcloud.google.com/go/speech v1.9.0/go.mod h1:xQ0jTc[34m↵[0m  [34m│ [38;5;231mcloud.google.com/go/speech v1.9.0/go.mod h1:xQ0jTc[34m↵[0m
[34m[38;5;231mmnRFFM2RfX/U+rk6FQNUF6DQlydUSyoooSpco=[0m               [34m│ [38;5;231mmnRFFM2RfX/U+rk6FQNUF6DQlydUSyoooSpco=[0m
[34m[48;5;52;38;5;231mcloud.google.com/go/storage v1.[48;5;124m28[48;5;52m.[48;5;124m1[48;5;52m h1:[48;5;124mF5QDG5ChchaAV[48;5;52;34m↵[0m[34m│ [48;5;22;38;5;231mcloud.google.com/go/storage v1.[48;5;28m29[48;5;22m.[48;5;28m0[48;5;22m h1:[48;5;28m6weCgzRvMg7[48;5;22;34m↵[0m[48;5;22m[0K[0m
[34m[48;5;124;38;5;231mQhINh24U99OWHURqrW8OmQcGKXcbgI[48;5;52m=[0m[48;5;52m                      [0m[34m│ [48;5;28;38;5;231mlzuUurI4697AqIRPU1SvzHhynwpW31jI[48;5;22m=[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mcloud.google.com/go/storage v1.[48;5;124m28[48;5;52m.[48;5;124m1[48;5;52m/go.mod h1:[48;5;124mQnisd4[48;5;52;34m↵[0m[34m│ [48;5;22;38;5;231mcloud.google.com/go/storage v1.[48;5;28m29[48;5;22m.[48;5;28m0[48;5;22m/go.mod h1:[48;5;28m4puE[48;5;22;34m↵[0m[48;5;22m[0K[0m
[34m[48;5;124;38;5;231mCqDdo6BGs2AD5LLnEsmSQ80wQ5ogcBBKhU86Y[48;5;52m=[0m[48;5;52m               [0m[34m│ [48;5;28;38;5;231mjyTKnku6gfKoTfNOU/W+a9JyuVNxjpS5GBrB8h4[48;5;22m=[0m[48;5;22m[0K[0m
[34m[38;5;231mcloud.google.com/go/storagetransfer v1.6.0/go.mod [34m↵[0m  [34m│ [38;5;231mcloud.google.com/go/storagetransfer v1.6.0/go.mod [34m↵[0m
[34m[38;5;231mh1:y77xm4CQV/ZhFZH75PLEXY0ROiS7Gh6pSKrM8dJyg6I=[0m      [34m│ [38;5;231mh1:y77xm4CQV/ZhFZH75PLEXY0ROiS7Gh6pSKrM8dJyg6I=[0m
[34m[38;5;231mcloud.google.com/go/talent v1.4.0/go.mod h1:ezFtAg[34m↵[0m  [34m│ [38;5;231mcloud.google.com/go/talent v1.4.0/go.mod h1:ezFtAg[34m↵[0m
[34m[38;5;231mVuRf8jRsvyE6EwmbTK5LKciD4KVnHuDEFmOOA=[0m               [34m│ [38;5;231mVuRf8jRsvyE6EwmbTK5LKciD4KVnHuDEFmOOA=[0m
[34m[38;5;231mcloud.google.com/go/texttospeech v1.5.0/go.mod h1:[34m↵[0m  [34m│ [38;5;231mcloud.google.com/go/texttospeech v1.5.0/go.mod h1:[34m↵[0m
[34m[38;5;231moKPLhR4n4ZdQqWKURdwxMy0uiTS1xU161C8W57Wkea4=[0m         [34m│ [38;5;231moKPLhR4n4ZdQqWKURdwxMy0uiTS1xU161C8W57Wkea4=[0m

[36m────────────────────────────────────────────────────────────────────────────────────[0m[36m┐[0m
[34m151[0m:[38;5;231m gioui.org/shader v1.0.6/go.mod h1:mWdiME581d/kV7/iEhLmUgUK5iZ09XR5XpduXzbePVM= [0m[36m│[0m
[36m────────────────────────────────────────────────────────────────────────────────────[0m[36m┴[0m[36m─────────────────────[0m
[34m[38;5;231mgit.sr.ht/~sbinet/gg v0.3.1 h1:LNhjNn8DerC8f9DHLz6[34m↵[0m  [34m│ [38;5;231mgit.sr.ht/~sbinet/gg v0.3.1 h1:LNhjNn8DerC8f9DHLz6[34m↵[0m
[34m[38;5;231mlS0YYul/b602DUxDgGkd/Aik=[0m                            [34m│ [38;5;231mlS0YYul/b602DUxDgGkd/Aik=[0m
[34m[38;5;231mgit.sr.ht/~sbinet/gg v0.3.1/go.mod h1:KGYtlADtqsqA[34m↵[0m  [34m│ [38;5;231mgit.sr.ht/~sbinet/gg v0.3.1/go.mod h1:KGYtlADtqsqA[34m↵[0m
[34m[38;5;231mNL9ueOFkWymvzUvLMQllU5Ixo+8v3pc=[0m                     [34m│ [38;5;231mNL9ueOFkWymvzUvLMQllU5Ixo+8v3pc=[0m
[34m[38;5;231mgithub.com/Azure/azure-amqp-common-go/v3 v3.2.3/go[34m↵[0m  [34m│ [38;5;231mgithub.com/Azure/azure-amqp-common-go/v3 v3.2.3/go[34m↵[0m
[34m[38;5;231m.mod h1:7rPmbSfszeovxGfc5fSAXE4ehlXQZHpMja2OtxC2Ta[34m↵[0m  [34m│ [38;5;231m.mod h1:7rPmbSfszeovxGfc5fSAXE4ehlXQZHpMja2OtxC2Ta[34m↵[0m
[34m[38;5;231ms=[0m                                                   [34m│ [38;5;231ms=[0m
[34m[48;5;52;38;5;231mgithub.com/Azure/azure-sdk-for-go [48;5;124mv67[48;5;52m.[48;5;124m3[48;5;52m.0+incompatib[34m↵[0m[34m│ [48;5;22;38;5;231mgithub.com/Azure/azure-sdk-for-go [48;5;28mv68[48;5;22m.[48;5;28m0[48;5;22m.0+incompat[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mle/go.mod h1:9XXNKU+eRnpl9moKnB4QOLf1HestfXbmab5FXxi[34m↵[0m[34m│ [48;5;22;38;5;231mible/go.mod h1:9XXNKU+eRnpl9moKnB4QOLf1HestfXbmab5[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mDBjc=[0m[48;5;52m                                                [0m[34m│ [48;5;22;38;5;231mFXxiDBjc=[0m[48;5;22m[0K[0m
[34m[38;5;231mgithub.com/Azure/azure-sdk-for-go/sdk/azcore v1.3.[34m↵[0m  [34m│ [38;5;231mgithub.com/Azure/azure-sdk-for-go/sdk/azcore v1.3.[34m↵[0m
[34m[38;5;231m0/go.mod h1:tZoQYdDZNOiIjdSn0dVWVfl0NEPGOJqVLzSrcF[34m↵[0m  [34m│ [38;5;231m0/go.mod h1:tZoQYdDZNOiIjdSn0dVWVfl0NEPGOJqVLzSrcF[34m↵[0m
[34m[38;5;231mk4Is0=[0m                                               [34m│ [38;5;231mk4Is0=[0m
[34m[38;5;231mgithub.com/Azure/azure-sdk-for-go/sdk/azidentity v[34m↵[0m  [34m│ [38;5;231mgithub.com/Azure/azure-sdk-for-go/sdk/azidentity v[34m↵[0m
[34m[38;5;231m1.2.0/go.mod h1:NBanQUfSWiWn3QEpWDTCU0IjBECKOYvl2R[34m↵[0m  [34m│ [38;5;231m1.2.0/go.mod h1:NBanQUfSWiWn3QEpWDTCU0IjBECKOYvl2R[34m↵[0m
[34m[38;5;231m8xdRtMtiM=[0m                                           [34m│ [38;5;231m8xdRtMtiM=[0m
[34m[38;5;231mgithub.com/Azure/azure-sdk-for-go/sdk/internal v1.[34m↵[0m  [34m│ [38;5;231mgithub.com/Azure/azure-sdk-for-go/sdk/internal v1.[34m↵[0m
[34m[38;5;231m1.2/go.mod h1:eWRD7oawr1Mu1sLCawqVc0CUiF43ia3qQMxL[34m↵[0m  [34m│ [38;5;231m1.2/go.mod h1:eWRD7oawr1Mu1sLCawqVc0CUiF43ia3qQMxL[34m↵[0m
[34m[38;5;231mscsKQ9w=[0m                                             [34m│ [38;5;231mscsKQ9w=[0m

[36m──────────────────────────────────────────────────────────────────────────────────────[0m[36m┐[0m
[34m187[0m:[38;5;231m github.com/ajstarks/svgo v0.0.0-20211024235047-1546f124cd8b/go.mod h1:1KcenG0jGW [0m[36m│[0m
[36m──────────────────────────────────────────────────────────────────────────────────────[0m[36m┴[0m[36m───────────────────[0m
[34m[38;5;231mgithub.com/akavel/rsrc v0.8.0/go.mod h1:uLoCtb9J+E[34m↵[0m  [34m│ [38;5;231mgithub.com/akavel/rsrc v0.8.0/go.mod h1:uLoCtb9J+E[34m↵[0m
[34m[38;5;231myAqh+26kdrTgmzRBFPGOolLWKpdxkKq+c=[0m                   [34m│ [38;5;231myAqh+26kdrTgmzRBFPGOolLWKpdxkKq+c=[0m
[34m[38;5;231mgithub.com/antihax/optional v1.0.0/go.mod h1:uupD/[34m↵[0m  [34m│ [38;5;231mgithub.com/antihax/optional v1.0.0/go.mod h1:uupD/[34m↵[0m
[34m[38;5;231m76wgC+ih3iEmQUL+0Ugr19nfwCT1kdvxnR2qWY=[0m              [34m│ [38;5;231m76wgC+ih3iEmQUL+0Ugr19nfwCT1kdvxnR2qWY=[0m
[34m[38;5;231mgithub.com/armon/go-socks5 v0.0.0-20160902184237-e[34m↵[0m  [34m│ [38;5;231mgithub.com/armon/go-socks5 v0.0.0-20160902184237-e[34m↵[0m
[34m[38;5;231m75332964ef5 h1:0CwZNZbxp69SHPdPJAN/hZIm0C4OItdklCF[34m↵[0m  [34m│ [38;5;231m75332964ef5 h1:0CwZNZbxp69SHPdPJAN/hZIm0C4OItdklCF[34m↵[0m
[34m[38;5;231mmMRWYpio=[0m                                            [34m│ [38;5;231mmMRWYpio=[0m
[34m[48;5;52;38;5;231mgithub.com/aws/aws-sdk-go v1.44.[48;5;124m181[48;5;52m h1:[48;5;124mw4OzE8bwIVo62[48;5;52;34m↵[0m[34m│ [48;5;22;38;5;231mgithub.com/aws/aws-sdk-go v1.44.[48;5;28m185[48;5;22m h1:[48;5;28mstasiou+Ucx[48;5;22;34m↵[0m[48;5;22m[0K[0m
[34m[48;5;124;38;5;231mgUTAp/uEFO2HSsUtf1pjXpSs36cluY[48;5;52m=[0m[48;5;52m                      [0m[34m│ [48;5;28;38;5;231m2A0RyXRyPph4sLCBxVQK7DPPK8tNcl5g[48;5;22m=[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mgithub.com/aws/aws-sdk-go v1.44.[48;5;124m181[48;5;52m/go.mod h1:aVsgQc[34m↵[0m[34m│ [48;5;22;38;5;231mgithub.com/aws/aws-sdk-go v1.44.[48;5;28m185[48;5;22m/go.mod h1:aVsg[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mEevwlmQ7qHE9I3h+dtQgpqhFB+i8Phjh7fkwI=[0m[48;5;52m               [0m[34m│ [48;5;22;38;5;231mQcEevwlmQ7qHE9I3h+dtQgpqhFB+i8Phjh7fkwI=[0m[48;5;22m[0K[0m
[34m[38;5;231mgithub.com/aws/aws-sdk-go-v2 v1.17.3 h1:shN7NlnVzv[34m↵[0m  [34m│ [38;5;231mgithub.com/aws/aws-sdk-go-v2 v1.17.3 h1:shN7NlnVzv[34m↵[0m
[34m[38;5;231mDUgPQ+1rLMSxY8OWRNDRYtiqe0p/PgrhY=[0m                   [34m│ [38;5;231mDUgPQ+1rLMSxY8OWRNDRYtiqe0p/PgrhY=[0m
[34m[38;5;231mgithub.com/aws/aws-sdk-go-v2 v1.17.3/go.mod h1:uzb[34m↵[0m  [34m│ [38;5;231mgithub.com/aws/aws-sdk-go-v2 v1.17.3/go.mod h1:uzb[34m↵[0m
[34m[38;5;231mQtefpm44goOPmdKyAlXSNcwlRgF3ePWVW6EtJvvw=[0m            [34m│ [38;5;231mQtefpm44goOPmdKyAlXSNcwlRgF3ePWVW6EtJvvw=[0m
[34m[38;5;231mgithub.com/aws/aws-sdk-go-v2/aws/protocol/eventstr[34m↵[0m  [34m│ [38;5;231mgithub.com/aws/aws-sdk-go-v2/aws/protocol/eventstr[34m↵[0m
[34m[38;5;231meam v1.4.10 h1:dK82zF6kkPeCo8J1e+tGx4JdvDIQzj7ygIo[34m↵[0m  [34m│ [38;5;231meam v1.4.10 h1:dK82zF6kkPeCo8J1e+tGx4JdvDIQzj7ygIo[34m↵[0m
[34m[38;5;231mLg8WMuGs=[0m                                            [34m│ [38;5;231mLg8WMuGs=[0m
[34m[38;5;231mgithub.com/aws/aws-sdk-go-v2/aws/protocol/eventstr[34m↵[0m  [34m│ [38;5;231mgithub.com/aws/aws-sdk-go-v2/aws/protocol/eventstr[34m↵[0m
[34m[38;5;231meam v1.4.10/go.mod h1:VeTZetY5KRJLuD/7fkQXMU6Mw7H5[34m↵[0m  [34m│ [38;5;231meam v1.4.10/go.mod h1:VeTZetY5KRJLuD/7fkQXMU6Mw7H5[34m↵[0m
[34m[38;5;231mm/KP2J5Iy9osMno=[0m                                     [34m│ [38;5;231mm/KP2J5Iy9osMno=[0m
[34m[48;5;52;38;5;231mgithub.com/aws/aws-sdk-go-v2/config v1.18.[48;5;124m8[48;5;52m h1:[48;5;124mlDpy0[48;5;52;34m↵[0m[34m│ [48;5;22;38;5;231mgithub.com/aws/aws-sdk-go-v2/config v1.18.[48;5;28m9[48;5;22m h1:[48;5;28mpd+[48;5;22;34m↵[0m[48;5;22m[0K[0m
[34m[48;5;124;38;5;231mWM8AHsywOnVrOHaSMfpaiV2igOw8D7svkFkXVA[48;5;52m=[0m[48;5;52m              [0m[34m│ [48;5;28;38;5;231mQUO1dvro6vGOuhgglJV6adGunU95xSTSzsQGhKpY[48;5;22m=[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mgithub.com/aws/aws-sdk-go-v2/config v1.18.[48;5;124m8[48;5;52m/go.mod h[34m↵[0m[34m│ [48;5;22;38;5;231mgithub.com/aws/aws-sdk-go-v2/config v1.18.[48;5;28m9[48;5;22m/go.mod[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231m1:[48;5;124m5XCmmyutmzzgkpk[48;5;52m/[48;5;124m6NYTjeWb6lgo9N170m1j6pQkIBs[48;5;52m=[0m[48;5;52m       [0m[34m│ [48;5;22;38;5;231m h1:[48;5;28m2Lx9yaA[48;5;22m/[48;5;28mMcDeQS8ft+edKrmOd5ry1v1euFQ+oGwUxsM[48;5;22m=[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mgithub.com/aws/aws-sdk-go-v2/credentials v1.13.[48;5;124m8[48;5;52m h1:[34m↵[0m[34m│ [48;5;22;38;5;231mgithub.com/aws/aws-sdk-go-v2/credentials v1.13.[48;5;28m9[48;5;22m h[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;124;38;5;231mvTrwTvv5qAwjWIGhZDSBH[48;5;52m/[48;5;124moQHuIQjGmD232k01FUh6A[48;5;52m=[0m[48;5;52m         [0m[34m│ [48;5;22;38;5;231m1:[48;5;28moxM[48;5;22m/[48;5;28mC8eXGsiHH+u0gZGo1++QTFPf+N5MUb1tfaaQMpU[48;5;22m=[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mgithub.com/aws/aws-sdk-go-v2/credentials v1.13.[48;5;124m8[48;5;52m/go.[34m↵[0m[34m│ [48;5;22;38;5;231mgithub.com/aws/aws-sdk-go-v2/credentials v1.13.[48;5;28m9[48;5;22m/g[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mmod h1:[48;5;124mlVa4OHbvgjVot4gmh1uouF1ubgexSCN92P6CJQpT0t8[48;5;52m=[0m[48;5;52m  [0m[34m│ [48;5;22;38;5;231mo.mod h1:[48;5;28m45DrDZTok50mEx4Uw59ym7n11Oy7G4gt0Pez2Z4kt[48;5;22;34m↵[0m[48;5;22m[0K[0m
[34m[0m                                                     [34m│ [48;5;28;38;5;231mAA[48;5;22m=[0m[48;5;22m[0K[0m
[34m[38;5;231mgithub.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.1[34m↵[0m  [34m│ [38;5;231mgithub.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.1[34m↵[0m
[34m[38;5;231m2.21 h1:j9wi1kQ8b+e0FBVHxCqCGo4kxDU175hoDHcWAi0sau[34m↵[0m  [34m│ [38;5;231m2.21 h1:j9wi1kQ8b+e0FBVHxCqCGo4kxDU175hoDHcWAi0sau[34m↵[0m
[34m[38;5;231mU=[0m                                                   [34m│ [38;5;231mU=[0m
[34m[38;5;231mgithub.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.1[34m↵[0m  [34m│ [38;5;231mgithub.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.1[34m↵[0m
[34m[38;5;231m2.21/go.mod h1:ugwW57Z5Z48bpvUyZuaPy4Kv+vEfJWnIrky[34m↵[0m  [34m│ [38;5;231m2.21/go.mod h1:ugwW57Z5Z48bpvUyZuaPy4Kv+vEfJWnIrky[34m↵[0m
[34m[38;5;231m7RmkBvJg=[0m                                            [34m│ [38;5;231m7RmkBvJg=[0m
[34m[48;5;52;38;5;231mgithub.com/aws/aws-sdk-go-v2/feature/s3/manager v1.1[34m↵[0m[34m│ [48;5;22;38;5;231mgithub.com/aws/aws-sdk-go-v2/feature/s3/manager v1[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231m1.[48;5;124m47[48;5;52m h1:[48;5;124mE884ndKWVGt8IhtUuGhXbEsmaCvdAAkTTUDu7uAok1g[48;5;52m=[0m[48;5;52m [0m[34m│ [48;5;22;38;5;231m.11.[48;5;28m48[48;5;22m h1:[48;5;28m3IGeA7Vh+gpp6Ptf0slDgNwFVTJEu81IiGl1v5yG[48;5;22;34m↵[0m[48;5;22m[0K[0m
[34m[0m                                                     [34m│ [48;5;28;38;5;231mZ3A[48;5;22m=[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mgithub.com/aws/aws-sdk-go-v2/feature/s3/manager v1.1[34m↵[0m[34m│ [48;5;22;38;5;231mgithub.com/aws/aws-sdk-go-v2/feature/s3/manager v1[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231m1.[48;5;124m47[48;5;52m/go.mod h1:[48;5;124mKybsEsmXLO0u75FyS3F0sY4OQ97syDe8z+ISq[48;5;52;34m↵[0m[34m│ [48;5;22;38;5;231m.11.[48;5;28m48[48;5;22m/go.mod h1:[48;5;28mkZ8I3L92ide4A8rLSEHofGn43eLE7E/m9[48;5;22;34m↵[0m[48;5;22m[0K[0m
[34m[48;5;124;38;5;231m8oEczA[48;5;52m=[0m[48;5;52m                                              [0m[34m│ [48;5;28;38;5;231mH986uub0ns[48;5;22m=[0m[48;5;22m[0K[0m
[34m[38;5;231mgithub.com/aws/aws-sdk-go-v2/internal/configsource[34m↵[0m  [34m│ [38;5;231mgithub.com/aws/aws-sdk-go-v2/internal/configsource[34m↵[0m
[34m[38;5;231ms v1.1.27 h1:I3cakv2Uy1vNmmhRQmFptYDxOvBnwCdNwyw63[34m↵[0m  [34m│ [38;5;231ms v1.1.27 h1:I3cakv2Uy1vNmmhRQmFptYDxOvBnwCdNwyw63[34m↵[0m
[34m[38;5;231mN0RaRU=[0m                                              [34m│ [38;5;231mN0RaRU=[0m
[34m[38;5;231mgithub.com/aws/aws-sdk-go-v2/internal/configsource[34m↵[0m  [34m│ [38;5;231mgithub.com/aws/aws-sdk-go-v2/internal/configsource[34m↵[0m
[34m[38;5;231ms v1.1.27/go.mod h1:a1/UpzeyBBerajpnP5nGZa9mGzsBn5[34m↵[0m  [34m│ [38;5;231ms v1.1.27/go.mod h1:a1/UpzeyBBerajpnP5nGZa9mGzsBn5[34m↵[0m
[34m[38;5;231mcOKxm6NWQsvoI=[0m                                       [34m│ [38;5;231mcOKxm6NWQsvoI=[0m
[34m[38;5;231mgithub.com/aws/aws-sdk-go-v2/internal/endpoints/v2[34m↵[0m  [34m│ [38;5;231mgithub.com/aws/aws-sdk-go-v2/internal/endpoints/v2[34m↵[0m
[34m[38;5;231m v2.4.21 h1:5NbbMrIzmUn/TXFqAle6mgrH5m9cOvMLRGL7pn[34m↵[0m  [34m│ [38;5;231m v2.4.21 h1:5NbbMrIzmUn/TXFqAle6mgrH5m9cOvMLRGL7pn[34m↵[0m
[34m[38;5;231mG8tRE=[0m                                               [34m│ [38;5;231mG8tRE=[0m

[36m──────────────────────────────────────────────────────────────────────────────────────[0m[36m┐[0m
[34m217[0m:[38;5;231m github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.21 h1:5C6XgTViS [0m[36m│[0m
[36m──────────────────────────────────────────────────────────────────────────────────────[0m[36m┴[0m[36m───────────────────[0m
[34m[38;5;231mgithub.com/aws/aws-sdk-go-v2/service/internal/pres[34m↵[0m  [34m│ [38;5;231mgithub.com/aws/aws-sdk-go-v2/service/internal/pres[34m↵[0m
[34m[38;5;231migned-url v1.9.21/go.mod h1:lRToEJsn+DRA9lW4O9L9+/[34m↵[0m  [34m│ [38;5;231migned-url v1.9.21/go.mod h1:lRToEJsn+DRA9lW4O9L9+/[34m↵[0m
[34m[38;5;231m3hjTkUzlzyzHqn8MTds5k=[0m                               [34m│ [38;5;231m3hjTkUzlzyzHqn8MTds5k=[0m
[34m[38;5;231mgithub.com/aws/aws-sdk-go-v2/service/internal/s3sh[34m↵[0m  [34m│ [38;5;231mgithub.com/aws/aws-sdk-go-v2/service/internal/s3sh[34m↵[0m
[34m[38;5;231mared v1.13.21 h1:vY5siRXvW5TrOKm2qKEf9tliBfdLxdfy0[34m↵[0m  [34m│ [38;5;231mared v1.13.21 h1:vY5siRXvW5TrOKm2qKEf9tliBfdLxdfy0[34m↵[0m
[34m[38;5;231mi02LOcmqUo=[0m                                          [34m│ [38;5;231mi02LOcmqUo=[0m
[34m[38;5;231mgithub.com/aws/aws-sdk-go-v2/service/internal/s3sh[34m↵[0m  [34m│ [38;5;231mgithub.com/aws/aws-sdk-go-v2/service/internal/s3sh[34m↵[0m
[34m[38;5;231mared v1.13.21/go.mod h1:WZvNXT1XuH8dnJM0HvOlvk+RNn[34m↵[0m  [34m│ [38;5;231mared v1.13.21/go.mod h1:WZvNXT1XuH8dnJM0HvOlvk+RNn[34m↵[0m
[34m[38;5;231m7NbAPvA/ACO0QarSc=[0m                                   [34m│ [38;5;231m7NbAPvA/ACO0QarSc=[0m
[34m[48;5;52;38;5;231mgithub.com/aws/aws-sdk-go-v2/service/kms v1.20.[48;5;124m0[48;5;52m/go.[34m↵[0m[34m│ [48;5;22;38;5;231mgithub.com/aws/aws-sdk-go-v2/service/kms v1.20.[48;5;28m1[48;5;22m/g[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mmod h1:13sjgMH7Xu4e46+0BEDhSnNh+cImHSYS5PpBjV3oXcU=[0m[48;5;52m  [0m[34m│ [48;5;22;38;5;231mo.mod h1:13sjgMH7Xu4e46+0BEDhSnNh+cImHSYS5PpBjV3oX[34m↵[0m[48;5;22m[0K[0m
[34m[0m                                                     [34m│ [48;5;22;38;5;231mcU=[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mgithub.com/aws/aws-sdk-go-v2/service/s3 v1.30.[48;5;124m0[48;5;52m h1:[48;5;124mw[48;5;52;34m↵[0m[34m│ [48;5;22;38;5;231mgithub.com/aws/aws-sdk-go-v2/service/s3 v1.30.[48;5;28m1[48;5;22m h1[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;124;38;5;231mddsyuESfviaiXk3w9N6/4iRwTg[48;5;52m/[48;5;124ma3gktjODY6jYQBo[48;5;52m=[0m[48;5;52m          [0m[34m│ [48;5;22;38;5;231m:[48;5;28mkIgvVY7PHx4gIb0na[48;5;22m/[48;5;28mQ9gTWJWauTwhKdaqJjX8PkIY8[48;5;22m=[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mgithub.com/aws/aws-sdk-go-v2/service/s3 v1.30.[48;5;124m0[48;5;52m/go.m[34m↵[0m[34m│ [48;5;22;38;5;231mgithub.com/aws/aws-sdk-go-v2/service/s3 v1.30.[48;5;28m1[48;5;22m/go[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mod h1:L2l2/q76teehcW7YEsgsDjqdsDTERJeX3nOMIFlgGUE=[0m[48;5;52m   [0m[34m│ [48;5;22;38;5;231m.mod h1:L2l2/q76teehcW7YEsgsDjqdsDTERJeX3nOMIFlgGU[34m↵[0m[48;5;22m[0K[0m
[34m[0m                                                     [34m│ [48;5;22;38;5;231mE=[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mgithub.com/aws/aws-sdk-go-v2/service/secretsmanager [34m↵[0m[34m│ [48;5;22;38;5;231mgithub.com/aws/aws-sdk-go-v2/service/secretsmanage[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mv1.18.[48;5;124m1[48;5;52m/go.mod h1:jAeo/PdIJZuDSwsvxJS94G4d6h8tStj7WX[34m↵[0m[34m│ [48;5;22;38;5;231mr v1.18.[48;5;28m2[48;5;22m/go.mod h1:jAeo/PdIJZuDSwsvxJS94G4d6h8tSt[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mVuKwLHWU8=[0m[48;5;52m                                           [0m[34m│ [48;5;22;38;5;231mj7WXVuKwLHWU8=[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mgithub.com/aws/aws-sdk-go-v2/service/sns v1.19.[48;5;124m0[48;5;52m/go.[34m↵[0m[34m│ [48;5;22;38;5;231mgithub.com/aws/aws-sdk-go-v2/service/sns v1.19.[48;5;28m1[48;5;22m/g[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mmod h1:iTh9DgwDnFqF5LfFHNXWAxLe9zV0/XcWaMCWXIRDqXA=[0m[48;5;52m  [0m[34m│ [48;5;22;38;5;231mo.mod h1:iTh9DgwDnFqF5LfFHNXWAxLe9zV0/XcWaMCWXIRDq[34m↵[0m[48;5;22m[0K[0m
[34m[0m                                                     [34m│ [48;5;22;38;5;231mXA=[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mgithub.com/aws/aws-sdk-go-v2/service/sqs v1.20.[48;5;124m0[48;5;52m/go.[34m↵[0m[34m│ [48;5;22;38;5;231mgithub.com/aws/aws-sdk-go-v2/service/sqs v1.20.[48;5;28m1[48;5;22m/g[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mmod h1:jQhN5f4p3PALMNlUtfb/0wGIFlV7vGtJlPDVfxfNfPY=[0m[48;5;52m  [0m[34m│ [48;5;22;38;5;231mo.mod h1:jQhN5f4p3PALMNlUtfb/0wGIFlV7vGtJlPDVfxfNf[34m↵[0m[48;5;22m[0K[0m
[34m[0m                                                     [34m│ [48;5;22;38;5;231mPY=[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mgithub.com/aws/aws-sdk-go-v2/service/ssm v1.35.[48;5;124m0[48;5;52m/go.[34m↵[0m[34m│ [48;5;22;38;5;231mgithub.com/aws/aws-sdk-go-v2/service/ssm v1.35.[48;5;28m1[48;5;22m/g[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mmod h1:Hf7wSogKP1XCJ9GgW8erZDL6IZ1NLwLN7bYdV/Gn/LI=[0m[48;5;52m  [0m[34m│ [48;5;22;38;5;231mo.mod h1:Hf7wSogKP1XCJ9GgW8erZDL6IZ1NLwLN7bYdV/Gn/[34m↵[0m[48;5;22m[0K[0m
[34m[0m                                                     [34m│ [48;5;22;38;5;231mLI=[0m[48;5;22m[0K[0m
[34m[38;5;231mgithub.com/aws/aws-sdk-go-v2/service/sso v1.12.0 h[34m↵[0m  [34m│ [38;5;231mgithub.com/aws/aws-sdk-go-v2/service/sso v1.12.0 h[34m↵[0m
[34m[38;5;231m1:/2gzjhQowRLarkkBOGPXSRnb8sQ2RVsjdG1C/UliK/c=[0m       [34m│ [38;5;231m1:/2gzjhQowRLarkkBOGPXSRnb8sQ2RVsjdG1C/UliK/c=[0m
[34m[38;5;231mgithub.com/aws/aws-sdk-go-v2/service/sso v1.12.0/g[34m↵[0m  [34m│ [38;5;231mgithub.com/aws/aws-sdk-go-v2/service/sso v1.12.0/g[34m↵[0m
[34m[38;5;231mo.mod h1:wo/B7uUm/7zw/dWhBJ4FXuw1sySU5lyIhVg1Bu2yL[34m↵[0m  [34m│ [38;5;231mo.mod h1:wo/B7uUm/7zw/dWhBJ4FXuw1sySU5lyIhVg1Bu2yL[34m↵[0m
[34m[38;5;231m9A=[0m                                                  [34m│ [38;5;231m9A=[0m
[34m[38;5;231mgithub.com/aws/aws-sdk-go-v2/service/ssooidc v1.14[34m↵[0m  [34m│ [38;5;231mgithub.com/aws/aws-sdk-go-v2/service/ssooidc v1.14[34m↵[0m
[34m[38;5;231m.0 h1:Jfly6mRxk2ZOSlbCvZfKNS7TukSx1mIzhSsqZ/IGSZI=[0m   [34m│ [38;5;231m.0 h1:Jfly6mRxk2ZOSlbCvZfKNS7TukSx1mIzhSsqZ/IGSZI=[0m
[34m[38;5;231mgithub.com/aws/aws-sdk-go-v2/service/ssooidc v1.14[34m↵[0m  [34m│ [38;5;231mgithub.com/aws/aws-sdk-go-v2/service/ssooidc v1.14[34m↵[0m
[34m[38;5;231m.0/go.mod h1:TZSH7xLO7+phDtViY/KUp9WGCJMQkLJ/VpgkT[34m↵[0m  [34m│ [38;5;231m.0/go.mod h1:TZSH7xLO7+phDtViY/KUp9WGCJMQkLJ/VpgkT[34m↵[0m
[34m[38;5;231mFd5gh8=[0m                                              [34m│ [38;5;231mFd5gh8=[0m
[34m[48;5;52;38;5;231mgithub.com/aws/aws-sdk-go-v2/service/sts v1.18.[48;5;124m0[48;5;52m h1:[34m↵[0m[34m│ [48;5;22;38;5;231mgithub.com/aws/aws-sdk-go-v2/service/sts v1.18.[48;5;28m1[48;5;22m h[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;124;38;5;231mkOO++CYo50RcTFISESluhWEi5Prhg+gaSs4whWabiZU[48;5;52m=[0m[48;5;52m         [0m[34m│ [48;5;22;38;5;231m1:[48;5;28mq3xG67qnKp1gsYSJY5AtTvFKY2IlmGPGrTw/Wy8EjeQ[48;5;22m=[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mgithub.com/aws/aws-sdk-go-v2/service/sts v1.18.[48;5;124m0[48;5;52m/go.[34m↵[0m[34m│ [48;5;22;38;5;231mgithub.com/aws/aws-sdk-go-v2/service/sts v1.18.[48;5;28m1[48;5;22m/g[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mmod h1:+lGbb3+1ugwKrNTWcf2RT05Xmp543B06zDFTwiTLp7I=[0m[48;5;52m  [0m[34m│ [48;5;22;38;5;231mo.mod h1:+lGbb3+1ugwKrNTWcf2RT05Xmp543B06zDFTwiTLp[34m↵[0m[48;5;22m[0K[0m
[34m[0m                                                     [34m│ [48;5;22;38;5;231m7I=[0m[48;5;22m[0K[0m
[34m[38;5;231mgithub.com/aws/smithy-go v1.13.5 h1:hgz0X/DX0dGqTY[34m↵[0m  [34m│ [38;5;231mgithub.com/aws/smithy-go v1.13.5 h1:hgz0X/DX0dGqTY[34m↵[0m
[34m[38;5;231mpGALqXJoRKRj5oQ7150i5FdTePzO8=[0m                       [34m│ [38;5;231mpGALqXJoRKRj5oQ7150i5FdTePzO8=[0m
[34m[38;5;231mgithub.com/aws/smithy-go v1.13.5/go.mod h1:Tg+OJXh[34m↵[0m  [34m│ [38;5;231mgithub.com/aws/smithy-go v1.13.5/go.mod h1:Tg+OJXh[34m↵[0m
[34m[38;5;231m4MB2R/uN61Ko2f6hTZwB/ZYGOtib8J3gBHzA=[0m                [34m│ [38;5;231m4MB2R/uN61Ko2f6hTZwB/ZYGOtib8J3gBHzA=[0m
[34m[38;5;231mgithub.com/benbjohnson/clock v1.3.0 h1:ip6w0uFQknc[34m↵[0m  [34m│ [38;5;231mgithub.com/benbjohnson/clock v1.3.0 h1:ip6w0uFQknc[34m↵[0m
[34m[38;5;231mKQ979AypyG0ER7mqUSBdKLOgAle/AT8A=[0m                    [34m│ [38;5;231mKQ979AypyG0ER7mqUSBdKLOgAle/AT8A=[0m

[36m──────────────────────────────────────────────────────────────────────────────────────[0m[36m┐[0m
[34m450[0m:[38;5;231m github.com/json-iterator/go v1.1.12/go.mod h1:e30LSqwooZae/UwlEbR2852Gd8hjQvJoHm [0m[36m│[0m
[36m──────────────────────────────────────────────────────────────────────────────────────[0m[36m┴[0m[36m───────────────────[0m
[34m[38;5;231mgithub.com/jstemmer/go-junit-report v1.0.0/go.mod [34m↵[0m  [34m│ [38;5;231mgithub.com/jstemmer/go-junit-report v1.0.0/go.mod [34m↵[0m
[34m[38;5;231mh1:Brl9GWCQeLvo8nXZwPNNblvFj/XSXhF0NWZEnDohbsk=[0m      [34m│ [38;5;231mh1:Brl9GWCQeLvo8nXZwPNNblvFj/XSXhF0NWZEnDohbsk=[0m
[34m[38;5;231mgithub.com/kisielk/errcheck v1.6.3/go.mod h1:nXw/i[34m↵[0m  [34m│ [38;5;231mgithub.com/kisielk/errcheck v1.6.3/go.mod h1:nXw/i[34m↵[0m
[34m[38;5;231m/MfnvRHqXa7XXmQMUB0oNFGuBrNI8d8NLy0LPw=[0m              [34m│ [38;5;231m/MfnvRHqXa7XXmQMUB0oNFGuBrNI8d8NLy0LPw=[0m
[34m[38;5;231mgithub.com/kisielk/gotool v1.0.0/go.mod h1:XhKaO+M[34m↵[0m  [34m│ [38;5;231mgithub.com/kisielk/gotool v1.0.0/go.mod h1:XhKaO+M[34m↵[0m
[34m[38;5;231mFFWcvkIS/tQcRk01m1F5IRFswLeQ+oQHNcck=[0m                [34m│ [38;5;231mFFWcvkIS/tQcRk01m1F5IRFswLeQ+oQHNcck=[0m
[34m[48;5;52;38;5;231mgithub.com/klauspost/compress v1.15.15-0.20230116143[34m↵[0m[34m│ [0m
[34m[48;5;52;38;5;231m836-fbae784ff625 h1:UB6LqX5EKpt3veNECC3tjVQ57x6REVBF[34m↵[0m[34m│ [0m
[34m[48;5;52;38;5;231mtAa+pT9dKFk=[0m[48;5;52m                                         [0m[34m│ [0m
[34m[0m                                                     [34m│ [48;5;22;38;5;231mgithub.com/klauspost/compress v1.15.16-0.202301211[34m↵[0m[48;5;22m[0K[0m
[34m[0m                                                     [34m│ [48;5;22;38;5;231m71712-fe37dc6783c8 h1:bjJDGObEyTZ52MP9jEjuoAoqsHS/[34m↵[0m[48;5;22m[0K[0m
[34m[0m                                                     [34m│ [48;5;22;38;5;231mHcLtUlwz3biCsC0=[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mgithub.com/klauspost/compress v1.15.[48;5;124m15[48;5;52m-0.[48;5;124m20230116143[48;5;52;34m↵[0m[34m│ [48;5;22;38;5;231mgithub.com/klauspost/compress v1.15.[48;5;28m16[48;5;22m-0.[48;5;28m202301211[48;5;22;34m↵[0m[48;5;22m[0K[0m
[34m[48;5;124;38;5;231m836[48;5;52m-[48;5;124mfbae784ff625[48;5;52m/go.mod h1:ZcK2JAFqKOpnBlxcLsJzYfrS9[34m↵[0m[34m│ [48;5;28;38;5;231m71712[48;5;22m-[48;5;28mfe37dc6783c8[48;5;22m/go.mod h1:ZcK2JAFqKOpnBlxcLsJzY[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mX1akm9fHZNnD9+Vo/4=[0m[48;5;52m                                  [0m[34m│ [48;5;22;38;5;231mfrS9X1akm9fHZNnD9+Vo/4=[0m[48;5;22m[0K[0m
[34m[38;5;231mgithub.com/klauspost/cpuid/v2 v2.2.3 h1:sxCkb+qR91[34m↵[0m  [34m│ [38;5;231mgithub.com/klauspost/cpuid/v2 v2.2.3 h1:sxCkb+qR91[34m↵[0m
[34m[38;5;231mz4vsqw4vGGZlDgPz3G7gjaLyK3V8y70BU=[0m                   [34m│ [38;5;231mz4vsqw4vGGZlDgPz3G7gjaLyK3V8y70BU=[0m
[34m[38;5;231mgithub.com/klauspost/cpuid/v2 v2.2.3/go.mod h1:RVV[34m↵[0m  [34m│ [38;5;231mgithub.com/klauspost/cpuid/v2 v2.2.3/go.mod h1:RVV[34m↵[0m
[34m[38;5;231moqg1df56z8g3pUjL/3lE5UfnlrJX8tyFgg4nqhuY=[0m            [34m│ [38;5;231moqg1df56z8g3pUjL/3lE5UfnlrJX8tyFgg4nqhuY=[0m
[34m[38;5;231mgithub.com/kpango/fastime v1.1.6 h1:lAw1Tiwnlbsx1x[34m↵[0m  [34m│ [38;5;231mgithub.com/kpango/fastime v1.1.6 h1:lAw1Tiwnlbsx1x[34m↵[0m
[34m[38;5;231mZs6W9eM7/8niwabknewbmLkh/yTVo=[0m                       [34m│ [38;5;231mZs6W9eM7/8niwabknewbmLkh/yTVo=[0m

[36m────────────────────────────────────────────────────────────────────────────────────[0m[36m┐[0m
[34m512[0m:[38;5;231m github.com/onsi/ginkgo v1.16.5 h1:8xi0RTUf59SOSfEtZMvwTvXYMzG4gV23XVHOZiXNtnE= [0m[36m│[0m
[36m────────────────────────────────────────────────────────────────────────────────────[0m[36m┴[0m[36m─────────────────────[0m
[34m[38;5;231mgithub.com/onsi/ginkgo v1.16.5/go.mod h1:+E8gABHa3[34m↵[0m  [34m│ [38;5;231mgithub.com/onsi/ginkgo v1.16.5/go.mod h1:+E8gABHa3[34m↵[0m
[34m[38;5;231mK6zRBolWtd+ROzc/U5bkGt0FwiG042wbpU=[0m                  [34m│ [38;5;231mK6zRBolWtd+ROzc/U5bkGt0FwiG042wbpU=[0m
[34m[38;5;231mgithub.com/onsi/ginkgo/v2 v2.7.0 h1:/XxtEV3I3Eif/H[34m↵[0m  [34m│ [38;5;231mgithub.com/onsi/ginkgo/v2 v2.7.0 h1:/XxtEV3I3Eif/H[34m↵[0m
[34m[38;5;231mobnVx9YmJgk8ENdRsuUmM+fLCFNow=[0m                       [34m│ [38;5;231mobnVx9YmJgk8ENdRsuUmM+fLCFNow=[0m
[34m[38;5;231mgithub.com/onsi/ginkgo/v2 v2.7.0/go.mod h1:yjiuMwP[34m↵[0m  [34m│ [38;5;231mgithub.com/onsi/ginkgo/v2 v2.7.0/go.mod h1:yjiuMwP[34m↵[0m
[34m[38;5;231mokqY1XauOgju45q3sJt6VzQ/Fict1LFVcsAo=[0m                [34m│ [38;5;231mokqY1XauOgju45q3sJt6VzQ/Fict1LFVcsAo=[0m
[34m[48;5;52;38;5;231mgithub.com/onsi/gomega v1.[48;5;124m25[48;5;52m.0 h1:[48;5;124mVw7br2PCDYijJHSfBO[48;5;52;34m↵[0m[34m│ [48;5;22;38;5;231mgithub.com/onsi/gomega v1.[48;5;28m26[48;5;22m.0 h1:[48;5;28m03cDLK28U6hWvCAn[48;5;22;34m↵[0m[48;5;22m[0K[0m
[34m[48;5;124;38;5;231mWhov+8cAnUf8MfMaIOV323l6Y[48;5;52m=[0m[48;5;52m                           [0m[34m│ [48;5;28;38;5;231ms6NeydX3zIm4SF3ci69ulidS32Q[48;5;22m=[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mgithub.com/onsi/gomega v1.[48;5;124m25[48;5;52m.0/go.mod h1:r+zV744Re+D[34m↵[0m[34m│ [48;5;22;38;5;231mgithub.com/onsi/gomega v1.[48;5;28m26[48;5;22m.0/go.mod h1:r+zV744Re[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231miYCIPRlYOTxn0YkOLcAnW8k1xXdMPGhM=[0m[48;5;52m                    [0m[34m│ [48;5;22;38;5;231m+DiYCIPRlYOTxn0YkOLcAnW8k1xXdMPGhM=[0m[48;5;22m[0K[0m
[34m[38;5;231mgithub.com/pelletier/go-toml/v2 v2.0.6/go.mod h1:e[34m↵[0m  [34m│ [38;5;231mgithub.com/pelletier/go-toml/v2 v2.0.6/go.mod h1:e[34m↵[0m
[34m[38;5;231mumQOmlWiOPt5WriQQqoM5y18pDHwha2N+QD+EUNTek=[0m          [34m│ [38;5;231mumQOmlWiOPt5WriQQqoM5y18pDHwha2N+QD+EUNTek=[0m
[34m[38;5;231mgithub.com/peterbourgon/diskv v2.0.1+incompatible [34m↵[0m  [34m│ [38;5;231mgithub.com/peterbourgon/diskv v2.0.1+incompatible [34m↵[0m
[34m[38;5;231mh1:UBdAOUP5p4RWqPBg048CAvpKN+vxiaj6gdUUzhl4XmI=[0m      [34m│ [38;5;231mh1:UBdAOUP5p4RWqPBg048CAvpKN+vxiaj6gdUUzhl4XmI=[0m
[34m[38;5;231mgithub.com/peterbourgon/diskv v2.0.1+incompatible/[34m↵[0m  [34m│ [38;5;231mgithub.com/peterbourgon/diskv v2.0.1+incompatible/[34m↵[0m
[34m[38;5;231mgo.mod h1:uqqh8zWWbv1HBMNONnaR/tNboyR3/BZd58JJSHlU[34m↵[0m  [34m│ [38;5;231mgo.mod h1:uqqh8zWWbv1HBMNONnaR/tNboyR3/BZd58JJSHlU[34m↵[0m
[34m[38;5;231mSCU=[0m                                                 [34m│ [38;5;231mSCU=[0m

[36m──────────────────────────────────────────────────────────────────────────────────────[0m[36m┐[0m
[34m582[0m:[38;5;231m github.com/stretchr/testify v1.8.1/go.mod h1:w2LPCIKwWwSfY2zedu0+kehJoqGctiVI29o [0m[36m│[0m
[36m──────────────────────────────────────────────────────────────────────────────────────[0m[36m┴[0m[36m───────────────────[0m
[34m[38;5;231mgithub.com/tedsuo/ifrit v0.0.0-20180802180643-bea9[34m↵[0m  [34m│ [38;5;231mgithub.com/tedsuo/ifrit v0.0.0-20180802180643-bea9[34m↵[0m
[34m[38;5;231m4bb476cc/go.mod h1:eyZnKCc955uh98WQvzOm0dgAeLnf2O0[34m↵[0m  [34m│ [38;5;231m4bb476cc/go.mod h1:eyZnKCc955uh98WQvzOm0dgAeLnf2O0[34m↵[0m
[34m[38;5;231mRz0LPoC5ze+0=[0m                                        [34m│ [38;5;231mRz0LPoC5ze+0=[0m
[34m[38;5;231mgithub.com/ugorji/go/codec v1.2.8/go.mod h1:UNopzC[34m↵[0m  [34m│ [38;5;231mgithub.com/ugorji/go/codec v1.2.8/go.mod h1:UNopzC[34m↵[0m
[34m[38;5;231mgEMSXjBc6AOMqYvWC1ktqTAfzJZUZgYf6w6lg=[0m               [34m│ [38;5;231mgEMSXjBc6AOMqYvWC1ktqTAfzJZUZgYf6w6lg=[0m
[34m[38;5;231mgithub.com/urfave/cli/v2 v2.3.0/go.mod h1:LJmUH05z[34m↵[0m  [34m│ [38;5;231mgithub.com/urfave/cli/v2 v2.3.0/go.mod h1:LJmUH05z[34m↵[0m
[34m[38;5;231mAU44vOAcrfzZQKsZbVcdbOG8rtL3/XcUArI=[0m                 [34m│ [38;5;231mAU44vOAcrfzZQKsZbVcdbOG8rtL3/XcUArI=[0m
[34m[48;5;52;38;5;231mgithub.com/vdaas/vald-client-go v1.[48;5;124m6[48;5;52m.[48;5;124m3[48;5;52m h1:[48;5;124m93aY9jOWlt[48;5;52;34m↵[0m[34m│ [48;5;22;38;5;231mgithub.com/vdaas/vald-client-go v1.[48;5;28m7[48;5;22m.[48;5;28m0[48;5;22m h1:[48;5;28mmFkd0/E7[48;5;22;34m↵[0m[48;5;22m[0K[0m
[34m[48;5;124;38;5;231mrO2b8hnIrz5P3NcxWAHMbqY8AnhNpud7w[48;5;52m=[0m[48;5;52m                   [0m[34m│ [48;5;28;38;5;231mOHsAq6of04mUuYTzlbsW+j+daPjRfaau0SA[48;5;22m=[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mgithub.com/vdaas/vald-client-go v1.[48;5;124m6[48;5;52m.[48;5;124m3[48;5;52m/go.mod h1:[48;5;124mWiE[48;5;52;34m↵[0m[34m│ [48;5;22;38;5;231mgithub.com/vdaas/vald-client-go v1.[48;5;28m7[48;5;22m.[48;5;28m0[48;5;22m/go.mod h1:[48;5;28m1[48;5;22;34m↵[0m[48;5;22m[0K[0m
[34m[48;5;124;38;5;231m3uVM1gjAEi4wbQi3S7lwfASR4BMiUvdOsM34XGqw[48;5;52m=[0m[48;5;52m            [0m[34m│ [48;5;28;38;5;231m5SmvrXHrbBmKQYwG3n/uiNN3kj5r56m3kZ5630msPA[48;5;22m=[0m[48;5;22m[0K[0m
[34m[38;5;231mgithub.com/xeipuuv/gojsonpointer v0.0.0-2019090519[34m↵[0m  [34m│ [38;5;231mgithub.com/xeipuuv/gojsonpointer v0.0.0-2019090519[34m↵[0m
[34m[38;5;231m4746-02993c407bfb/go.mod h1:N2zxlSyiKSe5eX1tZViRH5[34m↵[0m  [34m│ [38;5;231m4746-02993c407bfb/go.mod h1:N2zxlSyiKSe5eX1tZViRH5[34m↵[0m
[34m[38;5;231mQA0qijqEDrYZiPEAiq3wU=[0m                               [34m│ [38;5;231mQA0qijqEDrYZiPEAiq3wU=[0m
[34m[38;5;231mgithub.com/xeipuuv/gojsonreference v0.0.0-20180127[34m↵[0m  [34m│ [38;5;231mgithub.com/xeipuuv/gojsonreference v0.0.0-20180127[34m↵[0m
[34m[38;5;231m040603-bd5ef7bd5415/go.mod h1:GwrjFmJcFw6At/Gs6z4y[34m↵[0m  [34m│ [38;5;231m040603-bd5ef7bd5415/go.mod h1:GwrjFmJcFw6At/Gs6z4y[34m↵[0m
[34m[38;5;231mjiIwzuJ1/+UwLxMQDVQXShQ=[0m                             [34m│ [38;5;231mjiIwzuJ1/+UwLxMQDVQXShQ=[0m
[34m[38;5;231mgithub.com/xeipuuv/gojsonschema v1.2.0/go.mod h1:a[34m↵[0m  [34m│ [38;5;231mgithub.com/xeipuuv/gojsonschema v1.2.0/go.mod h1:a[34m↵[0m
[34m[38;5;231mnYRn/JVcOK2ZgGU+IjEV4nwlhoK5sQluxsYJ78Id3Y=[0m          [34m│ [38;5;231mnYRn/JVcOK2ZgGU+IjEV4nwlhoK5sQluxsYJ78Id3Y=[0m

[36m──────────────────────────────────────────────────────────────────────────────────────[0m[36m┐[0m
[34m619[0m:[38;5;231m go.opentelemetry.io/otel/trace v1.11.1 h1:ofxdnzsNrGBYXbP7t7zpUK281+go5rF7dvdIZX [0m[36m│[0m
[36m──────────────────────────────────────────────────────────────────────────────────────[0m[36m┴[0m[36m───────────────────[0m
[34m[38;5;231mgo.opentelemetry.io/otel/trace v1.11.1/go.mod h1:f[34m↵[0m  [34m│ [38;5;231mgo.opentelemetry.io/otel/trace v1.11.1/go.mod h1:f[34m↵[0m
[34m[38;5;231m/Q9G7vzk5u91PhbmKbg1Qn0rzH1LJ4vbPHFGkTPtOk=[0m          [34m│ [38;5;231m/Q9G7vzk5u91PhbmKbg1Qn0rzH1LJ4vbPHFGkTPtOk=[0m
[34m[38;5;231mgo.opentelemetry.io/proto/otlp v0.19.0 h1:IVN6GR+m[34m↵[0m  [34m│ [38;5;231mgo.opentelemetry.io/proto/otlp v0.19.0 h1:IVN6GR+m[34m↵[0m
[34m[38;5;231mhC4s5yfcTbmzHYODqvWAp3ZedA2SJPI1Nnw=[0m                 [34m│ [38;5;231mhC4s5yfcTbmzHYODqvWAp3ZedA2SJPI1Nnw=[0m
[34m[38;5;231mgo.opentelemetry.io/proto/otlp v0.19.0/go.mod h1:H[34m↵[0m  [34m│ [38;5;231mgo.opentelemetry.io/proto/otlp v0.19.0/go.mod h1:H[34m↵[0m
[34m[38;5;231m7XAot3MsfNsj7EXtrA2q5xSNQ10UqI405h3+duxN4U=[0m          [34m│ [38;5;231m7XAot3MsfNsj7EXtrA2q5xSNQ10UqI405h3+duxN4U=[0m
[34m[48;5;52;38;5;231mgo.starlark.net v0.0.0-20230112144946-fae38c8a6d89 h[34m↵[0m[34m│ [0m
[34m[48;5;52;38;5;231m1:qNFW0Bm9gXeA/h8lIzOiqvx7cMs/Xz5fkMgJpOo89qI=[0m[48;5;52m       [0m[34m│ [0m
[34m[48;5;52;38;5;231mgo.starlark.net v0.0.0-20230112144946-fae38c8a6d89/g[34m↵[0m[34m│ [0m
[34m[48;5;52;38;5;231mo.mod h1:kIVgS18CjmEC3PqMd5kaJSGEifyV/CeB9x506ZJ1Vbk=[0m[34m│ [0m
[34m[0m                                                     [34m│ [48;5;22;38;5;231mgo.starlark.net v0.0.0-20230122040757-066229b0515d[34m↵[0m[48;5;22m[0K[0m
[34m[0m                                                     [34m│ [48;5;22;38;5;231m h1:CTI+tbxvlfu7QlBj+4QjF8YPHoDh71h0/l2tXOM2k0o=[0m[48;5;22m[0K[0m
[34m[0m                                                     [34m│ [48;5;22;38;5;231mgo.starlark.net v0.0.0-20230122040757-066229b0515d[34m↵[0m[48;5;22m[0K[0m
[34m[0m                                                     [34m│ [48;5;22;38;5;231m/go.mod h1:1NtVfE+l6AHFaY4GmUPGHeLIW8/THkXnym5iweV[34m↵[0m[48;5;22m[0K[0m
[34m[0m                                                     [34m│ [48;5;22;38;5;231mgCwU=[0m[48;5;22m[0K[0m
[34m[38;5;231mgo.uber.org/atomic v1.10.0 h1:9qC72Qh0+3MqyJbAn8YU[34m↵[0m  [34m│ [38;5;231mgo.uber.org/atomic v1.10.0 h1:9qC72Qh0+3MqyJbAn8YU[34m↵[0m
[34m[38;5;231m5xVq1frD8bn3JtD2oXtafVQ=[0m                             [34m│ [38;5;231m5xVq1frD8bn3JtD2oXtafVQ=[0m
[34m[38;5;231mgo.uber.org/atomic v1.10.0/go.mod h1:LUxbIzbOniOlM[34m↵[0m  [34m│ [38;5;231mgo.uber.org/atomic v1.10.0/go.mod h1:LUxbIzbOniOlM[34m↵[0m
[34m[38;5;231mKjJjyPfpl4v+PKK2cNJn91OQbhoJI0=[0m                      [34m│ [38;5;231mKjJjyPfpl4v+PKK2cNJn91OQbhoJI0=[0m
[34m[38;5;231mgo.uber.org/automaxprocs v1.5.1 h1:e1YG66Lrk73dn4q[34m↵[0m  [34m│ [38;5;231mgo.uber.org/automaxprocs v1.5.1 h1:e1YG66Lrk73dn4q[34m↵[0m
[34m[38;5;231mhg8WFSvhF0JuFQF0ERIp4rpuV8Qk=[0m                        [34m│ [38;5;231mhg8WFSvhF0JuFQF0ERIp4rpuV8Qk=[0m

[36m─────────────────────────────────────────────────────────────────────────[0m[36m┐[0m
[34m639[0m:[38;5;231m gocloud.dev v0.28.0 h1:PjL1f9zu8epY1pFCIHdrQnJRZzRcDyAr18hNTkXIKlQ= [0m[36m│[0m
[36m─────────────────────────────────────────────────────────────────────────[0m[36m┴[0m[36m────────────────────────────────[0m
[34m[38;5;231mgocloud.dev v0.28.0/go.mod h1:nzSs01FpRYyIb/OqXLNN[34m↵[0m  [34m│ [38;5;231mgocloud.dev v0.28.0/go.mod h1:nzSs01FpRYyIb/OqXLNN[34m↵[0m
[34m[38;5;231ma+NMPZG9CdTUY/pGLgSpIN0=[0m                             [34m│ [38;5;231ma+NMPZG9CdTUY/pGLgSpIN0=[0m
[34m[38;5;231mgolang.org/x/crypto v0.5.0 h1:U/0M97KRkSFvyD/3FSmd[34m↵[0m  [34m│ [38;5;231mgolang.org/x/crypto v0.5.0 h1:U/0M97KRkSFvyD/3FSmd[34m↵[0m
[34m[38;5;231mP5W5swImpNgle/EHFhOsQPE=[0m                             [34m│ [38;5;231mP5W5swImpNgle/EHFhOsQPE=[0m
[34m[38;5;231mgolang.org/x/crypto v0.5.0/go.mod h1:NK/OQwhpMQP3M[34m↵[0m  [34m│ [38;5;231mgolang.org/x/crypto v0.5.0/go.mod h1:NK/OQwhpMQP3M[34m↵[0m
[34m[38;5;231mwtdjgLlYHnH9ebylxKWv3e0fK+mkQU=[0m                      [34m│ [38;5;231mwtdjgLlYHnH9ebylxKWv3e0fK+mkQU=[0m
[34m[48;5;52;38;5;231mgolang.org/x/exp v0.0.0-20230116083435-1de6713980de [34m↵[0m[34m│ [0m
[34m[48;5;52;38;5;231mh1:DBWn//IJw30uYCgERoxCg84hWtA97F4wMiKOIh00Uf0=[0m[48;5;52m      [0m[34m│ [0m
[34m[0m                                                     [34m│ [48;5;22;38;5;231mgolang.org/x/exp v0.0.0-20230118134722-a68e582fa15[34m↵[0m[48;5;22m[0K[0m
[34m[0m                                                     [34m│ [48;5;22;38;5;231m7 h1:fiNkyhJPUvxbRPbCqY/D9qdjmPzfHcpK3P4bM4gioSY=[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mgolang.org/x/exp v0.0.0-[48;5;124m20230116083435[48;5;52m-[48;5;124m1de6713980de[48;5;52m/[34m↵[0m[34m│ [48;5;22;38;5;231mgolang.org/x/exp v0.0.0-[48;5;28m20230118134722[48;5;22m-[48;5;28ma68e582fa15[48;5;22;34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mgo.mod h1:CxIveKay+FTh1D0yPZemJVgC/95VzuuOLq5Qi4xnoY[34m↵[0m[34m│ [48;5;28;38;5;231m7[48;5;22m/go.mod h1:CxIveKay+FTh1D0yPZemJVgC/95VzuuOLq5Qi4[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mc=[0m[48;5;52m                                                   [0m[34m│ [48;5;22;38;5;231mxnoYc=[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mgolang.org/x/exp/typeparams v0.0.0-20230116083435-1d[34m↵[0m[34m│ [0m
[34m[48;5;52;38;5;231me6713980de h1:aPL/oParTf1KrgrBJAeS7OcWqgQDbDVhSjq6rS[34m↵[0m[34m│ [0m
[34m[48;5;52;38;5;231mEEsQE=[0m[48;5;52m                                               [0m[34m│ [0m
[34m[0m                                                     [34m│ [48;5;22;38;5;231mgolang.org/x/exp/typeparams v0.0.0-20230118134722-[34m↵[0m[48;5;22m[0K[0m
[34m[0m                                                     [34m│ [48;5;22;38;5;231ma68e582fa157 h1:BKmw9kHvJFeDya3z07CXNwKhiReL0LnalM[34m↵[0m[48;5;22m[0K[0m
[34m[0m                                                     [34m│ [48;5;22;38;5;231mGf2B16dnM=[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mgolang.org/x/exp/typeparams v0.0.0-[48;5;124m20230116083435[48;5;52m-[48;5;124m1d[48;5;52;34m↵[0m[34m│ [48;5;22;38;5;231mgolang.org/x/exp/typeparams v0.0.0-[48;5;28m20230118134722[48;5;22m-[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;124;38;5;231me6713980de[48;5;52m/go.mod h1:AbB0pIl9nAr9wVwH+Z2ZpaocVmF5I4G[34m↵[0m[34m│ [48;5;28;38;5;231ma68e582fa157[48;5;22m/go.mod h1:AbB0pIl9nAr9wVwH+Z2ZpaocVmF[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231myWCDIsVjR0bk=[0m[48;5;52m                                        [0m[34m│ [48;5;22;38;5;231m5I4GyWCDIsVjR0bk=[0m[48;5;22m[0K[0m
[34m[38;5;231mgolang.org/x/image v0.3.0 h1:HTDXbdK9bjfSWkPzDJIw8[34m↵[0m  [34m│ [38;5;231mgolang.org/x/image v0.3.0 h1:HTDXbdK9bjfSWkPzDJIw8[34m↵[0m
[34m[38;5;231m9W8CAtfFGduujWs33NLLsg=[0m                              [34m│ [38;5;231m9W8CAtfFGduujWs33NLLsg=[0m
[34m[38;5;231mgolang.org/x/image v0.3.0/go.mod h1:fXd9211C/0VTlY[34m↵[0m  [34m│ [38;5;231mgolang.org/x/image v0.3.0/go.mod h1:fXd9211C/0VTlY[34m↵[0m
[34m[38;5;231muAcOhW8dY/RtEJqODXOWBDpmYBf+A=[0m                       [34m│ [38;5;231muAcOhW8dY/RtEJqODXOWBDpmYBf+A=[0m
[34m[38;5;231mgolang.org/x/lint v0.0.0-20210508222113-6edffad5e6[34m↵[0m  [34m│ [38;5;231mgolang.org/x/lint v0.0.0-20210508222113-6edffad5e6[34m↵[0m
[34m[38;5;231m16 h1:VLliZ0d+/avPrXXH+OakdXhpJuEoBZuwh1m2j7U6Iug=[0m   [34m│ [38;5;231m16 h1:VLliZ0d+/avPrXXH+OakdXhpJuEoBZuwh1m2j7U6Iug=[0m

[36m──────────────────────────────────────────────────────────────────────────────────────[0m[36m┐[0m
[34m674[0m:[38;5;231m gonum.org/v1/hdf5 v0.0.0-20210714002203-8c5d23bc6946 h1:vJpL69PeUullhJyKtTjHjENE [0m[36m│[0m
[36m──────────────────────────────────────────────────────────────────────────────────────[0m[36m┴[0m[36m───────────────────[0m
[34m[38;5;231mgonum.org/v1/hdf5 v0.0.0-20210714002203-8c5d23bc69[34m↵[0m  [34m│ [38;5;231mgonum.org/v1/hdf5 v0.0.0-20210714002203-8c5d23bc69[34m↵[0m
[34m[38;5;231m46/go.mod h1:BQUWDHIAygjdt1HnUPQ0eWqLN2n5FwJycrpYU[34m↵[0m  [34m│ [38;5;231m46/go.mod h1:BQUWDHIAygjdt1HnUPQ0eWqLN2n5FwJycrpYU[34m↵[0m
[34m[38;5;231mVUOx2I=[0m                                              [34m│ [38;5;231mVUOx2I=[0m
[34m[38;5;231mgonum.org/v1/plot v0.12.0 h1:y1ZNmfz/xHuHvtgFe8USZ[34m↵[0m  [34m│ [38;5;231mgonum.org/v1/plot v0.12.0 h1:y1ZNmfz/xHuHvtgFe8USZ[34m↵[0m
[34m[38;5;231mVyykQo5ERXPnspQNVK15Og=[0m                              [34m│ [38;5;231mVyykQo5ERXPnspQNVK15Og=[0m
[34m[38;5;231mgonum.org/v1/plot v0.12.0/go.mod h1:PgiMf9+3A3PnZd[34m↵[0m  [34m│ [38;5;231mgonum.org/v1/plot v0.12.0/go.mod h1:PgiMf9+3A3PnZd[34m↵[0m
[34m[38;5;231mJIciIXmyN1FwdAA6rXELSN761oQkw=[0m                       [34m│ [38;5;231mJIciIXmyN1FwdAA6rXELSN761oQkw=[0m
[34m[48;5;52;38;5;231mgoogle.golang.org/api v0.[48;5;124m107[48;5;52m.0 h1:[48;5;124mI2SlFjD8ZWabaIFOfe[48;5;52;34m↵[0m[34m│ [48;5;22;38;5;231mgoogle.golang.org/api v0.[48;5;28m108[48;5;22m.0 h1:[48;5;28mWVBc[48;5;22m/[48;5;28mfaN0DkKtR43[48;5;22;34m↵[0m[48;5;22m[0K[0m
[34m[48;5;124;38;5;231mEDg3pf0BHJDh6iYQ1ic3Yu[48;5;52m/[48;5;124mUU[48;5;52m=[0m[48;5;52m                           [0m[34m│ [48;5;28;38;5;231mQ/7+tPny9ZoLZdIiAyG5Q9vFClg[48;5;22m=[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mgoogle.golang.org/api v0.[48;5;124m107[48;5;52m.0/go.mod h1:2Ts0XTHNVWx[34m↵[0m[34m│ [48;5;22;38;5;231mgoogle.golang.org/api v0.[48;5;28m108[48;5;22m.0/go.mod h1:2Ts0XTHNV[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mypznxWOYUeI4g3WdP9Pk2Qk58+a/O9MY=[0m[48;5;52m                    [0m[34m│ [48;5;22;38;5;231mWxypznxWOYUeI4g3WdP9Pk2Qk58+a/O9MY=[0m[48;5;22m[0K[0m
[34m[38;5;231mgoogle.golang.org/appengine v1.6.7 h1:FZR1q0exgwxz[34m↵[0m  [34m│ [38;5;231mgoogle.golang.org/appengine v1.6.7 h1:FZR1q0exgwxz[34m↵[0m
[34m[38;5;231mPzp/aF+VccGrSfxfPpkBqjIIEq3ru6c=[0m                     [34m│ [38;5;231mPzp/aF+VccGrSfxfPpkBqjIIEq3ru6c=[0m
[34m[38;5;231mgoogle.golang.org/appengine v1.6.7/go.mod h1:8WjMM[34m↵[0m  [34m│ [38;5;231mgoogle.golang.org/appengine v1.6.7/go.mod h1:8WjMM[34m↵[0m
[34m[38;5;231mxjGQR8xUklV/ARdw2HLXBOI7O7uCIDZVag1xfc=[0m              [34m│ [38;5;231mxjGQR8xUklV/ARdw2HLXBOI7O7uCIDZVag1xfc=[0m
[34m[48;5;52;38;5;231mgoogle.golang.org/genproto v0.0.0-20230117162540-28d[34m↵[0m[34m│ [0m
[34m[48;5;52;38;5;231m6b9783ac4 h1:yF0uHwqqYt2tIL2F4hxRWA1ZFX43SEunWAK8MnQ[34m↵[0m[34m│ [0m
[34m[48;5;52;38;5;231miclk=[0m[48;5;52m                                                [0m[34m│ [0m
[34m[0m                                                     [34m│ [48;5;22;38;5;231mgoogle.golang.org/genproto v0.0.0-20230123190316-2[34m↵[0m[48;5;22m[0K[0m
[34m[0m                                                     [34m│ [48;5;22;38;5;231mc411cf9d197 h1:BwjeHhu4HS48EZmu1nS7flldBIDPC3qn+Hq[34m↵[0m[48;5;22m[0K[0m
[34m[0m                                                     [34m│ [48;5;22;38;5;231maSQ1K4x8=[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mgoogle.golang.org/genproto v0.0.0-[48;5;124m20230117162540[48;5;52m-[48;5;124m28d[48;5;52;34m↵[0m[34m│ [48;5;22;38;5;231mgoogle.golang.org/genproto v0.0.0-[48;5;28m20230123190316[48;5;22m-[48;5;28m2[48;5;22;34m↵[0m[48;5;22m[0K[0m
[34m[48;5;124;38;5;231m6b9783ac4[48;5;52m/go.mod h1:RGgjbofJ8xD9Sq1VVhDM1Vok1vRONV+r[34m↵[0m[34m│ [48;5;28;38;5;231mc411cf9d197[48;5;22m/go.mod h1:RGgjbofJ8xD9Sq1VVhDM1Vok1vRO[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mg+CjzG4SZKM=[0m[48;5;52m                                         [0m[34m│ [48;5;22;38;5;231mNV+rg+CjzG4SZKM=[0m[48;5;22m[0K[0m
[34m[38;5;231mgoogle.golang.org/grpc v1.52.0 h1:kd48UiU7EHsV4rnL[34m↵[0m  [34m│ [38;5;231mgoogle.golang.org/grpc v1.52.0 h1:kd48UiU7EHsV4rnL[34m↵[0m
[34m[38;5;231myOJRuP/Il/UHE7gdDAQ+SZI7nZk=[0m                         [34m│ [38;5;231myOJRuP/Il/UHE7gdDAQ+SZI7nZk=[0m
[34m[38;5;231mgoogle.golang.org/grpc v1.52.0/go.mod h1:pu6fVzoFb[34m↵[0m  [34m│ [38;5;231mgoogle.golang.org/grpc v1.52.0/go.mod h1:pu6fVzoFb[34m↵[0m
[34m[38;5;231m+NBYNAvQL08ic+lvB2IojljRYuun5vorUY=[0m                  [34m│ [38;5;231m+NBYNAvQL08ic+lvB2IojljRYuun5vorUY=[0m
[34m[38;5;231mgoogle.golang.org/grpc/cmd/protoc-gen-go-grpc v1.2[34m↵[0m  [34m│ [38;5;231mgoogle.golang.org/grpc/cmd/protoc-gen-go-grpc v1.2[34m↵[0m
[34m[38;5;231m.0/go.mod h1:DNq5QpG7LJqD2AamLZ7zvKE0DEpVl2BSEVjFy[34m↵[0m  [34m│ [38;5;231m.0/go.mod h1:DNq5QpG7LJqD2AamLZ7zvKE0DEpVl2BSEVjFy[34m↵[0m
[34m[38;5;231mcAAjRY=[0m                                              [34m│ [38;5;231mcAAjRY=[0m

[36m──────────────────────────────────────────────────────────────────────────────────────[0m[36m┐[0m
[34m687[0m:[38;5;231m google.golang.org/protobuf v1.28.1 h1:d0NfwRgPtno5B1Wa6L2DAG+KivqkdutMf1UhdNx175 [0m[36m│[0m
[36m──────────────────────────────────────────────────────────────────────────────────────[0m[36m┴[0m[36m───────────────────[0m
[34m[38;5;231mgoogle.golang.org/protobuf v1.28.1/go.mod h1:HV8QO[34m↵[0m  [34m│ [38;5;231mgoogle.golang.org/protobuf v1.28.1/go.mod h1:HV8QO[34m↵[0m
[34m[38;5;231md/L58Z+nl8r43ehVNZIU/HEI6OcFqwMG9pJV4I=[0m              [34m│ [38;5;231md/L58Z+nl8r43ehVNZIU/HEI6OcFqwMG9pJV4I=[0m
[34m[38;5;231mgopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c[34m↵[0m  [34m│ [38;5;231mgopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c[34m↵[0m
[34m[38;5;231m6c h1:Hei/4ADfdWqJk1ZMxUNpqntNwaWcugrBjAiHlqqRiVk=[0m   [34m│ [38;5;231m6c h1:Hei/4ADfdWqJk1ZMxUNpqntNwaWcugrBjAiHlqqRiVk=[0m
[34m[38;5;231mgopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c[34m↵[0m  [34m│ [38;5;231mgopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c[34m↵[0m
[34m[38;5;231m6c/go.mod h1:JHkPIbrfpd72SG/EVd6muEfDQjcINNoR0C8j2[34m↵[0m  [34m│ [38;5;231m6c/go.mod h1:JHkPIbrfpd72SG/EVd6muEfDQjcINNoR0C8j2[34m↵[0m
[34m[38;5;231mr3qZ4Q=[0m                                              [34m│ [38;5;231mr3qZ4Q=[0m
[34m[48;5;52;38;5;231mgopkg.in/inconshreveable/log15.v2 v2.[48;5;124m0[48;5;52m.0[48;5;124m-20221122034[48;5;52;34m↵[0m[34m│ [48;5;22;38;5;231mgopkg.in/inconshreveable/log15.v2 v2.[48;5;28m16[48;5;22m.0/go.mod h[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;124;38;5;231m931-555555054819[48;5;52m/go.mod h1:aPpfJ7XW+gOuirDoZ8gHhLh3k[34m↵[0m[34m│ [48;5;22;38;5;231m1:aPpfJ7XW+gOuirDoZ8gHhLh3kZ1B08FtV2bbmy7Jv3s=[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mZ1B08FtV2bbmy7Jv3s=[0m[48;5;52m                                  [0m[34m│ [0m
[34m[38;5;231mgopkg.in/inf.v0 v0.9.1 h1:73M5CoZyi3ZLMOyDlQh031Cx[34m↵[0m  [34m│ [38;5;231mgopkg.in/inf.v0 v0.9.1 h1:73M5CoZyi3ZLMOyDlQh031Cx[34m↵[0m
[34m[38;5;231m6N9NDJ2Vvfl76EDAgDc=[0m                                 [34m│ [38;5;231m6N9NDJ2Vvfl76EDAgDc=[0m
[34m[38;5;231mgopkg.in/inf.v0 v0.9.1/go.mod h1:cWUDdTG/fYaXco+Dc[34m↵[0m  [34m│ [38;5;231mgopkg.in/inf.v0 v0.9.1/go.mod h1:cWUDdTG/fYaXco+Dc[34m↵[0m
[34m[38;5;231mufb5Vnc6Gp2YChqWtbxRZE0mXw=[0m                          [34m│ [38;5;231mufb5Vnc6Gp2YChqWtbxRZE0mXw=[0m
[34m[38;5;231mgopkg.in/natefinch/npipe.v2 v2.0.0-20160621034901-[34m↵[0m  [34m│ [38;5;231mgopkg.in/natefinch/npipe.v2 v2.0.0-20160621034901-[34m↵[0m
[34m[38;5;231mc1b8fa8bdcce/go.mod h1:5AcXVHNjg+BDxry382+8OKon8SE[34m↵[0m  [34m│ [38;5;231mc1b8fa8bdcce/go.mod h1:5AcXVHNjg+BDxry382+8OKon8SE[34m↵[0m
[34m[38;5;231mWiKktQR07RKPsv1c=[0m                                    [34m│ [38;5;231mWiKktQR07RKPsv1c=[0m

[36m──────────────────────────────────────────────────────────────────────────────[0m[36m┐[0m
[34m712[0m:[38;5;231m k8s.io/client-go v0.26.0 h1:lT1D3OfO+wIi9UFolCrifbjUUgu7CpLca0AD8ghRLI8= [0m[36m│[0m
[36m──────────────────────────────────────────────────────────────────────────────[0m[36m┴[0m[36m───────────────────────────[0m
[34m[38;5;231mk8s.io/client-go v0.26.0/go.mod h1:I2Sh57A79EQsDmn[34m↵[0m  [34m│ [38;5;231mk8s.io/client-go v0.26.0/go.mod h1:I2Sh57A79EQsDmn[34m↵[0m
[34m[38;5;231m7F7ASpmru1cceh3ocVT9KlX2jEZg=[0m                        [34m│ [38;5;231m7F7ASpmru1cceh3ocVT9KlX2jEZg=[0m
[34m[38;5;231mk8s.io/component-base v0.26.0 h1:0IkChOCohtDHttmKu[34m↵[0m  [34m│ [38;5;231mk8s.io/component-base v0.26.0 h1:0IkChOCohtDHttmKu[34m↵[0m
[34m[38;5;231mz+EP3j3+qKmV55rM9gIFTXA7Vs=[0m                          [34m│ [38;5;231mz+EP3j3+qKmV55rM9gIFTXA7Vs=[0m
[34m[38;5;231mk8s.io/component-base v0.26.0/go.mod h1:lqHwlfV1/h[34m↵[0m  [34m│ [38;5;231mk8s.io/component-base v0.26.0/go.mod h1:lqHwlfV1/h[34m↵[0m
[34m[38;5;231maa14F/Z5Zizk5QmzaVf23nQzCwVOQpfC8=[0m                   [34m│ [38;5;231maa14F/Z5Zizk5QmzaVf23nQzCwVOQpfC8=[0m
[34m[48;5;52;38;5;231mk8s.io/klog/v2 v2.80.1 h1:atnLQ121W371wYYFawwYx1aEY2[34m↴[0m[34m│ [0m
[34m[0m[48;5;52m                                  [34m…[38;5;231meUfs4l3J72wtgAwV4=[0m[34m│ [0m
[34m[0m                                                     [34m│ [48;5;22;38;5;231mk8s.io/klog/v2 v2.90.0 h1:VkTxIV/FjRXn1fgNNcKGM8cf[34m↵[0m[48;5;22m[0K[0m
[34m[0m                                                     [34m│ [48;5;22;38;5;231mmL1Z33ZjXRTVxKCoF5M=[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mk8s.io/klog/v2 v2.[48;5;124m80[48;5;52m.[48;5;124m1[48;5;52m/go.mod h1:y1WjHnz7Dj687irZUWR[34m↵[0m[34m│ [48;5;22;38;5;231mk8s.io/klog/v2 v2.[48;5;28m90[48;5;22m.[48;5;28m0[48;5;22m/go.mod h1:y1WjHnz7Dj687irZU[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231m/WLkLc5N1YHtjLdmgWjndZn0=[0m[48;5;52m                            [0m[34m│ [48;5;22;38;5;231mWR/WLkLc5N1YHtjLdmgWjndZn0=[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mk8s.io/kube-openapi v0.0.0-20230117224833-444ee995c1[34m↵[0m[34m│ [0m
[34m[48;5;52;38;5;231m20 h1:bB/6AuV41SUB7Qm9fJaX/RRSZLo2ft9KRQilTD+VGEg=[0m[48;5;52m   [0m[34m│ [0m
[34m[0m                                                     [34m│ [48;5;22;38;5;231mk8s.io/kube-openapi v0.0.0-20230123231816-1cb3ae25[34m↵[0m[48;5;22m[0K[0m
[34m[0m                                                     [34m│ [48;5;22;38;5;231md79a h1:s6zvHjyDQX1NtVT88pvw2tddqhqY0Bz0Gbnn+yctsF[34m↵[0m[48;5;22m[0K[0m
[34m[0m                                                     [34m│ [48;5;22;38;5;231mU=[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mk8s.io/kube-openapi v0.0.0-[48;5;124m20230117224833[48;5;52m-[48;5;124m444ee995c1[48;5;52;34m↵[0m[34m│ [48;5;22;38;5;231mk8s.io/kube-openapi v0.0.0-[48;5;28m20230123231816[48;5;22m-[48;5;28m1cb3ae25[48;5;22;34m↵[0m[48;5;22m[0K[0m
[34m[48;5;124;38;5;231m20[48;5;52m/go.mod h1:/BYxry62FuDzmI+i9B+X2pqfySRmSOW2ARmj5Zb[34m↵[0m[34m│ [48;5;28;38;5;231md79a[48;5;22m/go.mod h1:/BYxry62FuDzmI+i9B+X2pqfySRmSOW2ARm[34m↵[0m[48;5;22m[0K[0m
[34m[48;5;52;38;5;231mqhj0=[0m[48;5;52m                                                [0m[34m│ [48;5;22;38;5;231mj5Zbqhj0=[0m[48;5;22m[0K[0m
[34m[38;5;231mk8s.io/metrics v0.26.0 h1:U/NzZHKDrIVGL93AUMRkqqXj[34m↵[0m  [34m│ [38;5;231mk8s.io/metrics v0.26.0 h1:U/NzZHKDrIVGL93AUMRkqqXj[34m↵[0m
[34m[38;5;231mOah3wGvjSnKmG/5NVCs=[0m                                 [34m│ [38;5;231mOah3wGvjSnKmG/5NVCs=[0m
[34m[38;5;231mk8s.io/metrics v0.26.0/go.mod h1:cf5MlG4ZgWaEFZrR9[34m↵[0m  [34m│ [38;5;231mk8s.io/metrics v0.26.0/go.mod h1:cf5MlG4ZgWaEFZrR9[34m↵[0m
[34m[38;5;231m+sOImhZ2ICMpIdNurA+D8snIs8=[0m                          [34m│ [38;5;231m+sOImhZ2ICMpIdNurA+D8snIs8=[0m
[34m[38;5;231mk8s.io/utils v0.0.0-20221128185143-99ec85e7a448 h1[34m↵[0m  [34m│ [38;5;231mk8s.io/utils v0.0.0-20221128185143-99ec85e7a448 h1[34m↵[0m
[34m[38;5;231m:KTgPnR10d5zhztWptI952TNtt/4u5h3IzDXkdIMuo2Y=[0m        [34m│ [38;5;231m:KTgPnR10d5zhztWptI952TNtt/4u5h3IzDXkdIMuo2Y=[0m

[1;4;33mk8s/agent/configmap.yaml[0m

[36m──────────────[0m[36m┐[0m
[34m20[0m:[38;5;231m [38;5;203mmetadata[38;5;231m: [0m[36m│[0m
[36m──────────────[0m[36m┴[0m[36m───────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[38;5;231m  [38;5;203mname[38;5;231m: [38;5;186mvald-agent-ngt-config[0m                        [34m│ [38;5;231m  [38;5;203mname[38;5;231m: [38;5;186mvald-agent-ngt-config[0m
[34m[38;5;231m  [38;5;203mlabels[38;5;231m:[0m                                            [34m│ [38;5;231m  [38;5;203mlabels[38;5;231m:[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/name[38;5;231m: [38;5;186mvald[0m                     [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/name[38;5;231m: [38;5;186mvald[0m
[34m[48;5;52;38;5;231m    [38;5;203mhelm.sh/chart[38;5;231m: [38;5;186mvald-v1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                       [0m[34m│ [48;5;22;38;5;231m    [38;5;203mhelm.sh/chart[38;5;231m: [38;5;186mvald-v1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/managed-by[38;5;231m: [38;5;186mHelm[0m               [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/managed-by[38;5;231m: [38;5;186mHelm[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/instance[38;5;231m: [38;5;186mrelease-name[0m         [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/instance[38;5;231m: [38;5;186mrelease-name[0m
[34m[48;5;52;38;5;231m    [38;5;203mapp.kubernetes.io/version[38;5;231m: [38;5;186mv1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                [0m[34m│ [48;5;22;38;5;231m    [38;5;203mapp.kubernetes.io/version[38;5;231m: [38;5;186mv1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/component[38;5;231m: [38;5;186magent[0m               [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/component[38;5;231m: [38;5;186magent[0m
[34m[38;5;203mdata[38;5;231m:[0m                                                [34m│ [38;5;203mdata[38;5;231m:[0m
[34m[38;5;231m  [38;5;203mconfig.yaml[38;5;231m: [38;5;203m|[0m                                     [34m│ [38;5;231m  [38;5;203mconfig.yaml[38;5;231m: [38;5;203m|[0m

[1;4;33mk8s/agent/pdb.yaml[0m

[36m──────────────[0m[36m┐[0m
[34m20[0m:[38;5;231m [38;5;203mmetadata[38;5;231m: [0m[36m│[0m
[36m──────────────[0m[36m┴[0m[36m───────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[38;5;231m  [38;5;203mname[38;5;231m: [38;5;186mvald-agent-ngt[0m                               [34m│ [38;5;231m  [38;5;203mname[38;5;231m: [38;5;186mvald-agent-ngt[0m
[34m[38;5;231m  [38;5;203mlabels[38;5;231m:[0m                                            [34m│ [38;5;231m  [38;5;203mlabels[38;5;231m:[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/name[38;5;231m: [38;5;186mvald[0m                     [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/name[38;5;231m: [38;5;186mvald[0m
[34m[48;5;52;38;5;231m    [38;5;203mhelm.sh/chart[38;5;231m: [38;5;186mvald-v1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                       [0m[34m│ [48;5;22;38;5;231m    [38;5;203mhelm.sh/chart[38;5;231m: [38;5;186mvald-v1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/managed-by[38;5;231m: [38;5;186mHelm[0m               [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/managed-by[38;5;231m: [38;5;186mHelm[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/instance[38;5;231m: [38;5;186mrelease-name[0m         [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/instance[38;5;231m: [38;5;186mrelease-name[0m
[34m[48;5;52;38;5;231m    [38;5;203mapp.kubernetes.io/version[38;5;231m: [38;5;186mv1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                [0m[34m│ [48;5;22;38;5;231m    [38;5;203mapp.kubernetes.io/version[38;5;231m: [38;5;186mv1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/component[38;5;231m: [38;5;186magent[0m               [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/component[38;5;231m: [38;5;186magent[0m
[34m[38;5;203mspec[38;5;231m:[0m                                                [34m│ [38;5;203mspec[38;5;231m:[0m
[34m[38;5;231m  [38;5;203mmaxUnavailable[38;5;231m: [38;5;141m1[0m                                  [34m│ [38;5;231m  [38;5;203mmaxUnavailable[38;5;231m: [38;5;141m1[0m

[1;4;33mk8s/agent/priorityclass.yaml[0m

[36m──────────────[0m[36m┐[0m
[34m20[0m:[38;5;231m [38;5;203mmetadata[38;5;231m: [0m[36m│[0m
[36m──────────────[0m[36m┴[0m[36m───────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[38;5;231m  [38;5;203mname[38;5;231m: [38;5;186mdefault-vald-agent-ngt-priority[0m              [34m│ [38;5;231m  [38;5;203mname[38;5;231m: [38;5;186mdefault-vald-agent-ngt-priority[0m
[34m[38;5;231m  [38;5;203mlabels[38;5;231m:[0m                                            [34m│ [38;5;231m  [38;5;203mlabels[38;5;231m:[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/name[38;5;231m: [38;5;186mvald[0m                     [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/name[38;5;231m: [38;5;186mvald[0m
[34m[48;5;52;38;5;231m    [38;5;203mhelm.sh/chart[38;5;231m: [38;5;186mvald-v1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                       [0m[34m│ [48;5;22;38;5;231m    [38;5;203mhelm.sh/chart[38;5;231m: [38;5;186mvald-v1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/managed-by[38;5;231m: [38;5;186mHelm[0m               [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/managed-by[38;5;231m: [38;5;186mHelm[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/instance[38;5;231m: [38;5;186mrelease-name[0m         [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/instance[38;5;231m: [38;5;186mrelease-name[0m
[34m[48;5;52;38;5;231m    [38;5;203mapp.kubernetes.io/version[38;5;231m: [38;5;186mv1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                [0m[34m│ [48;5;22;38;5;231m    [38;5;203mapp.kubernetes.io/version[38;5;231m: [38;5;186mv1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/component[38;5;231m: [38;5;186magent[0m               [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/component[38;5;231m: [38;5;186magent[0m
[34m[38;5;203mvalue[38;5;231m: [38;5;186m1e+09[0m                                         [34m│ [38;5;203mvalue[38;5;231m: [38;5;186m1e+09[0m
[34m[38;5;203mpreemptionPolicy[38;5;231m: [38;5;186mNever[0m                              [34m│ [38;5;203mpreemptionPolicy[38;5;231m: [38;5;186mNever[0m

[1;4;33mk8s/agent/statefulset.yaml[0m

[36m──────────────[0m[36m┐[0m
[34m21[0m:[38;5;231m [38;5;203mmetadata[38;5;231m: [0m[36m│[0m
[36m──────────────[0m[36m┴[0m[36m───────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[38;5;231m  [38;5;203mlabels[38;5;231m:[0m                                            [34m│ [38;5;231m  [38;5;203mlabels[38;5;231m:[0m
[34m[38;5;231m    [38;5;203mapp[38;5;231m: [38;5;186mvald-agent-ngt[0m                              [34m│ [38;5;231m    [38;5;203mapp[38;5;231m: [38;5;186mvald-agent-ngt[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/name[38;5;231m: [38;5;186mvald[0m                     [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/name[38;5;231m: [38;5;186mvald[0m
[34m[48;5;52;38;5;231m    [38;5;203mhelm.sh/chart[38;5;231m: [38;5;186mvald-v1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                       [0m[34m│ [48;5;22;38;5;231m    [38;5;203mhelm.sh/chart[38;5;231m: [38;5;186mvald-v1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/managed-by[38;5;231m: [38;5;186mHelm[0m               [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/managed-by[38;5;231m: [38;5;186mHelm[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/instance[38;5;231m: [38;5;186mrelease-name[0m         [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/instance[38;5;231m: [38;5;186mrelease-name[0m
[34m[48;5;52;38;5;231m    [38;5;203mapp.kubernetes.io/version[38;5;231m: [38;5;186mv1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                [0m[34m│ [48;5;22;38;5;231m    [38;5;203mapp.kubernetes.io/version[38;5;231m: [38;5;186mv1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/component[38;5;231m: [38;5;186magent[0m               [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/component[38;5;231m: [38;5;186magent[0m
[34m[38;5;203mspec[38;5;231m:[0m                                                [34m│ [38;5;203mspec[38;5;231m:[0m
[34m[38;5;231m  [38;5;203mserviceName[38;5;231m: [38;5;186mvald-agent-ngt[0m                        [34m│ [38;5;231m  [38;5;203mserviceName[38;5;231m: [38;5;186mvald-agent-ngt[0m

[1;4;33mk8s/agent/svc.yaml[0m

[36m──────────────[0m[36m┐[0m
[34m20[0m:[38;5;231m [38;5;203mmetadata[38;5;231m: [0m[36m│[0m
[36m──────────────[0m[36m┴[0m[36m───────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[38;5;231m  [38;5;203mname[38;5;231m: [38;5;186mvald-agent-ngt[0m                               [34m│ [38;5;231m  [38;5;203mname[38;5;231m: [38;5;186mvald-agent-ngt[0m
[34m[38;5;231m  [38;5;203mlabels[38;5;231m:[0m                                            [34m│ [38;5;231m  [38;5;203mlabels[38;5;231m:[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/name[38;5;231m: [38;5;186mvald[0m                     [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/name[38;5;231m: [38;5;186mvald[0m
[34m[48;5;52;38;5;231m    [38;5;203mhelm.sh/chart[38;5;231m: [38;5;186mvald-v1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                       [0m[34m│ [48;5;22;38;5;231m    [38;5;203mhelm.sh/chart[38;5;231m: [38;5;186mvald-v1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/managed-by[38;5;231m: [38;5;186mHelm[0m               [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/managed-by[38;5;231m: [38;5;186mHelm[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/instance[38;5;231m: [38;5;186mrelease-name[0m         [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/instance[38;5;231m: [38;5;186mrelease-name[0m
[34m[48;5;52;38;5;231m    [38;5;203mapp.kubernetes.io/version[38;5;231m: [38;5;186mv1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                [0m[34m│ [48;5;22;38;5;231m    [38;5;203mapp.kubernetes.io/version[38;5;231m: [38;5;186mv1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/component[38;5;231m: [38;5;186magent[0m               [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/component[38;5;231m: [38;5;186magent[0m
[34m[38;5;203mspec[38;5;231m:[0m                                                [34m│ [38;5;203mspec[38;5;231m:[0m
[34m[38;5;231m  [38;5;203mports[38;5;231m:[0m                                             [34m│ [38;5;231m  [38;5;203mports[38;5;231m:[0m

[1;4;33mk8s/discoverer/clusterrole.yaml[0m

[36m──────────────[0m[36m┐[0m
[34m20[0m:[38;5;231m [38;5;203mmetadata[38;5;231m: [0m[36m│[0m
[36m──────────────[0m[36m┴[0m[36m───────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[38;5;231m  [38;5;203mname[38;5;231m: [38;5;186mdiscoverer[0m                                   [34m│ [38;5;231m  [38;5;203mname[38;5;231m: [38;5;186mdiscoverer[0m
[34m[38;5;231m  [38;5;203mlabels[38;5;231m:[0m                                            [34m│ [38;5;231m  [38;5;203mlabels[38;5;231m:[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/name[38;5;231m: [38;5;186mvald[0m                     [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/name[38;5;231m: [38;5;186mvald[0m
[34m[48;5;52;38;5;231m    [38;5;203mhelm.sh/chart[38;5;231m: [38;5;186mvald-v1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                       [0m[34m│ [48;5;22;38;5;231m    [38;5;203mhelm.sh/chart[38;5;231m: [38;5;186mvald-v1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/managed-by[38;5;231m: [38;5;186mHelm[0m               [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/managed-by[38;5;231m: [38;5;186mHelm[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/instance[38;5;231m: [38;5;186mrelease-name[0m         [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/instance[38;5;231m: [38;5;186mrelease-name[0m
[34m[48;5;52;38;5;231m    [38;5;203mapp.kubernetes.io/version[38;5;231m: [38;5;186mv1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                [0m[34m│ [48;5;22;38;5;231m    [38;5;203mapp.kubernetes.io/version[38;5;231m: [38;5;186mv1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/component[38;5;231m: [38;5;186mdiscoverer[0m          [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/component[38;5;231m: [38;5;186mdiscoverer[0m
[34m[38;5;203mrules[38;5;231m:[0m                                               [34m│ [38;5;203mrules[38;5;231m:[0m
[34m[38;5;231m  - [38;5;203mapiGroups[38;5;231m:[0m                                       [34m│ [38;5;231m  - [38;5;203mapiGroups[38;5;231m:[0m

[1;4;33mk8s/discoverer/clusterrolebinding.yaml[0m

[36m──────────────[0m[36m┐[0m
[34m20[0m:[38;5;231m [38;5;203mmetadata[38;5;231m: [0m[36m│[0m
[36m──────────────[0m[36m┴[0m[36m───────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[38;5;231m  [38;5;203mname[38;5;231m: [38;5;186mdiscoverer[0m                                   [34m│ [38;5;231m  [38;5;203mname[38;5;231m: [38;5;186mdiscoverer[0m
[34m[38;5;231m  [38;5;203mlabels[38;5;231m:[0m                                            [34m│ [38;5;231m  [38;5;203mlabels[38;5;231m:[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/name[38;5;231m: [38;5;186mvald[0m                     [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/name[38;5;231m: [38;5;186mvald[0m
[34m[48;5;52;38;5;231m    [38;5;203mhelm.sh/chart[38;5;231m: [38;5;186mvald-v1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                       [0m[34m│ [48;5;22;38;5;231m    [38;5;203mhelm.sh/chart[38;5;231m: [38;5;186mvald-v1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/managed-by[38;5;231m: [38;5;186mHelm[0m               [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/managed-by[38;5;231m: [38;5;186mHelm[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/instance[38;5;231m: [38;5;186mrelease-name[0m         [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/instance[38;5;231m: [38;5;186mrelease-name[0m
[34m[48;5;52;38;5;231m    [38;5;203mapp.kubernetes.io/version[38;5;231m: [38;5;186mv1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                [0m[34m│ [48;5;22;38;5;231m    [38;5;203mapp.kubernetes.io/version[38;5;231m: [38;5;186mv1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/component[38;5;231m: [38;5;186mdiscoverer[0m          [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/component[38;5;231m: [38;5;186mdiscoverer[0m
[34m[38;5;203mroleRef[38;5;231m:[0m                                             [34m│ [38;5;203mroleRef[38;5;231m:[0m
[34m[38;5;231m  [38;5;203mapiGroup[38;5;231m: [38;5;186mrbac.authorization.k8s.io[0m                [34m│ [38;5;231m  [38;5;203mapiGroup[38;5;231m: [38;5;186mrbac.authorization.k8s.io[0m

[1;4;33mk8s/discoverer/configmap.yaml[0m

[36m──────────────[0m[36m┐[0m
[34m20[0m:[38;5;231m [38;5;203mmetadata[38;5;231m: [0m[36m│[0m
[36m──────────────[0m[36m┴[0m[36m───────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[38;5;231m  [38;5;203mname[38;5;231m: [38;5;186mvald-discoverer-config[0m                       [34m│ [38;5;231m  [38;5;203mname[38;5;231m: [38;5;186mvald-discoverer-config[0m
[34m[38;5;231m  [38;5;203mlabels[38;5;231m:[0m                                            [34m│ [38;5;231m  [38;5;203mlabels[38;5;231m:[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/name[38;5;231m: [38;5;186mvald[0m                     [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/name[38;5;231m: [38;5;186mvald[0m
[34m[48;5;52;38;5;231m    [38;5;203mhelm.sh/chart[38;5;231m: [38;5;186mvald-v1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                       [0m[34m│ [48;5;22;38;5;231m    [38;5;203mhelm.sh/chart[38;5;231m: [38;5;186mvald-v1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/managed-by[38;5;231m: [38;5;186mHelm[0m               [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/managed-by[38;5;231m: [38;5;186mHelm[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/instance[38;5;231m: [38;5;186mrelease-name[0m         [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/instance[38;5;231m: [38;5;186mrelease-name[0m
[34m[48;5;52;38;5;231m    [38;5;203mapp.kubernetes.io/version[38;5;231m: [38;5;186mv1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                [0m[34m│ [48;5;22;38;5;231m    [38;5;203mapp.kubernetes.io/version[38;5;231m: [38;5;186mv1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/component[38;5;231m: [38;5;186mdiscoverer[0m          [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/component[38;5;231m: [38;5;186mdiscoverer[0m
[34m[38;5;203mdata[38;5;231m:[0m                                                [34m│ [38;5;203mdata[38;5;231m:[0m
[34m[38;5;231m  [38;5;203mconfig.yaml[38;5;231m: [38;5;203m|[0m                                     [34m│ [38;5;231m  [38;5;203mconfig.yaml[38;5;231m: [38;5;203m|[0m

[1;4;33mk8s/discoverer/deployment.yaml[0m

[36m──────────────[0m[36m┐[0m
[34m21[0m:[38;5;231m [38;5;203mmetadata[38;5;231m: [0m[36m│[0m
[36m──────────────[0m[36m┴[0m[36m───────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[38;5;231m  [38;5;203mlabels[38;5;231m:[0m                                            [34m│ [38;5;231m  [38;5;203mlabels[38;5;231m:[0m
[34m[38;5;231m    [38;5;203mapp[38;5;231m: [38;5;186mvald-discoverer[0m                             [34m│ [38;5;231m    [38;5;203mapp[38;5;231m: [38;5;186mvald-discoverer[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/name[38;5;231m: [38;5;186mvald[0m                     [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/name[38;5;231m: [38;5;186mvald[0m
[34m[48;5;52;38;5;231m    [38;5;203mhelm.sh/chart[38;5;231m: [38;5;186mvald-v1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                       [0m[34m│ [48;5;22;38;5;231m    [38;5;203mhelm.sh/chart[38;5;231m: [38;5;186mvald-v1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/managed-by[38;5;231m: [38;5;186mHelm[0m               [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/managed-by[38;5;231m: [38;5;186mHelm[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/instance[38;5;231m: [38;5;186mrelease-name[0m         [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/instance[38;5;231m: [38;5;186mrelease-name[0m
[34m[48;5;52;38;5;231m    [38;5;203mapp.kubernetes.io/version[38;5;231m: [38;5;186mv1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                [0m[34m│ [48;5;22;38;5;231m    [38;5;203mapp.kubernetes.io/version[38;5;231m: [38;5;186mv1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/component[38;5;231m: [38;5;186mdiscoverer[0m          [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/component[38;5;231m: [38;5;186mdiscoverer[0m
[34m[38;5;203mspec[38;5;231m:[0m                                                [34m│ [38;5;203mspec[38;5;231m:[0m
[34m[38;5;231m  [38;5;203mprogressDeadlineSeconds[38;5;231m: [38;5;141m600[0m                       [34m│ [38;5;231m  [38;5;203mprogressDeadlineSeconds[38;5;231m: [38;5;141m600[0m

[1;4;33mk8s/discoverer/pdb.yaml[0m

[36m──────────────[0m[36m┐[0m
[34m20[0m:[38;5;231m [38;5;203mmetadata[38;5;231m: [0m[36m│[0m
[36m──────────────[0m[36m┴[0m[36m───────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[38;5;231m  [38;5;203mname[38;5;231m: [38;5;186mvald-discoverer[0m                              [34m│ [38;5;231m  [38;5;203mname[38;5;231m: [38;5;186mvald-discoverer[0m
[34m[38;5;231m  [38;5;203mlabels[38;5;231m:[0m                                            [34m│ [38;5;231m  [38;5;203mlabels[38;5;231m:[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/name[38;5;231m: [38;5;186mvald[0m                     [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/name[38;5;231m: [38;5;186mvald[0m
[34m[48;5;52;38;5;231m    [38;5;203mhelm.sh/chart[38;5;231m: [38;5;186mvald-v1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                       [0m[34m│ [48;5;22;38;5;231m    [38;5;203mhelm.sh/chart[38;5;231m: [38;5;186mvald-v1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/managed-by[38;5;231m: [38;5;186mHelm[0m               [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/managed-by[38;5;231m: [38;5;186mHelm[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/instance[38;5;231m: [38;5;186mrelease-name[0m         [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/instance[38;5;231m: [38;5;186mrelease-name[0m
[34m[48;5;52;38;5;231m    [38;5;203mapp.kubernetes.io/version[38;5;231m: [38;5;186mv1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                [0m[34m│ [48;5;22;38;5;231m    [38;5;203mapp.kubernetes.io/version[38;5;231m: [38;5;186mv1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/component[38;5;231m: [38;5;186mdiscoverer[0m          [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/component[38;5;231m: [38;5;186mdiscoverer[0m
[34m[38;5;203mspec[38;5;231m:[0m                                                [34m│ [38;5;203mspec[38;5;231m:[0m
[34m[38;5;231m  [38;5;203mmaxUnavailable[38;5;231m: [38;5;186m50%[0m                                [34m│ [38;5;231m  [38;5;203mmaxUnavailable[38;5;231m: [38;5;186m50%[0m

[1;4;33mk8s/discoverer/priorityclass.yaml[0m

[36m──────────────[0m[36m┐[0m
[34m20[0m:[38;5;231m [38;5;203mmetadata[38;5;231m: [0m[36m│[0m
[36m──────────────[0m[36m┴[0m[36m───────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[38;5;231m  [38;5;203mname[38;5;231m: [38;5;186mdefault-vald-discoverer-priority[0m             [34m│ [38;5;231m  [38;5;203mname[38;5;231m: [38;5;186mdefault-vald-discoverer-priority[0m
[34m[38;5;231m  [38;5;203mlabels[38;5;231m:[0m                                            [34m│ [38;5;231m  [38;5;203mlabels[38;5;231m:[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/name[38;5;231m: [38;5;186mvald[0m                     [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/name[38;5;231m: [38;5;186mvald[0m
[34m[48;5;52;38;5;231m    [38;5;203mhelm.sh/chart[38;5;231m: [38;5;186mvald-v1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                       [0m[34m│ [48;5;22;38;5;231m    [38;5;203mhelm.sh/chart[38;5;231m: [38;5;186mvald-v1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/managed-by[38;5;231m: [38;5;186mHelm[0m               [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/managed-by[38;5;231m: [38;5;186mHelm[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/instance[38;5;231m: [38;5;186mrelease-name[0m         [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/instance[38;5;231m: [38;5;186mrelease-name[0m
[34m[48;5;52;38;5;231m    [38;5;203mapp.kubernetes.io/version[38;5;231m: [38;5;186mv1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                [0m[34m│ [48;5;22;38;5;231m    [38;5;203mapp.kubernetes.io/version[38;5;231m: [38;5;186mv1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/component[38;5;231m: [38;5;186mdiscoverer[0m          [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/component[38;5;231m: [38;5;186mdiscoverer[0m
[34m[38;5;203mvalue[38;5;231m: [38;5;186m1e+06[0m                                         [34m│ [38;5;203mvalue[38;5;231m: [38;5;186m1e+06[0m
[34m[38;5;203mglobalDefault[38;5;231m: [38;5;141mfalse[0m                                 [34m│ [38;5;203mglobalDefault[38;5;231m: [38;5;141mfalse[0m

[1;4;33mk8s/discoverer/serviceaccount.yaml[0m

[36m──────────────[0m[36m┐[0m
[34m20[0m:[38;5;231m [38;5;203mmetadata[38;5;231m: [0m[36m│[0m
[36m──────────────[0m[36m┴[0m[36m───────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[38;5;231m  [38;5;203mname[38;5;231m: [38;5;186mvald[0m                                         [34m│ [38;5;231m  [38;5;203mname[38;5;231m: [38;5;186mvald[0m
[34m[38;5;231m  [38;5;203mlabels[38;5;231m:[0m                                            [34m│ [38;5;231m  [38;5;203mlabels[38;5;231m:[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/name[38;5;231m: [38;5;186mvald[0m                     [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/name[38;5;231m: [38;5;186mvald[0m
[34m[48;5;52;38;5;231m    [38;5;203mhelm.sh/chart[38;5;231m: [38;5;186mvald-v1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                       [0m[34m│ [48;5;22;38;5;231m    [38;5;203mhelm.sh/chart[38;5;231m: [38;5;186mvald-v1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/managed-by[38;5;231m: [38;5;186mHelm[0m               [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/managed-by[38;5;231m: [38;5;186mHelm[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/instance[38;5;231m: [38;5;186mrelease-name[0m         [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/instance[38;5;231m: [38;5;186mrelease-name[0m
[34m[48;5;52;38;5;231m    [38;5;203mapp.kubernetes.io/version[38;5;231m: [38;5;186mv1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                [0m[34m│ [48;5;22;38;5;231m    [38;5;203mapp.kubernetes.io/version[38;5;231m: [38;5;186mv1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/component[38;5;231m: [38;5;186mdiscoverer[0m          [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/component[38;5;231m: [38;5;186mdiscoverer[0m

[1;4;33mk8s/discoverer/svc.yaml[0m

[36m──────────────[0m[36m┐[0m
[34m20[0m:[38;5;231m [38;5;203mmetadata[38;5;231m: [0m[36m│[0m
[36m──────────────[0m[36m┴[0m[36m───────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[38;5;231m  [38;5;203mname[38;5;231m: [38;5;186mvald-discoverer[0m                              [34m│ [38;5;231m  [38;5;203mname[38;5;231m: [38;5;186mvald-discoverer[0m
[34m[38;5;231m  [38;5;203mlabels[38;5;231m:[0m                                            [34m│ [38;5;231m  [38;5;203mlabels[38;5;231m:[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/name[38;5;231m: [38;5;186mvald[0m                     [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/name[38;5;231m: [38;5;186mvald[0m
[34m[48;5;52;38;5;231m    [38;5;203mhelm.sh/chart[38;5;231m: [38;5;186mvald-v1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                       [0m[34m│ [48;5;22;38;5;231m    [38;5;203mhelm.sh/chart[38;5;231m: [38;5;186mvald-v1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/managed-by[38;5;231m: [38;5;186mHelm[0m               [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/managed-by[38;5;231m: [38;5;186mHelm[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/instance[38;5;231m: [38;5;186mrelease-name[0m         [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/instance[38;5;231m: [38;5;186mrelease-name[0m
[34m[48;5;52;38;5;231m    [38;5;203mapp.kubernetes.io/version[38;5;231m: [38;5;186mv1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                [0m[34m│ [48;5;22;38;5;231m    [38;5;203mapp.kubernetes.io/version[38;5;231m: [38;5;186mv1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/component[38;5;231m: [38;5;186mdiscoverer[0m          [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/component[38;5;231m: [38;5;186mdiscoverer[0m
[34m[38;5;203mspec[38;5;231m:[0m                                                [34m│ [38;5;203mspec[38;5;231m:[0m
[34m[38;5;231m  [38;5;203mports[38;5;231m:[0m                                             [34m│ [38;5;231m  [38;5;203mports[38;5;231m:[0m

[1;4;33mk8s/gateway/lb/configmap.yaml[0m

[36m──────────────[0m[36m┐[0m
[34m20[0m:[38;5;231m [38;5;203mmetadata[38;5;231m: [0m[36m│[0m
[36m──────────────[0m[36m┴[0m[36m───────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[38;5;231m  [38;5;203mname[38;5;231m: [38;5;186mvald-lb-gateway-config[0m                       [34m│ [38;5;231m  [38;5;203mname[38;5;231m: [38;5;186mvald-lb-gateway-config[0m
[34m[38;5;231m  [38;5;203mlabels[38;5;231m:[0m                                            [34m│ [38;5;231m  [38;5;203mlabels[38;5;231m:[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/name[38;5;231m: [38;5;186mvald[0m                     [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/name[38;5;231m: [38;5;186mvald[0m
[34m[48;5;52;38;5;231m    [38;5;203mhelm.sh/chart[38;5;231m: [38;5;186mvald-v1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                       [0m[34m│ [48;5;22;38;5;231m    [38;5;203mhelm.sh/chart[38;5;231m: [38;5;186mvald-v1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/managed-by[38;5;231m: [38;5;186mHelm[0m               [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/managed-by[38;5;231m: [38;5;186mHelm[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/instance[38;5;231m: [38;5;186mrelease-name[0m         [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/instance[38;5;231m: [38;5;186mrelease-name[0m
[34m[48;5;52;38;5;231m    [38;5;203mapp.kubernetes.io/version[38;5;231m: [38;5;186mv1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                [0m[34m│ [48;5;22;38;5;231m    [38;5;203mapp.kubernetes.io/version[38;5;231m: [38;5;186mv1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/component[38;5;231m: [38;5;186mgateway-lb[0m          [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/component[38;5;231m: [38;5;186mgateway-lb[0m
[34m[38;5;203mdata[38;5;231m:[0m                                                [34m│ [38;5;203mdata[38;5;231m:[0m
[34m[38;5;231m  [38;5;203mconfig.yaml[38;5;231m: [38;5;203m|[0m                                     [34m│ [38;5;231m  [38;5;203mconfig.yaml[38;5;231m: [38;5;203m|[0m

[1;4;33mk8s/gateway/lb/deployment.yaml[0m

[36m──────────────[0m[36m┐[0m
[34m21[0m:[38;5;231m [38;5;203mmetadata[38;5;231m: [0m[36m│[0m
[36m──────────────[0m[36m┴[0m[36m───────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[38;5;231m  [38;5;203mlabels[38;5;231m:[0m                                            [34m│ [38;5;231m  [38;5;203mlabels[38;5;231m:[0m
[34m[38;5;231m    [38;5;203mapp[38;5;231m: [38;5;186mvald-lb-gateway[0m                             [34m│ [38;5;231m    [38;5;203mapp[38;5;231m: [38;5;186mvald-lb-gateway[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/name[38;5;231m: [38;5;186mvald[0m                     [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/name[38;5;231m: [38;5;186mvald[0m
[34m[48;5;52;38;5;231m    [38;5;203mhelm.sh/chart[38;5;231m: [38;5;186mvald-v1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                       [0m[34m│ [48;5;22;38;5;231m    [38;5;203mhelm.sh/chart[38;5;231m: [38;5;186mvald-v1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/managed-by[38;5;231m: [38;5;186mHelm[0m               [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/managed-by[38;5;231m: [38;5;186mHelm[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/instance[38;5;231m: [38;5;186mrelease-name[0m         [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/instance[38;5;231m: [38;5;186mrelease-name[0m
[34m[48;5;52;38;5;231m    [38;5;203mapp.kubernetes.io/version[38;5;231m: [38;5;186mv1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                [0m[34m│ [48;5;22;38;5;231m    [38;5;203mapp.kubernetes.io/version[38;5;231m: [38;5;186mv1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/component[38;5;231m: [38;5;186mgateway-lb[0m          [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/component[38;5;231m: [38;5;186mgateway-lb[0m
[34m[38;5;203mspec[38;5;231m:[0m                                                [34m│ [38;5;203mspec[38;5;231m:[0m
[34m[38;5;231m  [38;5;203mprogressDeadlineSeconds[38;5;231m: [38;5;141m600[0m                       [34m│ [38;5;231m  [38;5;203mprogressDeadlineSeconds[38;5;231m: [38;5;141m600[0m

[1;4;33mk8s/gateway/lb/hpa.yaml[0m

[36m──────────────[0m[36m┐[0m
[34m20[0m:[38;5;231m [38;5;203mmetadata[38;5;231m: [0m[36m│[0m
[36m──────────────[0m[36m┴[0m[36m───────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[38;5;231m  [38;5;203mname[38;5;231m: [38;5;186mvald-lb-gateway[0m                              [34m│ [38;5;231m  [38;5;203mname[38;5;231m: [38;5;186mvald-lb-gateway[0m
[34m[38;5;231m  [38;5;203mlabels[38;5;231m:[0m                                            [34m│ [38;5;231m  [38;5;203mlabels[38;5;231m:[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/name[38;5;231m: [38;5;186mvald[0m                     [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/name[38;5;231m: [38;5;186mvald[0m
[34m[48;5;52;38;5;231m    [38;5;203mhelm.sh/chart[38;5;231m: [38;5;186mvald-v1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                       [0m[34m│ [48;5;22;38;5;231m    [38;5;203mhelm.sh/chart[38;5;231m: [38;5;186mvald-v1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/managed-by[38;5;231m: [38;5;186mHelm[0m               [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/managed-by[38;5;231m: [38;5;186mHelm[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/instance[38;5;231m: [38;5;186mrelease-name[0m         [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/instance[38;5;231m: [38;5;186mrelease-name[0m
[34m[48;5;52;38;5;231m    [38;5;203mapp.kubernetes.io/version[38;5;231m: [38;5;186mv1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                [0m[34m│ [48;5;22;38;5;231m    [38;5;203mapp.kubernetes.io/version[38;5;231m: [38;5;186mv1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/component[38;5;231m: [38;5;186mgateway-lb[0m          [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/component[38;5;231m: [38;5;186mgateway-lb[0m
[34m[38;5;203mspec[38;5;231m:[0m                                                [34m│ [38;5;203mspec[38;5;231m:[0m
[34m[38;5;231m  [38;5;203mmaxReplicas[38;5;231m: [38;5;141m9[0m                                     [34m│ [38;5;231m  [38;5;203mmaxReplicas[38;5;231m: [38;5;141m9[0m

[1;4;33mk8s/gateway/lb/pdb.yaml[0m

[36m──────────────[0m[36m┐[0m
[34m20[0m:[38;5;231m [38;5;203mmetadata[38;5;231m: [0m[36m│[0m
[36m──────────────[0m[36m┴[0m[36m───────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[38;5;231m  [38;5;203mname[38;5;231m: [38;5;186mvald-lb-gateway[0m                              [34m│ [38;5;231m  [38;5;203mname[38;5;231m: [38;5;186mvald-lb-gateway[0m
[34m[38;5;231m  [38;5;203mlabels[38;5;231m:[0m                                            [34m│ [38;5;231m  [38;5;203mlabels[38;5;231m:[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/name[38;5;231m: [38;5;186mvald[0m                     [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/name[38;5;231m: [38;5;186mvald[0m
[34m[48;5;52;38;5;231m    [38;5;203mhelm.sh/chart[38;5;231m: [38;5;186mvald-v1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                       [0m[34m│ [48;5;22;38;5;231m    [38;5;203mhelm.sh/chart[38;5;231m: [38;5;186mvald-v1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/managed-by[38;5;231m: [38;5;186mHelm[0m               [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/managed-by[38;5;231m: [38;5;186mHelm[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/instance[38;5;231m: [38;5;186mrelease-name[0m         [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/instance[38;5;231m: [38;5;186mrelease-name[0m
[34m[48;5;52;38;5;231m    [38;5;203mapp.kubernetes.io/version[38;5;231m: [38;5;186mv1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                [0m[34m│ [48;5;22;38;5;231m    [38;5;203mapp.kubernetes.io/version[38;5;231m: [38;5;186mv1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/component[38;5;231m: [38;5;186mgateway-lb[0m          [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/component[38;5;231m: [38;5;186mgateway-lb[0m
[34m[38;5;203mspec[38;5;231m:[0m                                                [34m│ [38;5;203mspec[38;5;231m:[0m
[34m[38;5;231m  [38;5;203mmaxUnavailable[38;5;231m: [38;5;186m50%[0m                                [34m│ [38;5;231m  [38;5;203mmaxUnavailable[38;5;231m: [38;5;186m50%[0m

[1;4;33mk8s/gateway/lb/priorityclass.yaml[0m

[36m──────────────[0m[36m┐[0m
[34m20[0m:[38;5;231m [38;5;203mmetadata[38;5;231m: [0m[36m│[0m
[36m──────────────[0m[36m┴[0m[36m───────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[38;5;231m  [38;5;203mname[38;5;231m: [38;5;186mdefault-vald-lb-gateway-priority[0m             [34m│ [38;5;231m  [38;5;203mname[38;5;231m: [38;5;186mdefault-vald-lb-gateway-priority[0m
[34m[38;5;231m  [38;5;203mlabels[38;5;231m:[0m                                            [34m│ [38;5;231m  [38;5;203mlabels[38;5;231m:[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/name[38;5;231m: [38;5;186mvald[0m                     [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/name[38;5;231m: [38;5;186mvald[0m
[34m[48;5;52;38;5;231m    [38;5;203mhelm.sh/chart[38;5;231m: [38;5;186mvald-v1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                       [0m[34m│ [48;5;22;38;5;231m    [38;5;203mhelm.sh/chart[38;5;231m: [38;5;186mvald-v1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/managed-by[38;5;231m: [38;5;186mHelm[0m               [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/managed-by[38;5;231m: [38;5;186mHelm[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/instance[38;5;231m: [38;5;186mrelease-name[0m         [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/instance[38;5;231m: [38;5;186mrelease-name[0m
[34m[48;5;52;38;5;231m    [38;5;203mapp.kubernetes.io/version[38;5;231m: [38;5;186mv1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                [0m[34m│ [48;5;22;38;5;231m    [38;5;203mapp.kubernetes.io/version[38;5;231m: [38;5;186mv1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/component[38;5;231m: [38;5;186mgateway-lb[0m          [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/component[38;5;231m: [38;5;186mgateway-lb[0m
[34m[38;5;203mvalue[38;5;231m: [38;5;186m1e+06[0m                                         [34m│ [38;5;203mvalue[38;5;231m: [38;5;186m1e+06[0m
[34m[38;5;203mglobalDefault[38;5;231m: [38;5;141mfalse[0m                                 [34m│ [38;5;203mglobalDefault[38;5;231m: [38;5;141mfalse[0m

[1;4;33mk8s/gateway/lb/svc.yaml[0m

[36m──────────────[0m[36m┐[0m
[34m20[0m:[38;5;231m [38;5;203mmetadata[38;5;231m: [0m[36m│[0m
[36m──────────────[0m[36m┴[0m[36m───────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[38;5;231m  [38;5;203mname[38;5;231m: [38;5;186mvald-lb-gateway[0m                              [34m│ [38;5;231m  [38;5;203mname[38;5;231m: [38;5;186mvald-lb-gateway[0m
[34m[38;5;231m  [38;5;203mlabels[38;5;231m:[0m                                            [34m│ [38;5;231m  [38;5;203mlabels[38;5;231m:[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/name[38;5;231m: [38;5;186mvald[0m                     [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/name[38;5;231m: [38;5;186mvald[0m
[34m[48;5;52;38;5;231m    [38;5;203mhelm.sh/chart[38;5;231m: [38;5;186mvald-v1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                       [0m[34m│ [48;5;22;38;5;231m    [38;5;203mhelm.sh/chart[38;5;231m: [38;5;186mvald-v1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/managed-by[38;5;231m: [38;5;186mHelm[0m               [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/managed-by[38;5;231m: [38;5;186mHelm[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/instance[38;5;231m: [38;5;186mrelease-name[0m         [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/instance[38;5;231m: [38;5;186mrelease-name[0m
[34m[48;5;52;38;5;231m    [38;5;203mapp.kubernetes.io/version[38;5;231m: [38;5;186mv1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                [0m[34m│ [48;5;22;38;5;231m    [38;5;203mapp.kubernetes.io/version[38;5;231m: [38;5;186mv1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/component[38;5;231m: [38;5;186mgateway-lb[0m          [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/component[38;5;231m: [38;5;186mgateway-lb[0m
[34m[38;5;203mspec[38;5;231m:[0m                                                [34m│ [38;5;203mspec[38;5;231m:[0m
[34m[38;5;231m  [38;5;203mports[38;5;231m:[0m                                             [34m│ [38;5;231m  [38;5;203mports[38;5;231m:[0m

[1;4;33mk8s/manager/index/configmap.yaml[0m

[36m──────────────[0m[36m┐[0m
[34m20[0m:[38;5;231m [38;5;203mmetadata[38;5;231m: [0m[36m│[0m
[36m──────────────[0m[36m┴[0m[36m───────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[38;5;231m  [38;5;203mname[38;5;231m: [38;5;186mvald-manager-index-config[0m                    [34m│ [38;5;231m  [38;5;203mname[38;5;231m: [38;5;186mvald-manager-index-config[0m
[34m[38;5;231m  [38;5;203mlabels[38;5;231m:[0m                                            [34m│ [38;5;231m  [38;5;203mlabels[38;5;231m:[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/name[38;5;231m: [38;5;186mvald[0m                     [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/name[38;5;231m: [38;5;186mvald[0m
[34m[48;5;52;38;5;231m    [38;5;203mhelm.sh/chart[38;5;231m: [38;5;186mvald-v1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                       [0m[34m│ [48;5;22;38;5;231m    [38;5;203mhelm.sh/chart[38;5;231m: [38;5;186mvald-v1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/managed-by[38;5;231m: [38;5;186mHelm[0m               [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/managed-by[38;5;231m: [38;5;186mHelm[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/instance[38;5;231m: [38;5;186mrelease-name[0m         [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/instance[38;5;231m: [38;5;186mrelease-name[0m
[34m[48;5;52;38;5;231m    [38;5;203mapp.kubernetes.io/version[38;5;231m: [38;5;186mv1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                [0m[34m│ [48;5;22;38;5;231m    [38;5;203mapp.kubernetes.io/version[38;5;231m: [38;5;186mv1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/component[38;5;231m: [38;5;186mmanager-index[0m       [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/component[38;5;231m: [38;5;186mmanager-index[0m
[34m[38;5;203mdata[38;5;231m:[0m                                                [34m│ [38;5;203mdata[38;5;231m:[0m
[34m[38;5;231m  [38;5;203mconfig.yaml[38;5;231m: [38;5;203m|[0m                                     [34m│ [38;5;231m  [38;5;203mconfig.yaml[38;5;231m: [38;5;203m|[0m

[1;4;33mk8s/manager/index/deployment.yaml[0m

[36m──────────────[0m[36m┐[0m
[34m21[0m:[38;5;231m [38;5;203mmetadata[38;5;231m: [0m[36m│[0m
[36m──────────────[0m[36m┴[0m[36m───────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[38;5;231m  [38;5;203mlabels[38;5;231m:[0m                                            [34m│ [38;5;231m  [38;5;203mlabels[38;5;231m:[0m
[34m[38;5;231m    [38;5;203mapp[38;5;231m: [38;5;186mvald-manager-index[0m                          [34m│ [38;5;231m    [38;5;203mapp[38;5;231m: [38;5;186mvald-manager-index[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/name[38;5;231m: [38;5;186mvald[0m                     [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/name[38;5;231m: [38;5;186mvald[0m
[34m[48;5;52;38;5;231m    [38;5;203mhelm.sh/chart[38;5;231m: [38;5;186mvald-v1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                       [0m[34m│ [48;5;22;38;5;231m    [38;5;203mhelm.sh/chart[38;5;231m: [38;5;186mvald-v1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/managed-by[38;5;231m: [38;5;186mHelm[0m               [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/managed-by[38;5;231m: [38;5;186mHelm[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/instance[38;5;231m: [38;5;186mrelease-name[0m         [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/instance[38;5;231m: [38;5;186mrelease-name[0m
[34m[48;5;52;38;5;231m    [38;5;203mapp.kubernetes.io/version[38;5;231m: [38;5;186mv1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                [0m[34m│ [48;5;22;38;5;231m    [38;5;203mapp.kubernetes.io/version[38;5;231m: [38;5;186mv1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/component[38;5;231m: [38;5;186mmanager-index[0m       [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/component[38;5;231m: [38;5;186mmanager-index[0m
[34m[38;5;203mspec[38;5;231m:[0m                                                [34m│ [38;5;203mspec[38;5;231m:[0m
[34m[38;5;231m  [38;5;203mprogressDeadlineSeconds[38;5;231m: [38;5;141m600[0m                       [34m│ [38;5;231m  [38;5;203mprogressDeadlineSeconds[38;5;231m: [38;5;141m600[0m

[1;4;33mk8s/manager/index/pdb.yaml[0m

[36m──────────────[0m[36m┐[0m
[34m20[0m:[38;5;231m [38;5;203mmetadata[38;5;231m: [0m[36m│[0m
[36m──────────────[0m[36m┴[0m[36m───────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[38;5;231m  [38;5;203mname[38;5;231m: [38;5;186mvald-manager-index[0m                           [34m│ [38;5;231m  [38;5;203mname[38;5;231m: [38;5;186mvald-manager-index[0m
[34m[38;5;231m  [38;5;203mlabels[38;5;231m:[0m                                            [34m│ [38;5;231m  [38;5;203mlabels[38;5;231m:[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/name[38;5;231m: [38;5;186mvald[0m                     [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/name[38;5;231m: [38;5;186mvald[0m
[34m[48;5;52;38;5;231m    [38;5;203mhelm.sh/chart[38;5;231m: [38;5;186mvald-v1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                       [0m[34m│ [48;5;22;38;5;231m    [38;5;203mhelm.sh/chart[38;5;231m: [38;5;186mvald-v1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/managed-by[38;5;231m: [38;5;186mHelm[0m               [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/managed-by[38;5;231m: [38;5;186mHelm[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/instance[38;5;231m: [38;5;186mrelease-name[0m         [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/instance[38;5;231m: [38;5;186mrelease-name[0m
[34m[48;5;52;38;5;231m    [38;5;203mapp.kubernetes.io/version[38;5;231m: [38;5;186mv1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                [0m[34m│ [48;5;22;38;5;231m    [38;5;203mapp.kubernetes.io/version[38;5;231m: [38;5;186mv1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/component[38;5;231m: [38;5;186mmanager-index[0m       [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/component[38;5;231m: [38;5;186mmanager-index[0m
[34m[38;5;203mspec[38;5;231m:[0m                                                [34m│ [38;5;203mspec[38;5;231m:[0m
[34m[38;5;231m  [38;5;203mmaxUnavailable[38;5;231m: [38;5;186m50%[0m                                [34m│ [38;5;231m  [38;5;203mmaxUnavailable[38;5;231m: [38;5;186m50%[0m

[1;4;33mk8s/manager/index/priorityclass.yaml[0m

[36m──────────────[0m[36m┐[0m
[34m20[0m:[38;5;231m [38;5;203mmetadata[38;5;231m: [0m[36m│[0m
[36m──────────────[0m[36m┴[0m[36m───────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[38;5;231m  [38;5;203mname[38;5;231m: [38;5;186mdefault-vald-manager-index-priority[0m          [34m│ [38;5;231m  [38;5;203mname[38;5;231m: [38;5;186mdefault-vald-manager-index-priority[0m
[34m[38;5;231m  [38;5;203mlabels[38;5;231m:[0m                                            [34m│ [38;5;231m  [38;5;203mlabels[38;5;231m:[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/name[38;5;231m: [38;5;186mvald[0m                     [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/name[38;5;231m: [38;5;186mvald[0m
[34m[48;5;52;38;5;231m    [38;5;203mhelm.sh/chart[38;5;231m: [38;5;186mvald-v1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                       [0m[34m│ [48;5;22;38;5;231m    [38;5;203mhelm.sh/chart[38;5;231m: [38;5;186mvald-v1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/managed-by[38;5;231m: [38;5;186mHelm[0m               [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/managed-by[38;5;231m: [38;5;186mHelm[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/instance[38;5;231m: [38;5;186mrelease-name[0m         [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/instance[38;5;231m: [38;5;186mrelease-name[0m
[34m[48;5;52;38;5;231m    [38;5;203mapp.kubernetes.io/version[38;5;231m: [38;5;186mv1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                [0m[34m│ [48;5;22;38;5;231m    [38;5;203mapp.kubernetes.io/version[38;5;231m: [38;5;186mv1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/component[38;5;231m: [38;5;186mmanager-index[0m       [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/component[38;5;231m: [38;5;186mmanager-index[0m
[34m[38;5;203mvalue[38;5;231m: [38;5;186m1e+06[0m                                         [34m│ [38;5;203mvalue[38;5;231m: [38;5;186m1e+06[0m
[34m[38;5;203mglobalDefault[38;5;231m: [38;5;141mfalse[0m                                 [34m│ [38;5;203mglobalDefault[38;5;231m: [38;5;141mfalse[0m

[1;4;33mk8s/manager/index/svc.yaml[0m

[36m──────────────[0m[36m┐[0m
[34m20[0m:[38;5;231m [38;5;203mmetadata[38;5;231m: [0m[36m│[0m
[36m──────────────[0m[36m┴[0m[36m───────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[38;5;231m  [38;5;203mname[38;5;231m: [38;5;186mvald-manager-index[0m                           [34m│ [38;5;231m  [38;5;203mname[38;5;231m: [38;5;186mvald-manager-index[0m
[34m[38;5;231m  [38;5;203mlabels[38;5;231m:[0m                                            [34m│ [38;5;231m  [38;5;203mlabels[38;5;231m:[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/name[38;5;231m: [38;5;186mvald[0m                     [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/name[38;5;231m: [38;5;186mvald[0m
[34m[48;5;52;38;5;231m    [38;5;203mhelm.sh/chart[38;5;231m: [38;5;186mvald-v1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                       [0m[34m│ [48;5;22;38;5;231m    [38;5;203mhelm.sh/chart[38;5;231m: [38;5;186mvald-v1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/managed-by[38;5;231m: [38;5;186mHelm[0m               [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/managed-by[38;5;231m: [38;5;186mHelm[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/instance[38;5;231m: [38;5;186mrelease-name[0m         [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/instance[38;5;231m: [38;5;186mrelease-name[0m
[34m[48;5;52;38;5;231m    [38;5;203mapp.kubernetes.io/version[38;5;231m: [38;5;186mv1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                [0m[34m│ [48;5;22;38;5;231m    [38;5;203mapp.kubernetes.io/version[38;5;231m: [38;5;186mv1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/component[38;5;231m: [38;5;186mmanager-index[0m       [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/component[38;5;231m: [38;5;186mmanager-index[0m
[34m[38;5;203mspec[38;5;231m:[0m                                                [34m│ [38;5;203mspec[38;5;231m:[0m
[34m[38;5;231m  [38;5;203mports[38;5;231m:[0m                                             [34m│ [38;5;231m  [38;5;203mports[38;5;231m:[0m

[1;4;33mk8s/operator/helm/operator.yaml[0m

[36m──────────────[0m[36m┐[0m
[34m22[0m:[38;5;231m [38;5;203mmetadata[38;5;231m: [0m[36m│[0m
[36m──────────────[0m[36m┴[0m[36m───────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[38;5;231m  [38;5;203mlabels[38;5;231m:[0m                                            [34m│ [38;5;231m  [38;5;203mlabels[38;5;231m:[0m
[34m[38;5;231m    [38;5;203mapp[38;5;231m: [38;5;186mvald-helm-operator[0m                          [34m│ [38;5;231m    [38;5;203mapp[38;5;231m: [38;5;186mvald-helm-operator[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/name[38;5;231m: [38;5;186mvald-helm-operator[0m       [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/name[38;5;231m: [38;5;186mvald-helm-operator[0m
[34m[48;5;52;38;5;231m    [38;5;203mhelm.sh/chart[38;5;231m: [38;5;186mvald-helm-operator-v1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m         [0m[34m│ [48;5;22;38;5;231m    [38;5;203mhelm.sh/chart[38;5;231m: [38;5;186mvald-helm-operator-v1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/managed-by[38;5;231m: [38;5;186mHelm[0m               [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/managed-by[38;5;231m: [38;5;186mHelm[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/instance[38;5;231m: [38;5;186mrelease-name[0m         [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/instance[38;5;231m: [38;5;186mrelease-name[0m
[34m[48;5;52;38;5;231m    [38;5;203mapp.kubernetes.io/version[38;5;231m: [38;5;186mv1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                [0m[34m│ [48;5;22;38;5;231m    [38;5;203mapp.kubernetes.io/version[38;5;231m: [38;5;186mv1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/component[38;5;231m: [38;5;186mhelm-operator[0m       [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/component[38;5;231m: [38;5;186mhelm-operator[0m
[34m[38;5;203mspec[38;5;231m:[0m                                                [34m│ [38;5;203mspec[38;5;231m:[0m
[34m[38;5;231m  [38;5;203mreplicas[38;5;231m: [38;5;141m2[0m                                        [34m│ [38;5;231m  [38;5;203mreplicas[38;5;231m: [38;5;141m2[0m

[36m──────────[0m[36m┐[0m
[34m43[0m:[38;5;231m [38;5;203mspec[38;5;231m: [0m[36m│[0m
[36m──────────[0m[36m┴[0m[36m───────────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[38;5;231m      [38;5;203mserviceAccountName[38;5;231m: [38;5;186mvald-helm-operator[0m         [34m│ [38;5;231m      [38;5;203mserviceAccountName[38;5;231m: [38;5;186mvald-helm-operator[0m
[34m[38;5;231m      [38;5;203mcontainers[38;5;231m:[0m                                    [34m│ [38;5;231m      [38;5;203mcontainers[38;5;231m:[0m
[34m[38;5;231m        - [38;5;203mname[38;5;231m: [38;5;186mvald-helm-operator[0m                   [34m│ [38;5;231m        - [38;5;203mname[38;5;231m: [38;5;186mvald-helm-operator[0m
[34m[48;5;52;38;5;231m          [38;5;203mimage[38;5;231m: [38;5;186m"vdaas/vald-helm-operator:v1.[48;5;124m7[48;5;52m.[48;5;124m0[48;5;52m"[0m[48;5;52m   [0m[34m│ [48;5;22;38;5;231m          [38;5;203mimage[38;5;231m: [38;5;186m"vdaas/vald-helm-operator:v1.[48;5;28m6[48;5;22m.[48;5;28m3[48;5;22m"[0m[48;5;22m[0K[0m
[34m[38;5;231m          [38;5;203mimagePullPolicy[38;5;231m: [38;5;186mAlways[0m                    [34m│ [38;5;231m          [38;5;203mimagePullPolicy[38;5;231m: [38;5;186mAlways[0m
[34m[38;5;231m          [38;5;203margs[38;5;231m:[0m                                      [34m│ [38;5;231m          [38;5;203margs[38;5;231m:[0m
[34m[38;5;231m            - [38;5;186m"run"[0m                                  [34m│ [38;5;231m            - [38;5;186m"run"[0m

[1;4;33mk8s/operator/helm/svc.yaml[0m

[36m──────────────[0m[36m┐[0m
[34m20[0m:[38;5;231m [38;5;203mmetadata[38;5;231m: [0m[36m│[0m
[36m──────────────[0m[36m┴[0m[36m───────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[38;5;231m  [38;5;203mname[38;5;231m: [38;5;186mvald-helm-operator[0m                           [34m│ [38;5;231m  [38;5;203mname[38;5;231m: [38;5;186mvald-helm-operator[0m
[34m[38;5;231m  [38;5;203mlabels[38;5;231m:[0m                                            [34m│ [38;5;231m  [38;5;203mlabels[38;5;231m:[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/name[38;5;231m: [38;5;186mvald-helm-operator[0m       [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/name[38;5;231m: [38;5;186mvald-helm-operator[0m
[34m[48;5;52;38;5;231m    [38;5;203mhelm.sh/chart[38;5;231m: [38;5;186mvald-helm-operator-v1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m         [0m[34m│ [48;5;22;38;5;231m    [38;5;203mhelm.sh/chart[38;5;231m: [38;5;186mvald-helm-operator-v1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/managed-by[38;5;231m: [38;5;186mHelm[0m               [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/managed-by[38;5;231m: [38;5;186mHelm[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/instance[38;5;231m: [38;5;186mrelease-name[0m         [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/instance[38;5;231m: [38;5;186mrelease-name[0m
[34m[48;5;52;38;5;231m    [38;5;203mapp.kubernetes.io/version[38;5;231m: [38;5;186mv1.[48;5;124m7[48;5;52m.[48;5;124m0[0m[48;5;52m                [0m[34m│ [48;5;22;38;5;231m    [38;5;203mapp.kubernetes.io/version[38;5;231m: [38;5;186mv1.[48;5;28m6[48;5;22m.[48;5;28m3[0m[48;5;22m[0K[0m
[34m[38;5;231m    [38;5;203mapp.kubernetes.io/component[38;5;231m: [38;5;186mhelm-operator[0m       [34m│ [38;5;231m    [38;5;203mapp.kubernetes.io/component[38;5;231m: [38;5;186mhelm-operator[0m
[34m[38;5;203mspec[38;5;231m:[0m                                                [34m│ [38;5;203mspec[38;5;231m:[0m
[34m[38;5;231m  [38;5;203mports[38;5;231m:[0m                                             [34m│ [38;5;231m  [38;5;203mports[38;5;231m:[0m

[1;4;33mversions/HELM_VERSION[0m

[36m───[0m[36m┐[0m
[34m1[0m: [36m│[0m
[36m───[0m[36m┴[0m[36m──────────────────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[48;5;52;38;5;231mv3.[48;5;124m10[48;5;52m.[48;5;124m3[0m[48;5;52m                                              [0m[34m│ [48;5;22;38;5;231mv3.[48;5;28m11[48;5;22m.[48;5;28m0[0m[48;5;22m[0K[0m

[1;4;33mversions/JAEGER_OPERATOR_VERSION[0m

[36m───[0m[36m┐[0m
[34m1[0m: [36m│[0m
[36m───[0m[36m┴[0m[36m──────────────────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[48;5;52;38;5;231m2.[48;5;124m38[48;5;52m.0[0m[48;5;52m                                               [0m[34m│ [48;5;22;38;5;231m2.[48;5;28m39[48;5;22m.0[0m[48;5;22m[0K[0m

[1;4;33mversions/KUBELINTER_VERSION[0m

[36m───[0m[36m┐[0m
[34m1[0m: [36m│[0m
[36m───[0m[36m┴[0m[36m──────────────────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[48;5;52;38;5;231m0.[48;5;124m5[48;5;52m.[48;5;124m1[0m[48;5;52m                                                [0m[34m│ [48;5;22;38;5;231m0.[48;5;28m6[48;5;22m.[48;5;28m0[0m[48;5;22m[0K[0m

[1;4;33mversions/TELEPRESENCE_VERSION[0m

[36m───[0m[36m┐[0m
[34m1[0m: [36m│[0m
[36m───[0m[36m┴[0m[36m──────────────────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[48;5;52;38;5;231m2.10.[48;5;124m3[48;5;52m-rc.[48;5;124m1[0m[48;5;52m                                          [0m[34m│ [48;5;22;38;5;231m2.10.[48;5;28m4[48;5;22m-rc.[48;5;28m0[0m[48;5;22m[0K[0m

[1;4;33mversions/VALDCLI_VERSION[0m

[36m───[0m[36m┐[0m
[34m1[0m: [36m│[0m
[36m───[0m[36m┴[0m[36m──────────────────────────────────────────────────────────────────────────────────────────────────────[0m
[34m[48;5;52;38;5;231mv1.[48;5;124m6[48;5;52m.[48;5;124m3[0m[48;5;52m                                               [0m[34m│ [48;5;22;38;5;231mv1.[48;5;28m7[48;5;22m.[48;5;28m0[0m[48;5;22m[0K[0m
