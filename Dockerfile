FROM golang:1.22.4-alpine as builder

WORKDIR app

# COPY go.mod, go.sum and download the dependencies
COPY go.* ./
RUN go mod download

# COPY All things inside the project and build
COPY . .
RUN go build -o /app/build/mychirpy .

FROM alpine:latest
COPY --from=builder /app/build/mychirpy /app/build/mychirpy

EXPOSE 8080
ENTRYPOINT [ "/app/build/mychirpy" ]