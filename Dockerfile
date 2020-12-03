FROM golang:latest

WORKDIR /go/src/app
COPY . .

RUN go build cmd/ad-hoc/main.go
# this is necessary otherwise the CMD won't be able to run the binary
RUN mv main ad-hoc.sh

CMD ["./ad-hoc.sh"]
