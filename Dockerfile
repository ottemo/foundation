############################################################
# Dockerfile to build Ottemo Foundation
# Based on Ubuntu
############################################################

# Use the latest official ubuntu base image from docker hub
FROM ubuntu:latest

RUN apt-get update
RUN apt-get -y install bzr git golang-go gcc
WORKDIR /opt/go
RUN mkdir -pv /opt/go/src/github.com/ottemo
RUN mkdir -pv /opt/go/bin
RUN mkdir -pv /opt/go/pkg
RUN mkdir -pv /opt/media
RUN git clone https://ottemo-dev:freshbox111222333@github.com/ottemo/foundation.git /opt/go/src/github.com/ottemo/foundation
# sqlite setting is just for dev, this will be removed in the future.
RUN echo "db.sqlite3.uri=ottemo.db" >> /opt/go/src/github.com/ottemo/foundation/ottemo.ini
RUN echo "media.fsmedia.folder=/opt/media" >> /opt/go/src/github.com/ottemo/foundation/ottemo.ini
ENV GOPATH /opt/go
RUN cd /opt/go/src/github.com/ottemo/foundation && go get -t 
RUN cd /opt/go/src/github.com/ottemo/foundation && go build && go install
WORKDIR /opt/go/src/github.com/ottemo/foundation

EXPOSE 3000
CMD go run main.go