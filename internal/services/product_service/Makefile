# RUN
run:
	cd cmd/app && go run .

# GOOSE
# https://github.com/pressly/goose
MIGRATION_NAME = initial
add_migration:
	cd internal/data/migrations && goose create $(MIGRATION_NAME) sql