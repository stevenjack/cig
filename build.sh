#!/bin/bash

if [[ $TRAVIS_BRANCH == 'master' ]]; then
  go get github.com/mitchellh/gox
  gox -os="darwin linux windows" -arch="amd64" -build-toolchain
  gox -os="darwin linux windows" -arch="amd64"
else
  go test
fi
