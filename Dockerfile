FROM golang:alpine

WORKDIR /app
RUN apk add --no-cache git
RUN go get -v github.com/gorilla/mux
COPY . /go/src/github.com/colbyleiske/cse138_assignment2/
RUN go install github.com/colbyleiske/cse138_assignment2/
EXPOSE 13800
ENTRYPOINT /go/bin/cse138_assignment2