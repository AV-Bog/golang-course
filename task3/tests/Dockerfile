FROM golang:1.25

WORKDIR /tests
COPY . .

RUN go mod download

CMD ["go", "test", "-v", "."]
