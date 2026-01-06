FROM docker.io/golang:1.24-alpine AS builder

# Enable CGO and install sqlite deps
RUN apk add --no-cache gcc musl-dev sqlite-dev git ca-certificates

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

RUN go build -o /out/api .

# ---------- runtime ----------
FROM alpine:3.20

RUN apk add --no-cache ca-certificates sqlite-libs

WORKDIR /app

COPY --from=builder /out/api /app/api
RUN mkdir -p /app/storage

EXPOSE 8080
CMD ["/app/api"]
