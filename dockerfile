FROM golang:1.20-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

# Build aplikasi
RUN go build -o main .

EXPOSE 8080

CMD ["./main"]