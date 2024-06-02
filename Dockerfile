# syntax=docker/dockerfile:1

FROM golang:latest
ENV HOST 0.0.0.0

RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y git

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

RUN pwd
RUN mkdir ./src/github.com && mkdir /go/src/github.com/raspincel
COPY go.mod go.sum /
RUN go get -u github.com/golang/dep/cmd/dep \
&&  git clone https://github.com/raspincel/es-ii-backend /go/src/github.com/raspincel/es-ii-backend

WORKDIR ./src/github.com/raspincel/es-ii-backend
RUN cp /go/bin/dep ./
COPY internal ./

RUN go get -u github.com/golang/dep/cmd/dep
RUN ./dep init
COPY go.mod go.sum ./
RUN go mod tidy
RUN ./dep ensure -v 
RUN go mod vendor
RUN pwd
RUN go build -o main

EXPOSE 8080
RUN ls .
CMD ["/go/src/github.com/raspincel/es-ii-backend/main"]