FROM golang:latest 

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o myapp

EXPOSE 9001

CMD ["./myapp"]
