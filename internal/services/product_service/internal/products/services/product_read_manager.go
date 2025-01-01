package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/core/pagination"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/postgres"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/products/models"
)

type IProductManager interface {
	GetProductsWithCategory(ctx context.Context, pageRequest *pagination.PageRequest) ([]models.ProductWithCategory, int, error)
	GetProductById(ctx context.Context, productId int) (*models.Product, error)
	GetProductsCount(ctx context.Context) (int, error)

	CreateProduct(ctx context.Context, product *models.Product) error
	UpdateProduct(ctx context.Context, product *models.Product) error
	DeleteProduct(ctx context.Context, productId int) error
}

type productManager struct {
	db postgres.IPgxDbConn
}

func NewProductManager(db postgres.IPgxDbConn) IProductManager {
	return productManager{
		db: db,
	}
}

func (u productManager) GetProductsWithCategory(ctx context.Context, pageRequest *pagination.PageRequest) ([]models.ProductWithCategory, int, error) {
	query := `SELECT %s FROM products p join categories c on p.category_id = c.id WHERE p.is_deleted = false %s %s %s`
	whereExpr := ""
	sortExpr := ""
	paginateExpr := ""
	count := 0

	args := pgx.NamedArgs{}

	if pageRequest != nil {
		if pageRequest.Filters != "" {
			whereExpr = `
				AND (
					p.normalized_name like @filters OR
					p.normalized_description like @filters
				)
			`

			args["filters"] = fmt.Sprintf("%%%s%%", strings.ToUpper(pageRequest.Filters))
		}

		if pageRequest.Sorting != "" {
			sortingFields := []string{"normalized_name", "normalized_description"}
			if err := pageRequest.SanitizeSorting(sortingFields...); err != nil {
				return nil, 0, err
			}

			sortExpr = fmt.Sprintf("ORDER BY p.%s", pageRequest.Sorting)
		}

		if pageRequest.SkipCount != 0 || pageRequest.MaxResultCount != 0 {
			paginateExpr = "LIMIT @limit OFFSET @offset"
			args["limit"] = pageRequest.MaxResultCount
			args["offset"] = pageRequest.SkipCount
		}

		if err := u.db.QueryRow(ctx, fmt.Sprintf(query, "Count(*)", whereExpr, "", ""), args).Scan(&count); err != nil {
			return nil, 0, fmt.Errorf("unable to count products: %w", err)
		}
	}

	selectQueryExpr := `
		p.id, p.name, p.description, p.price, p.stock_quantity, 
		p.category_id, p.created_at, p.created_by, p.updated_at, 
		p.updated_by, c.id, c.name`

	query = fmt.Sprintf(query, selectQueryExpr, whereExpr, sortExpr, paginateExpr)

	rows, err := u.db.Query(ctx, query, args)
	if err != nil {
		return nil, 0, fmt.Errorf("unable to query products: %w", err)
	}
	defer rows.Close()

	var products []models.ProductWithCategory

	for rows.Next() {
		var p models.ProductWithCategory
		// Scan the row into the struct fields
		err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.Description,
			&p.Price,
			&p.StockQuantity,
			&p.CategoryId,
			&p.CreatedAt,
			&p.CreatedBy,
			&p.UpdatedAt,
			&p.UpdatedBy,
			&p.CategoryFK.Id,
			&p.CategoryFK.Name,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("unable to scan products: %w", err)
		}

		products = append(products, p)
	}

	if count == 0 {
		count = len(products)
	}

	return products, count, err
}

func (u productManager) GetProductById(ctx context.Context, productId int) (*models.Product, error) {
	query := `SELECT * FROM products where id = @productId and is_deleted = false LIMIT 1`

	args := pgx.NamedArgs{
		"productId": productId,
	}
	rows, err := u.db.Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("unable to query product by id: %w", err)
	}
	defer rows.Close()

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Product])
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (u productManager) GetProductsCount(ctx context.Context) (int, error) {
	query := `SELECT Count(*) FROM products where is_deleted = false`

	var count int

	if err := u.db.QueryRow(ctx, query).Scan(&count); err != nil {
		return 0, fmt.Errorf("unable to count row: %w", err)
	}

	return count, nil
}
