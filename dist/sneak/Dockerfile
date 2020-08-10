ARG SNEAK_VERSION

FROM golang:alpine
ARG SNEAK_VERSION
ENV GO111MODULE=on

RUN apk update \
    && apk add gcc libc-dev

WORKDIR /go/src/github.com/kathleenfrench/sneak
COPY . .

RUN go build -mod=readonly -a -ldflags "-X github.com/kathleenfrench/sneak.Version=$SNEAK_VERSION" -o sneak

ENTRYPOINT [ "./sneak" ]