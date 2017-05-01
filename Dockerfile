# build stage
FROM golang:alpine AS build-env
WORKDIR /go
ENV GOPATH=/go
ADD . /go/src/github.com/fishnix/nixlight/
RUN go build -o nixlight src/github.com/fishnix/nixlight/nixlight.go

# final stage
FROM alpine
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build-env nixlight /app/nixlight
EXPOSE 80
ENTRYPOINT ["./nixlight"]
