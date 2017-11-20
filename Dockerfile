FROM alpine

COPY . /var/toxy

WORKDIR /var/toxy

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && apk add --update go git libc-dev \
    && export GOPATH=`pwd` \
    && export GOBIN=/usr/bin \
    && go get \
    && go install \
    && echo $GOPATH \
    && rm -r /var/toxy \
    && apk del --purge go git libc-dev

WORKDIR /var/config
