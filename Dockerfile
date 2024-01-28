FROM golang:alpine

ENV GIN_MODE=release
ENV PORT=8080

WORKDIR /go/src/
COPY . .

EXPOSE ${PORT}

RUN go mod download

RUN go build -o main .

CMD ["./main"]
