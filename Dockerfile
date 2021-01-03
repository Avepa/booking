FROM golang:latest

WORKDIR /app
COPY ./ /app

WORKDIR /app/src
RUN go mod download
RUN go build -o main .

EXPOSE 9000
CMD ["./main"]