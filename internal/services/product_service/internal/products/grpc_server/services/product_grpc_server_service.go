package services

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/postgres"
	productGrpc "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/products/grpc_server/protos"
	productModel "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/products/models"
)

type ProductGrpcServerService struct {
	db postgres.IPgxDbConn
}

func NewProductGrpcServerService(db *pgxpool.Pool) *ProductGrpcServerService {
	return &ProductGrpcServerService{db: db}
}

func (u *ProductGrpcServerService) GetAllProducts(ctx context.Context, req *productGrpc.GetAllProductsRequest) (*productGrpc.GetAllProductsResponse, error) {
	query := `SELECT * FROM products WHERE is_deleted = false`

	rows, err := u.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("unable to query products: %w", err)
	}
	defer rows.Close()

	products, err := pgx.CollectRows(rows, pgx.RowToStructByName[productModel.Product])
	if err != nil {
		return nil, err
	}

	var grpcProducts []*productGrpc.ProductModel

	for _, product := range products {
		grpcProduct := &productGrpc.ProductModel{
			Id:            int32(product.Id),
			Name:          product.Name,
			Description:   product.Description,
			Price:         product.Price.String(),
			StockQuantity: int32(product.StockQuantity),
		}

		grpcProducts = append(grpcProducts, grpcProduct)
	}

	var result = &productGrpc.GetAllProductsResponse{
		Products: grpcProducts,
	}

	return result, nil
}
