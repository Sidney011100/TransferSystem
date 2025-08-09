FROM golang:1.23.5

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o server .

CMD ["./server"]