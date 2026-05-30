FROM golang:1.26.2-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build with debug symbols
RUN go build -gcflags="all=-N -l" -o main cmd/api/main.go

# Install Delve
RUN go install github.com/go-delve/delve/cmd/dlv@latest


FROM golang:1.26.2-alpine AS debug

WORKDIR /app

COPY --from=build /app .
COPY --from=build /go/bin/dlv /usr/local/bin/dlv

EXPOSE 8080
EXPOSE 40000

CMD ["dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "./main"]