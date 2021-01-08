FROM golang:latest

WORKDIR /app
COPY ./ /app

RUN go mod download
WORKDIR /app/cmd
RUN go build -o main .

EXPOSE 9000
CMD ["./main"]