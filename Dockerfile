FROM rawmind/alpine-base:3.5-1
MAINTAINER Raul Sanchez <rawmind@gmail.com>

#Set environment
ENV SERVICE_NAME=rancher-stats \
    SERVICE_HOME=/opt/rancher-stats \
    SERVICE_USER=rancher \
    SERVICE_UID=10012 \
    SERVICE_GROUP=rancher \
    SERVICE_GID=10012 \
    GOMAXPROCS=2 \
    GOROOT=/usr/lib/go \
    GOPATH=/opt/src \
    GOBIN=/gopath/bin
ENV PATH=${PATH}:${SERVICE_HOME}

# Add files
ADD src /opt/src/
RUN apk add --no-cache git mercurial bzr make go musl-dev && \
    cd /opt/src && \
    go get && \
    go build -o rancher-stats && \
    mkdir ${SERVICE_HOME} && \
    mv rancher-stats ${SERVICE_HOME}/ && \
    cd /opt && \ 
    rm -rf /opt/src && \
    apk del --no-cache git mercurial bzr make go musl-dev

USER $SERVICE_USER
WORKDIR $SERVICE_HOME

