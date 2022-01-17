#!/bin/bash -e
tool="$1"
tooldir="$(dirname "${tool}")"
url="https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh"

command -v curl > /dev/null
echo "installing golangci-lint at ${tool}"
[ -d "${tooldir}" ] || exit 1
curl -JL "${url}" | sh -s -- -b "${tooldir}" v1.43.0
mv "${tooldir}"/golangci-lint "${tooldir}"/golangci-lint.tmp
chmod +x "${tooldir}"/golangci-lint.tmp
mv "${tooldir}"/golangci-lint.tmp "${tool}"
"${tool}" version
