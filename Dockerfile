FROM golang:latest AS builder

WORKDIR /build
ENV CGO_ENABLED=0

# Populate module cache
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o app .

FROM scratch
COPY --from=builder /build /

ENV GIN_MODE=debug

ENTRYPOINT ["/app"]
