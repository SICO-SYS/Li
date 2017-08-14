FROM golang:alpine

MAINTAINER sine "sinerwr@gmail.com"

RUN apk --update add git && \
    go-wrapper download github.com/SiCo-Ops/Li && \
    apk del git && \
    cd $GOPATH/src/github.com/SiCo-Ops/Li && \
    go-wrapper install && \
    rm -rf $GOPATH/src

EXPOSE 6666

VOLUME $GOPATH/bin/config.json

CMD Li
