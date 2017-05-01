# build stage
FROM golang:alpine AS build-env
ADD . /src
RUN cd /src && go build -o nixlight

# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /src/nixlight /app/
EXPOSE 80
ENTRYPOINT ./nixlight