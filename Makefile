# Run Services
run_identity_service:
	cd internal/services/identity_service/ && go run ./cmd/app/main.go

# Run Atlas
atlas_identity:
	cd internal/services/identity_service/ && atlas migrate diff initial --env gorm
