package users

import "webinar/graphql/server/internal/entities"

func buildCustomerEntity(dbModel customerModel) entities.Customer {
	return entities.Customer{
		ID:      dbModel.ID,
		UserID:  dbModel.UserID,
		Name:    dbModel.Name,
		Address: dbModel.Address,
	}
}

func buildSessionEntity(dbModel sessionModel) entities.Session {
	return entities.Session{
		UserID: dbModel.UserID,
		Role:   dbModel.Role,
	}
}
