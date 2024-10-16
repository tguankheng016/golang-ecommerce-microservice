# Run Services
run_identity_service:
	cd internal/services/identity_service/cmd/app && go run .

run_product_service:
	cd internal/services/product_service/cmd/app && go run .

# Run Atlas
# Change migration name before run
atlas_identity:
	cd internal/services/identity_service/ && atlas migrate diff migrationName --env gorm

atlas_product:
	cd internal/services/product_service/ && atlas migrate diff initial --env gorm

# Run Swaggo
swagger_identity:
	@echo Starting swagger generating
	cd internal/services/identity_service/ && swag init --parseDependency --parseInternal -g cmd/app/main.go -o docs

swagger_product:
	@echo Starting swagger generating
	cd internal/services/product_service/ && swag init --parseDependency --parseInternal -g cmd/app/main.go -o docs