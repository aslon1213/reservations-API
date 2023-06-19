FROM golang:latest
COPY . /app
WORKDIR /app
RUN go mod tidy 
RUN go build -o main .
CMD ["/app/main"]
EXPOSE 8080