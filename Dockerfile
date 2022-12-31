FROM golang:latest

WORKDIR /app

COPY go.mod ./

RUN go mod download

# Copy source
COPY . .

RUN go build -o main cmd/fakeovirt/main.go

ENV SERVICE=ovirt-imageio
ENV PORT=30001
EXPOSE 30001

CMD ["./main"]
