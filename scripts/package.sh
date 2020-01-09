#!/bin/bash
BINARY="terraform-provider-alienvault_${TRAVIS_TAG}"
GO111MODULE=on

GOOS=darwin GOARCH=amd64 go build -o "${BINARY}"
zip "${BINARY}_darwin_amd64.zip" "${BINARY}"
rm -f "${BINARY}"

GOOS=linux GOARCH=amd64 go build -o "${BINARY}"
zip "${BINARY}_linux_amd64.zip" "${BINARY}"
rm -f "${BINARY}"

GOOS=windows GOARCH=amd64 go build -o "${BINARY}"
zip "${BINARY}_windows_amd64.zip" "${BINARY}"
rm -f "${BINARY}"
