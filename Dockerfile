FROM golang:1.11.5-alpine3.10
LABEL maintainer="Siddhartha Basu <siddhartha-basu@northwestern.edu>"
ENV GOPROXY https://proxy.golang.org
RUN apk add --no-cache git build-base
RUN mkdir -p /modware-identity
WORKDIR /modware-identity
COPY go.mod ./
COPY go.sum ./
COPY *.go ./
RUN go mod download
ADD server server
ADD commands commands
ADD message message
ADD validate validate
ADD storage storage
RUN go build -o app main.go

FROM alpine:3.10
RUN apk --no-cache add ca-certificates
COPY --from=0 /modware-identity/app /usr/local/bin/
ENTRYPOINT ["/usr/local/bin/app"]
