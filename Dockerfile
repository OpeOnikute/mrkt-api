FROM golang:1.11

WORKDIR /go/src/app
COPY . .

RUN GO111MODULE=on

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["mrkt-api"]