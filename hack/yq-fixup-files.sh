#!/bin/bash -e
YQ=${YQ:-$(command -v yq)}

for yaml in "$@"; do
  if [ -f "${yaml}" ]; then
    # use yq as yaml formatter
    ${YQ} eval "${yaml}" > "${yaml}".tmp

    # require each yaml to start with '---'
    hdr=$(head -c 3 "${yaml}".tmp)
    if [ "${hdr}" != "---" ]; then
      echo "---" > "${yaml}".tmp2
      cat "${yaml}".tmp >> "${yaml}".tmp2
      mv "${yaml}".tmp2 "${yaml}".tmp
    fi

    # override original yaml file
    mv "${yaml}".tmp "${yaml}"
  fi
done
