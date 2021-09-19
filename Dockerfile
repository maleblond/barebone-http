FROM golang:1.17-alpine

EXPOSE 3000

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY *.go ./

RUN go build -o ./bin

CMD [ "./bin" ]
