# Stage 1 (Build)
FROM golang:1.21-alpine AS builder

ARG VERSION
RUN apk add --update --no-cache git make
WORKDIR /app/
COPY go.mod go.sum /app/
RUN go mod download
COPY . /app/
RUN CGO_ENABLED=0 go build \
    -v \
    -trimpath \
    -o autoportforward \
    cmd/main.go
RUN echo "ID=\"distroless\"" > /etc/os-release
# Stage 2 (Final)
FROM gcr.io/distroless/static:latest
COPY --from=builder /etc/os-release /etc/os-release

COPY --from=builder /app/autoportforward  /usr/bin/
CMD [ "/usr/bin/autoportforward" ]

EXPOSE 8080