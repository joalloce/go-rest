FROM golang:1.17-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download
RUN go mod verify

COPY *.go ./

RUN go build -o /go-fiber

CMD ["/go-fiber"]