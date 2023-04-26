FROM golang:1.20

WORKDIR /app

COPY ./ /app/
# install psql
RUN apt-get update
RUN apt-get -y install postgresql-client
# make wait for postgres
RUN chmod +x wait-for-postgres.sh
# build go app
RUN go mod download
RUN go build -o ./bin/api

CMD [ "./bin/api" ]