FROM golang:latest

WORKDIR /app

COPY ./ /app

ENV GOPATH=/go
ENV PATH=${GOPATH}/bin:${PATH}

RUN go mod download -x

RUN go get github.com/githubnemo/CompileDaemon

ENTRYPOINT CompileDaemon -exclude-dir=.git -exclude-dir=docs --build="go build main.go"
