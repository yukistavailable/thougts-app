FROM golang:1.13-alpine3.11

RUN apk add git \
    && git config --global http.postBuffer 16M \
    && go get -u "github.com/gin-gonic/gin" \
    && go get "github.com/jinzhu/gorm/dialects/postgres" \
    && go get "github.com/jinzhu/gorm" \
    && go get gopkg.in/ini.v1 \
    && go get gopkg.in/gorp.v2

WORKDIR $GOPATH/src/thoughts-app/api