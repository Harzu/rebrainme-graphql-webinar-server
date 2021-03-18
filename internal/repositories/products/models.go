package products

type orderProductModel struct {
	ID      int64  `db:"id"`
	OrderID int64  `db:"order_id"`
	Name    string `db:"name"`
	Price   int64  `db:"price"`
}
