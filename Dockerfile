FROM golang:latest

WORKDIR /app

COPY go.mod go.sum .env ./
RUN go mod download && go mod verify

COPY . .

RUN go build -o /gofiber

EXPOSE 3000

CMD [ "gofiber" ]