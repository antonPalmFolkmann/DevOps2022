FROM postgres:15.1 AS base
ENV POSTGRES_PASSWORD=password
ENV POSTGRES_USER=postgres
ENV POSTGRES_DB=minitwit
ENV PGPASSWORD=password

COPY ./src/database/dumps/06-01-2022 /restore/dumps
# All bash files and SQL files in /docker-entry-point-initdb.d/ are executed if the docker container is started
# and the database is empty.
COPY ./src/database/restore.sql /docker-entrypoint-initdb.d/restore.sql