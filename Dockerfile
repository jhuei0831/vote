FROM golang:latest

RUN mkdir -p /usr/local/go/src/vote
WORKDIR /usr/local/go/src/vote
ADD . /usr/local/go/src/vote

RUN go mod download
RUN go build ./main.go

EXPOSE 9443
CMD ["./main"]