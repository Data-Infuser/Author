FROM amd64/golang:1.14

ENV GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /go/src/author

COPY go.sum go.mod ./

RUN go mod download

COPY . .

# Build the Go app
RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 9090

# Command to run the executable
CMD ["./main"]
