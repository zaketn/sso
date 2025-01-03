migrate:
	@go run ./cmd/migrator --migrations-path=./migrations --storage-path=./storage/sso.db

migrate-tests:
	@go run ./cmd/migrator \
	--migrations-path=./tests/migrations \
	--migrations-table=migrations_tests \
	--storage-path=./storage/sso.db