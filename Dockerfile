FROM golang:latest

WORKDIR /app

COPY ./ /app

ENV GOPATH=/go
ENV PATH=${GOPATH}/bin:${PATH}

RUN go mod download -x

RUN go get github.com/githubnemo/CompileDaemon

RUN cd ${GOPATH} \
    GO111MODULE=on go get -u -v github.com/ipfs/ipfs-update \
    && ipfs-update install latest \
    || echo "ERROR: ===== ipfs-update failing again, will run without enabling IPFS ====="

ENTRYPOINT CompileDaemon -exclude-dir=.git -exclude-dir=docs --build="go build main.go"
