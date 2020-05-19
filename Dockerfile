FROM golang:latest

WORKDIR /app

COPY go.mod ./

RUN go mod download

# Copy source
COPY . .

RUN go build -o main cmd/fakeovirt/main.go

EXPOSE 12346

CMD ["./main"]
