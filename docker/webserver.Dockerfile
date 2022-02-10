# Bullseye is the latest, stable version as of 09/02/2022
FROM golang:bullseye as base
WORKDIR /src
COPY ./go/ ./
RUN go mod download

FROM base as test
CMD ["go", "test"]

FROM base as build
RUN go build -o /minitwit

FROM base as development
EXPOSE 8080
RUN go build -o /minitwit
ENTRYPOINT [ "/minitwit" ]