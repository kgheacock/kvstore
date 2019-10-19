FROM golang

WORKDIR /app
ADD . /go/src/github.com/colbyleiske/cse138_assignment2/
RUN go get -v github.com/gorilla/mux
RUN go install github.com/colbyleiske/cse138_assignment2/
ENTRYPOINT /go/bin/cse138_assignment2
