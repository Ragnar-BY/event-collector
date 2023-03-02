
FROM golang:1.20 as builder

RUN useradd server

RUN mkdir -p /server
COPY . /server

WORKDIR /server


RUN go mod download


# Build the binary with go build
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
    go build -ldflags '-s -w -extldflags "-static"' \
    -o /bin/server ./cmd/server


FROM alpine:latest

COPY --from=builder /bin/server /server

ENTRYPOINT ["/server"]