#!/bin/bash
IFS=$'\n\t'
self="$(readlink -f "$0")"
base="$(readlink -f "$(dirname "${self}")/../")"

cd "${base}" || exit 1
gittop=$(git rev-parse --show-toplevel > /dev/null 2>&1 && echo -n "git-repo")

gitrevision=""
gitdirty=""
if [ -n "${gittop}" ]; then
  gitrevision=$(git describe --abbrev=7 --always --dirty=+)
elif [ -f "${base}/GITREV" ]; then
  gitrevision="$(head -1 "${base}/GITREV")"
fi
echo -n "$gitrevision" | tr -d ' \t\v\n' ;

