migrate:
	@go run ./cmd/migrator --migrations-path=./migrations --storage-path=./storage/sso.db