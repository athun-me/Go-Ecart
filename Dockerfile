FROM golang:latest

WORKDIR /app

COPY . /app

COPY go.mod go.sum ./

RUN go mod download

RUN go build -o myapp .

EXPOSE $PORT

CMD ["./myapp"]
