#!/bin/bash -eu
#
# This script is executed as postAttachCommand in devcontainer.json
#

echo "creating symbolic link of config ZSHRC..."

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
