############################################################
# Dockerfile to build Ottemo Foundation
# Based on Ubuntu
############################################################

# Use the latest official ubuntu base image from docker hub
FROM ubuntu:latest

RUN apt-get update
RUN apt-get -y install bzr git golang-go gcc
WORKDIR /opt/ottemo/go
RUN mkdir -pv /opt/ottemo/go/src/github.com/ottemo
RUN mkdir -pv /opt/ottemo/go/bin
RUN mkdir -pv /opt/ottemo/go/pkg
RUN mkdir -pv /opt/ottemo/media

RUN git clone https://ottemo-dev:freshbox111222333@github.com/ottemo/foundation.git /opt/ottemo/go/src/github.com/ottemo/foundation
RUN echo "media.fsmedia.folder=/opt/ottemo/media" >> /opt/ottemo/go/bin/ottemo.ini
RUN echo "mongodb.db=ottemo-demo" >> /opt/ottemo/go/bin/ottemo.ini
RUN echo "mongodb.uri=mongodb://ottemo:ottemo2014@candidate.42.mongolayer.com:10243,candidate.43.mongolayer.com:10327/ottemo-demo" >> /opt/ottemo/go/bin/ottemo.ini

ENV GOPATH /opt/ottemo/go
RUN cd /opt/ottemo/go/src/github.com/ottemo/foundation && go get -t 
RUN cd $GOPATH/src/github.com/ottemo/foundation && go get gopkg.in/mgo.v2
cd $GOPATH/src/github.com/ottemo/foundation && go get gopkg.in/mgo.v2/bson
RUN cd /opt/ottemo/go/src/github.com/ottemo/foundation && go build -tags mongo

EXPOSE 3000
WORKDIR /opt/ottemo/go/bin
CMD ./foundation