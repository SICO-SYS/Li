FROM golang:alpine

MAINTAINER sine "sinerwr@gmail.com"

RUN apk --update add git && \
    go-wrapper download github.com/SiCo-Ops/Li && \
    apk del git && \
    cd $GOPATH/src/github.com/SiCo-Ops/Li && \
    cp config.sample.json $GOPATH/bin/config.json && \
    go-wrapper install && \
    cd / &&\
    rm -rf $GOPATH/src

WORKDIR $GOPATH/bin

CMD ["Li"]