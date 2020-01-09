#!/bin/bash

set -x

BINARY="terraform-provider-alienvault_${TRAVIS_TAG}"	BINARY="terraform-provider-alienvault_${TRAVIS_TAG}"
GO111MODULE=on	GO111MODULE=on


GOOS=darwin GOARCH=amd64 go build -o "${BINARY}"	package () {
zip "${BINARY}_darwin_amd64.zip" "${BINARY}"	  GOOS=$1
rm -f "${BINARY}"	  echo $GOOS

  dir="${BINARY}_${GOOS}_amd64"
GOOS=linux GOARCH=amd64 go build -o "${BINARY}"	  mkdir "${dir}"
zip "${BINARY}_linux_amd64.zip" "${BINARY}"	  GOOS=$1 GOARCH=amd64 go build -o "${dir}/${BINARY}"
rm -f "${BINARY}"	  zip "${BINARY}_${GOOS}_amd64.zip" -r "${dir}"
  rm -rf "./${dir:?}/"
}


GOOS=windows GOARCH=amd64 go build -o "${BINARY}"	package darwin
zip "${BINARY}_windows_amd64.zip" "${BINARY}"	package linux
rm -f "${BINARY}"
