FROM golang:1.22

WORKDIR /app
RUN go version
ENV $GOPATH=/

COPY . .

RUN go mod download
RUN go build -o coinshop ./cmd/main.go

EXPOSE 8000

CMD ["./musiclibrary"]