package users

import "webinar/graphql/server/internal/graph/model"

type customerModel struct {
	ID      int64  `db:"id"`
	UserID  int64  `db:"user_id"`
	Name    string `db:"name"`
	Address string `db:"address"`
}

type sessionModel struct {
	UserID int64      `db:"id"`
	Role   model.Role `db:"role"`
}
