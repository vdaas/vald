#!/bin/bash -eu

#
# Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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
:

#
# This script is executed as postAttachCommand in devcontainer.json
# This script does...
# - create symbolic link of config.yaml for easier development
# - add command history setting to .zshrc to persist history
#

echo "creating symbolic link of config..."

LINK_TARGET="$(pwd)/cmd/agent/core/ngt/sample.yaml"
LINK_SRC="/etc/server/config.yaml"

mkdir -p /etc/server

if [ ! -e "$LINK_SRC" ]; then
    ln -s "$LINK_TARGET" "$LINK_SRC"
    echo "created symbolic link: $LINK_SRC -> $LINK_TARGET"
else
    echo "$LINK_SRC already exists"
fi


echo "adding history setting to .zshrc..."

LINE1="export HISTFILE=/commandhistory/.zsh_history"
LINE2="setopt INC_APPEND_HISTORY"

ZSHRC="/root/.zshrc"

# write only if those lines don't exist
grep -qxF "$LINE1" "$ZSHRC" || echo "$LINE1" >> "$ZSHRC"
grep -qxF "$LINE2" "$ZSHRC" || echo "$LINE2" >> "$ZSHRC"

echo "added history setting to .zshrc"
