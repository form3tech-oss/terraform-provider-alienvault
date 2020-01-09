#!/bin/bash

set -x

BINARY="terraform-provider-alienvault_${TRAVIS_TAG}"
GO111MODULE=on

package () {
  GOOS=$1
  echo $GOOS
  dir="${BINARY}_${GOOS}_amd64"
  mkdir "${dir}"
  GOOS=$1 GOARCH=amd64 go build -o "${dir}/${BINARY}"
  zip "${BINARY}_${GOOS}_amd64.zip" -r "${dir}"
  rm -rf "./${dir:?}/"
}

package darwin
package linux
package windows
