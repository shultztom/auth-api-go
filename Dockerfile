FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

# Run tests during build
RUN go test ./... -v

# Build binary only if tests pass
RUN go build -o /main

RUN go build -o /main

EXPOSE 8080

CMD [ "/main" ]