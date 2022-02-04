FROM golang:1.16-alpine
RUN apk add git
RUN apk add --no-cache make
RUN mkdir /app
ADD . /app
WORKDIR /app/cmd/myapp
RUN go build -o main main.go
CMD ["/app/cmd/myapp/main"]