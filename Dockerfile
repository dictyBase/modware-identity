FROM golang:1.9.5-alpine3.7
LABEL maintainer="Siddhartha Basu <siddhartha-basu@northwestern.edu>"
RUN apk add --no-cache git build-base \
    && go get github.com/golang/dep/cmd/dep
RUN mkdir -p /go/src/github.com/dictyBase/modware-identity
WORKDIR /go/src/github.com/dictyBase/modware-identity
COPY Gopkg.* main.go ./
ADD server server
ADD commands commands
ADD message message
ADD validate validate
ADD storage storage
RUN dep ensure \
    && go build -o app

FROM alpine:3.7
RUN apk --no-cache add ca-certificates
COPY --from=0 /go/src/github.com/dictyBase/modware-identity/app /usr/local/bin/
ENTRYPOINT ["/usr/local/bin/app"]
