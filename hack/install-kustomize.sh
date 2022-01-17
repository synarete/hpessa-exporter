#!/bin/bash -e
tool="$1"
tooldir="$(dirname "${tool}")"
tmpdir="$2"

command -v go > /dev/null
echo "installing kustomize at ${tool}"
[ -d "${tooldir}" ] || exit 1
[ -d "${tmpdir}" ] || exit 2

cd "$tmpdir"
go mod init tmp
GOBIN="${tooldir}" go get sigs.k8s.io/kustomize/kustomize/v3@v3.5.4
"${tool}" version
