FROM golang:1.3
MAINTAINER Steven Jack <stevenmajack@gmail.com>

RUN apt-get update
RUN apt-get install git -y

RUN go get github.com/mitchellh/gox
RUN gox -build-toolchain

RUN go get github.com/codegangsta/cli
RUN go get gopkg.in/yaml.v2
RUN go get github.com/fatih/color
RUN go get github.com/mitchellh/go-homedir

WORKDIR /go
