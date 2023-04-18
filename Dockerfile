FROM golang:alpine

RUN mkdir /files
COPY storage.go /files
WORKDIR /files

RUN go build -o /files/storage storage.go
ENTRYPOINT ["/files/storage"]
