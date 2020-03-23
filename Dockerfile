FROM golang:latest

WORKDIR /app

COPY go.mod ./

RUN go mod download

# Copy source
COPY . .

RUN go build -o main .

EXPOSE 12346

CMD ["./main"]
