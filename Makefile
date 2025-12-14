# Makefile for Go CLI migration (Laravel style)

GO=go
CONSOLE=database/migrate.go

# Migration create
# usage: make migration NAME=create_users_table
migration:
	@$(GO) run $(CONSOLE) create $(NAME)

# Run all pending migrations (up)
migrate-up:
	@$(GO) run $(CONSOLE) up

# Run N pending migrations (up)
# usage: make migrate-up-n N=3
migrate-up-n:
	@$(GO) run $(CONSOLE) up $(N)

# Rollback last migration (down)
migrate-down:
	@$(GO) run $(CONSOLE) down

# Rollback N migrations (down)
# usage: make migrate-down-n N=3
migrate-down-n:
	@$(GO) run $(CONSOLE) down $(N)

# Show migration status
migrate-status:
	@$(GO) run $(CONSOLE) status

.PHONY: migration migrate-up migrate-up-n migrate-down migrate-down-n migrate-status