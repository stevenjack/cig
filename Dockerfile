FROM golang:1.4.2
MAINTAINER Steven Jack <stevenmajack@gmail.com>

RUN apt-get update && apt-get install git -yq

RUN go get github.com/mitchellh/gox
RUN gox -build-toolchain

WORKDIR /go
