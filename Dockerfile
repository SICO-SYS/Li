FROM golang:alpine

MAINTAINER sine "sinerwr@gmail.com"

RUN apk --update add git
RUN go-wrapper download github.com/SiCo-Ops/He
RUN apk del git

WORKDIR $GOPATH/src/github.com/SiCo-Ops/He

RUN go-wrapper install

WORKDIR $GOPATH/bin/

RUN rm -rf $GOPATH/src

EXPOSE 6666

VOLUME $GOPATH/bin/config.json

CMD Li
