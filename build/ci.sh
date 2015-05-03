#!/bin/bash

go test

if [[ $TRAVIS_TAG =~ ^v[0-9\.]+ ]]; then
  go get github.com/mitchellh/gox
  gox -os="darwin linux windows" -arch="amd64" -build-toolchain
  gox -os="darwin linux windows" -arch="amd64" -output="{{.Dir}}_{{.OS}}_x86_64"
fi
