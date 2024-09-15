FROM golang:latest

WORKDIR /go/src/app

COPY . .

RUN go build -o main ./src/cmd

EXPOSE 8080

CMD ["./main"]