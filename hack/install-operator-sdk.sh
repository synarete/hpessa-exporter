#!/bin/bash -e
# see: https://sdk.operatorframework.io/docs/installation/
tool="$1"
tooldir="$(dirname "${tool}")"
dl_url="https://github.com/operator-framework/operator-sdk/releases/download/v1.15.0"
arch=$(case $(uname -m) in x86_64) echo -n amd64 ;; aarch64) echo -n arm64 ;; *) echo -n "$(uname -m)" ;; esac)
osname=$(uname | awk '{print tolower($0)}')

command -v curl > /dev/null
echo "installing operator-stk at ${tool} (${osname}-${arch})"
[ -d "${tooldir}" ] || exit 1
cd "${tooldir}"
curl -LO "${dl_url}/operator-sdk_${osname}_${arch}"
chmod +x "operator-sdk_${osname}_${arch}"
mv "operator-sdk_${osname}_${arch}" "${tool}"
"${tool}" version
