ARG GO_VERSION=1.22.3
ARG ALPINE_VERSION=3.20

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS builder

WORKDIR /src

COPY go.* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 \
  GOOS=linux \
  go build \
  -ldflags="-s -w" \
  -o yeah \
  ./cmd/yeah/

FROM alpine:${ALPINE_VERSION}

COPY --from=builder /src/yeah /app/yeah

ENTRYPOINT ["/app/yeah"]
CMD ["-l", "-v", "debug"]
