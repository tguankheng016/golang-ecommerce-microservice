# RUN
run:
	cd cmd/app && go run .

# GOOSE
# https://github.com/pressly/goose
MIGRATION_NAME = added_external_user_id_to_users
add_migration:
	cd internal/data/migrations && goose create $(MIGRATION_NAME) sql