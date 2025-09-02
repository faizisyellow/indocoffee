include .env

MIGRATE_PATH = cmd/migrate/migrations

.PHONY: new-migration
migrate:
	@migrate create -ext sql -dir $(MIGRATE_PATH) -seq $(name)

.PHONY:migration-up
migrate-up:
	@migrate -path=$(MIGRATE_PATH) -database=mysql://root:$(DB_MIGRATE_AUTH)@/$(DATABASE_NAME) up

.PHONY:migration-down
migrate-down:
	@migrate -path=$(MIGRATE_PATH) -database=mysql://root:$(DB_MIGRATE_AUTH)@/$(DATABASE_NAME) down

.PHONY:migration-back
migrate-back:
	@migrate -path=$(MIGRATE_PATH) -database=mysql://root:$(DB_MIGRATE_AUTH)@/$(DATABASE_NAME) force $(no)


.PHONY: gen-docs
gen-docs:
	@swag init -g /main.go -d ./cmd/api,internal && swag fmt
