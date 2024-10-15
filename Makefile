# Run Services
run_identity_service:
	cd internal/services/identity_service/cmd/app && go run .

# Run Atlas
# Change migration name before run
atlas_identity:
	cd internal/services/identity_service/ && atlas migrate diff migrationName --env gorm

# swagger
swagger_identity:
	@echo Starting swagger generating
	cd internal/services/identity_service/ && swag init --parseDependency --parseInternal -g cmd/app/main.go -o docs