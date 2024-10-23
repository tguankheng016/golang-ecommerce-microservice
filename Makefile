# Run Services
run_identity_service:
	cd internal/services/identity_service/cmd/app && go run .

run_product_service:
	cd internal/services/product_service/cmd/app && go run .

# Run Atlas
# Change migration name before run
MIGRATION_NAME = add_external_user_id_to_users
atlas_identity:
	cd internal/services/identity_service/ && atlas migrate diff $(MIGRATION_NAME) --env gorm

atlas_product:
	cd internal/services/product_service/ && atlas migrate diff $(MIGRATION_NAME) --env gorm

# Run Swaggo
swagger_identity:
	@echo Starting swagger generating
	cd internal/services/identity_service/ && swag init --parseDependency --parseInternal -g cmd/app/main.go -o docs/v1

swagger_product:
	@echo Starting swagger generating
	cd internal/services/product_service/ && swag init --parseDependency --parseInternal -g cmd/app/main.go -o docs/v1

# Run GRPC
proto_identity_user_service:
	@echo Starting proto generating server
	cd internal/services/identity_service/users/grpc_server/protos && protoc --go_out=. --go-grpc_out=. --go-grpc_opt=require_unimplemented_servers=false *.proto
	@echo Starting proto generating client
	cd internal/services/product_service/users/grpc_client/protos && protoc --go_out=. --go-grpc_out=. --go-grpc_opt=require_unimplemented_servers=false *.proto

proto_identity_identity_service:
	@echo Starting proto generating client
	cd internal/pkg/security/jwt/grpc_client/protos && protoc --go_out=. --go-grpc_out=. --go-grpc_opt=require_unimplemented_servers=false *.proto
	cd internal/pkg/permissions/grpc_client/protos && protoc --go_out=. --go-grpc_out=. --go-grpc_opt=require_unimplemented_servers=false *.proto

