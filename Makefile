mock:
	mockery --all --keeptree

test:
	go test -v ./...

migrate:
	migrate -source file://internal/postgres/migrations \
			-database postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable \
			up
			
rollback:
	migrate -source file://internal/postgres/migrations \
			-database postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable \
			down
			
drop:
	migrate -source file://internal/postgres/migrations \
			-database postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable \
			drop
			
migration:
	@read -p "Enter migration name: " name; \
		migrate create -ext sql -dir internal/postgres/migrations -seq $$name
		
run:
	go run cmd/main.go
