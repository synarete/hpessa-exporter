#!/bin/bash -e
tool="$1"
tooldir="$(dirname "${tool}")"
tmpdir="$2"

command -v go > /dev/null
echo "installing yq at ${tool}"
[ -d "${tooldir}" ] || exit 1
[ -d "${tmpdir}" ] || exit 2

cd "$tmpdir"
go mod init tmp
GO111MODULE=on GOBIN="${tooldir}" go get \
  github.com/mikefarah/yq/v4

