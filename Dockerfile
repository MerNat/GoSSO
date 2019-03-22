FROM golang:alpine AS base_alpine

RUN apk add bash ca-certificates git gcc g++ libc-dev

WORKDIR /temp

ADD . /temp

RUN go build -v -o GoSSO && mkdir /final \
    && cp -r /temp/GoSSO /final \
    && cp -r /temp/config.json /final

FROM alpine

RUN apk update && apk add ca-certificates \
    && rm -rf /var/cache/apk/*

WORKDIR /mySSO

COPY --from=base_alpine /final /mySSO

EXPOSE 8080

ENTRYPOINT ./GoSSO