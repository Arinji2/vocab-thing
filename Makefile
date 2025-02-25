DB_FILE=db/app.db
MIGRATIONS_DIR=db/migrations
GOOSE_DRIVER=sqlite3
GOOSE_BIN=goose

.PHONY: db-status migrate-up migrate-down db-reset db-delete

db-status: ## Show database migration status
	@$(GOOSE_BIN) -dir=$(MIGRATIONS_DIR) $(GOOSE_DRIVER) $(DB_FILE) status

migrate-up: ## Run all pending database migrations
	@$(GOOSE_BIN) -dir=$(MIGRATIONS_DIR) $(GOOSE_DRIVER) $(DB_FILE) up

migrate-down: ## Rollback the last database migration
	@$(GOOSE_BIN) -dir=$(MIGRATIONS_DIR) $(GOOSE_DRIVER) $(DB_FILE) down

db-reset: ## Reset database by dropping all tables and re-migrating
	@$(GOOSE_BIN) -dir=$(MIGRATIONS_DIR) $(GOOSE_DRIVER) $(DB_FILE) reset

db-delete: ## Delete the database file
	@rm -f $(DB_FILE)
	@echo "Database file '$(DB_FILE)' deleted."

