FROM golang:1.23-alpine as builder

RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o rescope ./cmd/...
RUN chmod +x ./rescope
ENTRYPOINT ["/app/rescope"]