#
# Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
#
# Licensed under the Apache License, Version 2.0 (the "License");
# You may not use this file except in compliance with the License.
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

import os
import pathlib
import subprocess
import tempfile

repo = pathlib.Path(__file__).resolve().parents[2]
with tempfile.TemporaryDirectory() as td:
    root = pathlib.Path(td)
    (root/'versions/actions').mkdir(parents=True)
    (root/'bin').mkdir()
    git = root/'bin/git'
    git.write_text('''#!/bin/sh
case "$MODE" in
  network) exit 128;;
  missing) exit 0;;
  malformed) printf 'bad refs/tags/v1.2.3\\n';;
  branch) printf '%s %s\\n' "$(printf 'c%.0s' $(seq 1 40))" "refs/heads/$REF";;
  *) printf '%s %s\\n' "$(printf 'a%.0s' $(seq 1 40))" "refs/tags/$REF";
     if [ "$MODE" = annotated ]; then printf '%s %s\\n' "$(printf 'b%.0s' $(seq 1 40))" "refs/tags/$REF^{}"; fi;
     ;;
esac
''')
    git.chmod(0o755)
    makefile = root/'Makefile'
    makefile.write_text(f'''SHELL := /bin/bash
.SHELLFLAGS := -euo pipefail -c
SHELL_STRICT_EXEC := /bin/bash -euo pipefail -c
XARGS := xargs -I {{}} -P 4
XARGS_NO_RUN_IF_EMPTY := -r
CORES := 4
ROOTDIR := {root}
include {repo}/Makefile.d/functions.mk
pin:
\t@eval "$$GITHUB_ACTION_PIN_SCRIPT" _
update:
\t$(call update-github-actions,actions/checkout github/codeql-action/init)
''')
    env = dict(os.environ, PATH=str(root/'bin')+':'+os.environ['PATH'])
    for action, version, ref, mode, expected in [
        ('actions/checkout', '1.2.3', 'v1.2.3', 'lightweight', 'a'),
        ('github/codeql-action/init', '1.2.3', 'v1.2.3', 'annotated', 'b'),
        ('cirrus-actions/rebase', '1.8', '1.8', 'lightweight', 'a'),
        ('aquasecurity/trivy-action', 'master', 'master', 'branch', 'c'),
    ]:
        name = action.replace('/', '_').replace('-', '_').upper()
        (root/'versions/actions'/name).write_text(version+'\n')
        # A dedicated recipe sets $1, just as the production xargs invocation does.
        source = makefile.read_text().replace('@eval "$$GITHUB_ACTION_PIN_SCRIPT" _', f'@$(SHELL_STRICT_EXEC) \'eval "$$GITHUB_ACTION_PIN_SCRIPT"\' _ {action}')
        (root/'case.mk').write_text(source)
        result = subprocess.run(['make','-s','-f',str(root/'case.mk'),'pin'], env=dict(env, MODE=mode, REF=ref), text=True, capture_output=True)
        assert result.returncode == 0, result.stderr
        assert f'{action}@{expected*40} # {ref}' in result.stdout, result.stdout
        print('PASS', action, mode)
        for failure in ['network','missing','malformed']:
            result = subprocess.run(['make','-s','-f',str(root/'case.mk'),'pin'], env=dict(env, MODE=failure, REF=ref), text=True, capture_output=True)
            assert result.returncode != 0 and not result.stdout, (failure,result)
    print('PASS network, missing-ref and malformed-SHA failures')
    workflow = root/'.github/workflows/test.yaml'
    workflow.parent.mkdir(parents=True)
    original = 'uses: actions/checkout@' + 'd'*40 + ' # v1.0.0\n'
    workflow.write_text(original)
    (root/'.gitfiles').write_text('.github/workflows/test.yaml\n')
    with makefile.open('a') as f:
        f.write('override UPDATE_GITHUB_ACTION_SCRIPT = :\n')
    result = subprocess.run(['make', '-s', '-f', str(makefile), 'update'],
                            env=dict(env, MODE='missing', REF='v1.2.3'),
                            text=True, capture_output=True)
    assert result.returncode != 0, result.stdout
    assert workflow.read_text() == original
    print('PASS failed resolution preserves existing workflow pins')
    workflow.write_text(original + 'uses: github/codeql-action/init@old\n')
    result = subprocess.run(['make', '-s', '-f', str(makefile), 'update'],
                            env=dict(env, MODE='lightweight', REF='v1.2.3'),
                            text=True, capture_output=True)
    assert result.returncode == 0, result.stderr
    updated = workflow.read_text()
    assert updated == ('uses: actions/checkout@' + 'a'*40 + ' # v1.2.3\n'
                       + 'uses: github/codeql-action/init@' + 'a'*40 + ' # v1.2.3\n')
    result = subprocess.run(['make', '-s', '-f', str(makefile), 'update'],
                            env=dict(env, MODE='lightweight', REF='v1.2.3'),
                            text=True, capture_output=True)
    assert result.returncode == 0, result.stderr
    assert workflow.read_text() == updated
    print('PASS multiple actions are pinned and repeated updates are stable')
