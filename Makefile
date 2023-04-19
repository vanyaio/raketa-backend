build:
	@go build -o ./bin/api

run: build
	@./bin/api

docker:
	docker run --name raketadb \
	-e POSTGRES_HOST=localhost \
	-e POSTGRES_PASSWORD=postgres \
	-e POSTGRES_USER=postgres \
	-e POSTGRES_DB=raketadb \
	-p 5432:5432 -d postgres

docker-start:
	docker start raketadb

docker-exec: docker-start
	docker exec -it raketadb psql -U postgres raketadb

protob:
	@protoc -I ./proto --go_out=./proto --go_opt=paths=source_relative \
	--go-grpc_out=./proto --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=./proto --grpc-gateway_opt=paths=source_relative \
	proto/*.proto

# migrations
migrate-create:
	migrate create -ext sql -dir ./migrations -seq raketadb

migrate-up:
	migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/raketadb?sslmode=disable" up

migrate-down:
	migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/raketadb?sslmode=disable" down

evans:
	evans -r repl -p 50052

test:
	@go test -v ./...