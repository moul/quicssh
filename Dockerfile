# build
FROM            golang:1.19.0-alpine as builder
RUN             apk add --no-cache git gcc musl-dev make
ENV             GO111MODULE=on
WORKDIR         /go/src/moul.io/quicssh
COPY            go.* ./
RUN             go mod download
COPY            . ./
RUN             make install

# minimalist runtime
FROM            alpine:3.16.0
COPY            --from=builder /go/bin/quicssh /bin/
ENTRYPOINT      ["/bin/quicssh"]
CMD             []
