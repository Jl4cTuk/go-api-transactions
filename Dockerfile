FROM golang:1.24.4

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -C . -o base ./cmd/qual

EXPOSE 8080

CMD ["/build/base"]