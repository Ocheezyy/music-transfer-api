FROM golang:1.23-bookworm

WORKDIR /go/src/app

COPY . .

RUN rm -rf .vscode && rm -rf .example.env && rm -rf ReadMe.md && rm -rf clear_test_data.sql && rm -rf tests.rest && rm -rf migrate

RUN go build -o main main.go

CMD ["./main"]