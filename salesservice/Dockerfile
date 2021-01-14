FROM golang:1.15.6-alpine3.12 as builder
RUN apk add git

RUN mkdir /app
ADD . /app
WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates curl

RUN mkdir /app

WORKDIR /app/

COPY --from=builder /app/main ./cmd/sales/main.go

EXPOSE 8001

CMD ["./main"]
