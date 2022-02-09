# Bullseye is the latest, stable version as of 09/02/2022
FROM golang:bullseye

WORKDIR /src

COPY ./go/ ./
RUN go mod download

RUN go build -o /minitwit

EXPOSE 8080

ENTRYPOINT [ "/minitwit" ]