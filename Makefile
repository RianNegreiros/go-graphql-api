mock:
	mockery --all --keeptree

test:
	go test -v ./...