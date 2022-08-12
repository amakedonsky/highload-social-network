
FROM golang:alpine AS builder

ENV CGO_ENABLED 0
ENV GOOS linux

WORKDIR /build
ADD go.mod .
RUN go mod download
COPY . .
RUN go build -o main main.go

FROM alpine
WORKDIR /app
COPY --from=builder /build/main /app/main

EXPOSE 8090
CMD ["/app/main"]