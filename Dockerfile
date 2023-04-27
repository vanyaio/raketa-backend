FROM golang:1.20

WORKDIR /app

COPY ./ /app/

ENV PATH="$PATH:/usr/local/bin" \
    GRPC_PORT=:50052 \
    REST_PORT=:9090 \
    POSTGRES_HOST=raketadb \
    POSTGRES_DB=raketadb
# install psql
RUN apt-get update
RUN apt-get -y install postgresql-client

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
# make wait for postgres
RUN chmod +x wait-for-postgres.sh
# build go app
RUN go mod download
RUN go build -o ./bin/api

CMD ["sh", "-c", "./wait-for-postgres.sh raketadb migrate -path /app/migrations -database 'postgres://postgres:postgres@raketadb:5432/raketadb?sslmode=disable' up && ./bin/api"]
