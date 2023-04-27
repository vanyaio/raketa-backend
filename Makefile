build:
	@go build -o ./bin/api

run: build
	@./bin/api 

run-ports:
	go run main.go -grpc-port=$(grpc-port) -rest-port=$(rest-port)

postgres-up:
	docker run --name raketadb \
	-e POSTGRES_HOST=localhost \
	-e POSTGRES_PASSWORD=postgres \
	-e POSTGRES_USER=postgres \
	-e POSTGRES_DB=raketadb \
	-p 5432:5432 -d postgres

postgres-start:
	docker start raketadb

postgres-run: postgres-start
	docker exec -it raketadb psql -U postgres raketadb

postgres-stop:
	docker stop raketadb

postgres-del: postgres-stop
	docker rm raketadb

protob:
	@protoc -I ./proto --go_out=./proto --go_opt=paths=source_relative \
	--go-grpc_out=./proto --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=./proto --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out ./swagger \
    --openapiv2_opt logtostderr=true \
	proto/*.proto

# migrations
migrate-create:
	migrate create -ext sql -dir ./migrations -seq raketadb

migrate-up:
	migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/raketadb?sslmode=disable" up

migrate-down:
	migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/raketadb?sslmode=disable" down ${version}

# docker-migrations
docker-migrate-down:
	docker-compose exec raketa bash -c "migrate -path /app/migrations -database 'postgres://postgres:postgres@raketadb:5432/raketadb?sslmode=disable' down ${version}"

evans:
	evans -r repl -p 50052

test:
	@go test -v ./...