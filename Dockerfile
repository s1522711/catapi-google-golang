FROM golang:alpine3.22

LABEL maintainer="s1522711, <s1522711@protonmail.com>"

RUN apk add --no-cache --update curl ca-certificates openssl git tar bash sqlite fontconfig && adduser --disabled-password --home /home/container container

USER container
ENV USER=container HOME=/home/container

WORKDIR /home/container

COPY ./entrypoint.sh /entrypoint.sh

CMD ["/bin/bash", "/entrypoint.sh"]
