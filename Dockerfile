#### Metadata
FROM ubuntu:16.04
MAINTAINER Antoine Pourchet <antoine.pourchet@gmail.com>

#### Image Building
USER root
ENV HOME /root
ENV GOPATH /go

# apt-get
RUN apt-get update
RUN apt-get install sudo
RUN sudo apt-get install -y man git curl python
RUN sudo apt-get install -y golang
RUN sudo apt-get install -y make
RUN sudo apt-get install -y jq
RUN sudo apt-get install -y screen

# Specifics
VOLUME /pingme
ADD . /go/src/github.com/apourchet/pingme
RUN go install github.com/apourchet/pingme/cmd/pingmesrv

EXPOSE 1025

ENTRYPOINT /go/bin/pingmesrv
