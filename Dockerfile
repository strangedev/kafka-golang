FROM golang:1.13-alpine AS build

# build librdkafka-git
RUN apk --no-cache add --virtual build-deps bash git build-base pkgconfig
WORKDIR /build
RUN git clone https://github.com/edenhill/librdkafka.git \
    && cd librdkafka \
    && ./configure --prefix /usr \
    && make \
    && make install

# build golang sources
WORKDIR /go/src/app
COPY . .
RUN make

# only keep binaries
FROM alpine
COPY --from=build /usr/lib/librdkafka* /usr/lib/
COPY --from=build /usr/include/librdkafka /usr/include/librdkafka
COPY --from=build /go/src/app/bin /usr/local/bin