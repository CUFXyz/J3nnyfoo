FROM golang

WORKDIR /fooapp

COPY go.mod /fooapp
COPY go.sum /fooapp

RUN go mod tidy

COPY . /fooapp/

ENV port=9090
EXPOSE 9090

RUN go build -o foo_app cmd/main.go

CMD ["./foo_app"]