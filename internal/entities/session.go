package entities

import "webinar/graphql/server/internal/graph/model"

type Session struct {
	UserID int64
	Role   model.Role
}
