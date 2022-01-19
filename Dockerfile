FROM golang:latest

WORKDIR /app

COPY ./ /app

run go mod download -x

run go get github.com/githubnemo/CompileDaemon

ENTRYPOINT CompileDaemon -exclude-dir=.git -exclude-dir=docs --build="go build main.go"