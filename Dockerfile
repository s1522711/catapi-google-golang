FROM golang:alpine3.22

LABEL maintainer="s1522711, <s1522711@protonmail.com>"

RUN apk add --no-cache --update curl ca-certificates openssl git tar bash sqlite fontconfig && adduser --disabled-password --home /home/container container

USER container
ENV USER=container HOME=/home/container

WORKDIR /home/container

# Move the go path and go mod cache to the home directory
RUN mkdir -p /home/container/go
RUN mkdir -p /home/container/go/pkg/mod

# Export the environment variables
ENV GOPATH=$HOME/go
ENV GOMODCACHE=$HOME/go/pkg/mod
ENV PATH=$PATH:$GOPATH/bin

COPY ./entrypoint.sh /entrypoint.sh

CMD ["/bin/bash", "/entrypoint.sh"]
