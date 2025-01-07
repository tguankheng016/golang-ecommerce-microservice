package seed

import (
	"context"

	productGrpc "github.com/tguankheng016/go-ecommerce-microservice/internal/services/cart_service/internal/products/grpc_client/protos"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/cart_service/internal/products/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func SeedProducts(ctx context.Context, productGrpcClientService productGrpc.ProductGrpcServiceClient, database *mongo.Database) error {
	res, err := productGrpcClientService.GetAllProducts(ctx, &productGrpc.GetAllProductsRequest{})

	if err != nil {
		return err
	}

	if len(res.Products) == 0 {
		return nil
	}

	productCollection := database.Collection(models.ProductCollectionName)

	for _, grpcProduct := range res.Products {
		newProduct := models.Product{
			Id:          int(grpcProduct.Id),
			Name:        grpcProduct.Name,
			Description: grpcProduct.Description,
			Price:       grpcProduct.Price,
		}

		filter := bson.D{bson.E{Key: "id", Value: newProduct.Id}}

		var dbProduct models.Product
		err = productCollection.FindOne(ctx, filter).Decode(&dbProduct)

		if err != nil && err != mongo.ErrNoDocuments {
			return err
		}

		if err == mongo.ErrNoDocuments {
			_, err := productCollection.InsertOne(ctx, newProduct)
			if err != nil {
				return err
			}
		} else {
			update := bson.D{
				{
					Key: "$set",
					Value: bson.D{
						bson.E{Key: "name", Value: newProduct.Name},
						bson.E{Key: "description", Value: newProduct.Description},
						bson.E{Key: "price", Value: newProduct.Price},
					},
				},
			}

			_, err := productCollection.UpdateOne(ctx, filter, update)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
