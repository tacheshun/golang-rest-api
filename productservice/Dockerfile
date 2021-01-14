FROM golang:1.15.6-alpine3.12

RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o product .
CMD ["/app/product"]
