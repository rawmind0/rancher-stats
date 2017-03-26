FROM rawmind/alpine-base:3.5-1
MAINTAINER Raul Sanchez <rawmind@gmail.com>

#Set environment
ENV SERVICE_NAME=rancher-stats \
    SERVICE_USER=rancher \
    SERVICE_UID=10001 \
    SERVICE_GROUP=rancher \
    SERVICE_GID=10001 \
    GOMAXPROCS=2 \
    GOROOT=/usr/lib/go \
    GOPATH=/opt/src \
    GOBIN=/gopath/bin


# Add files
ADD src /opt/
RUN apk add --no-cache git mercurial bzr make go && \
    cd /opt/src && \
    go get && \
    go build -o rancher-stats && \
    mkdir /opt/rancher-stats && \
    mv rancher-stats /opt/rancher-stats && \
    cd /opt ; rm -rf /opt/src && \
    apk remove --no-cache git mercurial bzr make go


cd ${SERVICE_VOLUME} && \
    tar czvf ${SERVICE_ARCHIVE} * ; rm -rf ${SERVICE_VOLUME}/* 
