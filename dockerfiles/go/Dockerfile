FROM golang:1.13-alpine3.11

RUN apk add git \
    && go get -u "github.com/gin-gonic/gin" \
    && go get "github.com/jinzhu/gorm/dialects/postgres" \
    && go get "github.com/jinzhu/gorm"