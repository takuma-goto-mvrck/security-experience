FROM golang:1.9
ARG service
RUN go get github.com/go-sql-driver/mysql
RUN mkdir -p go/src/github.com/takuma-goto310/security-experience/Apps/$service
WORKDIR /go/src/github.com/takuma-goto310/security-experience/Apps/$service