package products

import "webinar/graphql/server/internal/entities"

func buildProductEntity(dbModel orderProductModel) entities.Product {
	return entities.Product{
		ID:    dbModel.ID,
		Name:  dbModel.Name,
		Price: dbModel.Price,
	}
}
