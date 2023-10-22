FROM golang:alpine

RUN apk add nodejs npm

WORKDIR /app/src

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY client ./client

COPY cmd/server ./cmd

RUN ls cmd

RUN cd client && npm install && npm run build

CMD cd .. && go run cmd/main.go
