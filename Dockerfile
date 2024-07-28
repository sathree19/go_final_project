FROM golang:1.21.0


ENV TODO_PORT=:7540

ENV TODO_DBFILE=dbS/scheduler.db

ENV TODO_PASSWORD=12345

WORKDIR /usr/src/app

COPY . .

RUN go mod download 

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /my_app

EXPOSE 7540

CMD ["/my_app"]
