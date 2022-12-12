FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

COPY *.go ./

RUN go mod download
RUN go mod verify

RUN go build -o /main

CMD ["/main"]