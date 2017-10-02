FROM golang:alpine

MAINTAINER sine "sinerwr@gmail.com"

RUN apk --update add git && \
    go-wrapper download github.com/SiCo-Ops/Li && \
    apk del git && \
    cd $GOPATH/src/github.com/SiCo-Ops/Li && \
    go-wrapper install && \
    cp *.json $GOPATH/bin/ && \
    cd / &&\
    rm -rf $GOPATH/src

WORKDIR $GOPATH/bin

VOLUME $GOPATH/bin

EXPOSE 6666

CMD ["Li"]