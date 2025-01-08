package v1

import (
	"context"
	"encoding/json"
	"net/http"

	v "github.com/RussellLuo/validating/v3"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jinzhu/copier"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/events"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/permissions"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/postgres"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/products/dtos"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/products/models"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/products/services"
)

// Request
type HumaCreateProductRequest struct {
	Body struct {
		dtos.CreateProductDto
	}
}

// Result
type HumaCreateProductResult struct {
	Body struct {
		Product dtos.ProductDto
	}
}

// Validator
func (e HumaCreateProductRequest) Schema() v.Schema {
	return v.Schema{
		v.F("id", e.Body.Id): v.Any(
			v.Zero[*int](),
			v.Nested(func(ptr *int) v.Validator { return v.Value(*ptr, v.Eq(0).Msg("Invalid product id")) }),
		).LastError(),
		v.F("product_name", e.Body.Name):        v.Nonzero[string]().Msg("Please enter the product name"),
		v.F("product_desc", e.Body.Description): v.Nonzero[string]().Msg("Please enter the product description"),
	}
}

// Handler
func MapRoute(
	api huma.API,
	pool *pgxpool.Pool,
	publisher message.Publisher,
) {
	huma.Register(
		api,
		huma.Operation{
			OperationID:   "CreateProduct",
			Method:        http.MethodPost,
			Path:          "/products/product",
			Summary:       "Create Product",
			Tags:          []string{"Products"},
			DefaultStatus: http.StatusOK,
			Security: []map[string][]string{
				{"bearer": {}},
			},
			Middlewares: huma.Middlewares{
				permissions.Authorize(api, permissions.PagesProductsCreate),
				postgres.SetupTransaction(api, pool),
			},
		},
		createProduct(publisher),
	)
}

func createProduct(publisher message.Publisher) func(context.Context, *HumaCreateProductRequest) (*HumaCreateProductResult, error) {
	return func(ctx context.Context, request *HumaCreateProductRequest) (*HumaCreateProductResult, error) {
		errs := v.Validate(request.Schema())
		for _, err := range errs {
			return nil, huma.Error400BadRequest(err.Message())
		}

		tx, err := postgres.GetTxFromCtx(ctx)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		productManager := services.NewProductManager(tx)

		var product models.Product
		if err := copier.Copy(&product, &request.Body); err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		if err := productManager.CreateProduct(ctx, &product); err != nil {
			return nil, huma.Error400BadRequest(err.Error())
		}

		var productDto dtos.ProductDto
		if err := copier.Copy(&productDto, &product); err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		productCreatedEvent := events.ProductCreatedEvent{
			Id:            product.Id,
			Name:          product.Name,
			Description:   product.Description,
			Price:         product.Price.String(),
			StockQuantity: product.StockQuantity,
		}

		payload, err := json.Marshal(productCreatedEvent)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		msg := message.NewMessage(watermill.NewUUID(), payload)
		publisher.Publish(events.ProductCreatedTopicV1, msg)

		result := HumaCreateProductResult{}
		result.Body.Product = productDto

		return &result, nil
	}
}
