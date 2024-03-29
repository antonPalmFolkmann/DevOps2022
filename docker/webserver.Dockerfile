# Bullseye is the latest, stable version as of 09/02/2022
FROM golang:bullseye as base

WORKDIR /src
COPY ./src ./

WORKDIR /client
COPY ./client ./

WORKDIR /src/webserver
RUN go mod tidy
RUN go mod download
RUN go build -o /minitwit

EXPOSE 8080
EXPOSE 8081
ENTRYPOINT [ "/minitwit" ]