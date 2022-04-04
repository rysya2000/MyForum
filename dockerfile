FROM golang

LABEL project="Forum"

WORKDIR /web

COPY . .

RUN go build cmd/main.go

CMD ["./main"]