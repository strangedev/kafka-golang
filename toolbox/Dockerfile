FROM golang:1.13-alpine

RUN apk --no-cache add --virtual build-deps bash git build-base pkgconfig
WORKDIR /build
RUN git clone https://github.com/edenhill/librdkafka.git \
    && cd librdkafka \
    && ./configure --prefix /usr \
    && make \
    && make install

WORKDIR /go/src/app

CMD ["tail", "-f", "/dev/null"]