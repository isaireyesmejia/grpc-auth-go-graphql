FROM golang:1.27-alpine AS constructor

RUN apk add --no-cache protobuf protobuf-dev

WORKDIR /app

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN protoc --go_out=. --go_opt=paths=source_relative \
           --go-grpc_out=. --go-grpc_opt=paths=source_relative \
           proto/auth.proto

RUN CGO_ENABLED=0 GOOS=linux go build -o servidor ./server/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=constructor /app/servidor .
EXPOSE 50051
CMD ["./servidor"]
