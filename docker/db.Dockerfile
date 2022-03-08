FROM postgres:latest AS base
ENV POSTGRES_PASSWORD=password
ENV POSTGRES_USER=postgres
ENV POSTGRES_DB=minitwit
ENV PGPASSWORD=password
COPY ./src/database/dumps/06-01-2022 /restore/dumps
COPY ./src/database/restore.sql /docker-entrypoint-initdb.d/restore.sql