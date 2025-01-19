#syntax=docker/dockerfile:1

FROM golang:1.23 as builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set destination for COPY
WORKDIR /build

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build
RUN go build -o receipt-app main.go

WORKDIR /dist

RUN cp /build/receipt-app .

FROM golang:1.23

COPY --from=builder /dist/receipt-app /

# Expose ports that app uses
EXPOSE 8080

# Run
CMD ["/receipt-app"]