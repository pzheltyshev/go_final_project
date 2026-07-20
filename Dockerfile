FROM golang:1.25.0

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o todo-app .

EXPOSE ${TODO_PORT}

CMD ["./todo-app"] 