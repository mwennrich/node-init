FROM golang:1.24-alpine as builder
RUN apk add make binutils

COPY / /work
WORKDIR /work
RUN make

FROM scratch
COPY --from=builder /work/bin/node-init /node-init
ENTRYPOINT ["/node-init","init"]
