# Choose whatever you want, version >= 1.16
FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN GOOS=linux GOARCH=amd64 go build -o halosuster .

expose 8080

CMD ["./halosuster"]
