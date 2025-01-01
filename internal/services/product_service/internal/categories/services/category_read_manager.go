package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/core/pagination"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/postgres"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/categories/models"
)

type ICategoryManager interface {
	GetCategories(ctx context.Context, pageRequest *pagination.PageRequest) ([]models.Category, int, error)
	GetCategoryById(ctx context.Context, categoryId int) (*models.Category, error)
	GetCategoryByName(ctx context.Context, name string) (*models.Category, error)
	GetCategoriesCount(ctx context.Context) (int, error)

	CreateCategory(ctx context.Context, category *models.Category) error
	UpdateCategory(ctx context.Context, category *models.Category) error
	DeleteCategory(ctx context.Context, categoryId int) error
}

type categoryManager struct {
	db postgres.IPgxDbConn
}

func NewCategoryManager(db postgres.IPgxDbConn) ICategoryManager {
	return categoryManager{
		db: db,
	}
}

func (u categoryManager) GetCategories(ctx context.Context, pageRequest *pagination.PageRequest) ([]models.Category, int, error) {
	query := `SELECT %s FROM categories WHERE is_deleted = false %s %s %s`
	whereExpr := ""
	sortExpr := ""
	paginateExpr := ""
	count := 0

	args := pgx.NamedArgs{}

	if pageRequest != nil {
		if pageRequest.Filters != "" {
			whereExpr = `
				AND (
					normalized_name like @filters
				)
			`

			args["filters"] = fmt.Sprintf("%%%s%%", strings.ToUpper(pageRequest.Filters))
		}

		if pageRequest.Sorting != "" {
			sortingFields := []string{"normalized_name"}
			if err := pageRequest.SanitizeSorting(sortingFields...); err != nil {
				return nil, 0, err
			}

			sortExpr = fmt.Sprintf("ORDER BY %s", pageRequest.Sorting)
		}

		if pageRequest.SkipCount != 0 || pageRequest.MaxResultCount != 0 {
			paginateExpr = "LIMIT @limit OFFSET @offset"
			args["limit"] = pageRequest.MaxResultCount
			args["offset"] = pageRequest.SkipCount
		}

		if err := u.db.QueryRow(ctx, fmt.Sprintf(query, "Count(*)", whereExpr, "", ""), args).Scan(&count); err != nil {
			return nil, 0, fmt.Errorf("unable to count categories: %w", err)
		}
	}

	query = fmt.Sprintf(query, "*", whereExpr, sortExpr, paginateExpr)

	rows, err := u.db.Query(ctx, query, args)
	if err != nil {
		return nil, 0, fmt.Errorf("unable to query categories: %w", err)
	}
	defer rows.Close()

	categories, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Category])

	if count == 0 {
		count = len(categories)
	}

	return categories, count, err
}

func (u categoryManager) GetCategoryById(ctx context.Context, categoryId int) (*models.Category, error) {
	query := `SELECT * FROM categories where id = @categoryId and is_deleted = false LIMIT 1`

	args := pgx.NamedArgs{
		"categoryId": categoryId,
	}
	rows, err := u.db.Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("unable to query category by id: %w", err)
	}
	defer rows.Close()

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Category])
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (u categoryManager) GetCategoryByName(ctx context.Context, name string) (*models.Category, error) {
	query := `SELECT * FROM categories where normalized_name = @name and is_deleted = false LIMIT 1`

	args := pgx.NamedArgs{
		"name": strings.ToUpper(name),
	}
	rows, err := u.db.Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("unable to query category by name: %w", err)
	}
	defer rows.Close()

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Category])
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (u categoryManager) GetCategoriesCount(ctx context.Context) (int, error) {
	query := `SELECT Count(*) FROM categories where is_deleted = false`

	var count int

	if err := u.db.QueryRow(ctx, query).Scan(&count); err != nil {
		return 0, fmt.Errorf("unable to count row: %w", err)
	}

	return count, nil
}
