FROM golang:alpine AS build-env
ENV GOPATH ./

FROM alpine
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN update-ca-certificates

WORKDIR /app/
COPY .env /app
ADD ./app /app/
RUN chmod +x /app/app

ENTRYPOINT ["./app"]