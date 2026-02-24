FROM golang:1.25-alpine AS builder

WORKDIR /build

RUN apk add --no-cache ca-certificates build-base git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -o manager \
    -ldflags="-s -w" \
    -trimpath \
    ./cmd/main.go

FROM gcr.io/distroless/static:nonroot

LABEL maintainer="operator-maintainers@example.com"
LABEL description="Kubernetes operator for agentic workloads"

COPY --from=builder /build/manager /manager

USER nonroot:nonroot

ENTRYPOINT ["/manager"]
