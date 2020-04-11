FROM golang:1.11

WORKDIR /go/src/app
COPY . .

RUN GO111MODULE=on go get -d -v ./...
RUN GO111MODULE=on go install -v ./...

CMD ["mrkt-api"]