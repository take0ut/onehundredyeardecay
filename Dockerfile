FROM --platform=linux/amd64 golang:alpine

RUN apk add nodejs npm

WORKDIR /app/src

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY client ./client

COPY cmd/server ./cmd/server

RUN cd client && npm install && npm run build-docker

RUN cp -r client/dist ./cmd/server

CMD cd /app/src && go run cmd/server/main.go
