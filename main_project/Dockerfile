FROM golang:latest

RUN mkdir /app
ADD . /app

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files 
COPY go.mod go.sum ./

RUN go mod download

COPY . .
COPY /config/local.env .

# Build the Go app
RUN go build -o main .



# This container exposes port 8080 to the outside world
EXPOSE 8000

# Run the binary program produced by `go install`
CMD ["/app/main"]