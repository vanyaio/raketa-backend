run:
	POSTGRES_HOST=localhost POSTGRES_DB=raketalocaldb go run main.go	

postgres-up:
	docker run --name raketalocaldb \
	-e POSTGRES_HOST=localhost \
	-e POSTGRES_PASSWORD=postgres \
	-e POSTGRES_USER=postgres \
	-e POSTGRES_DB=raketalocaldb \
	-p 5436:5432 -d postgres

postgres-local-start:
	docker start raketalocaldb

postgres-local-run: postgres-local-start
	docker exec -it raketalocaldb psql -U postgres raketalocaldb

postgres-run:
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
	migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/raketalocaldb?sslmode=disable" up

migrate-down:
	migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/raketalocaldb?sslmode=disable" down ${version}

# docker-migrations
docker-migrate-down:
	docker-compose exec raketa bash -c "migrate -path /app/migrations -database 'postgres://postgres:postgres@raketadb:5432/raketadb?sslmode=disable' down ${version}"

evans:
	evans -r repl -p 50052

test:
	@go test -v ./...