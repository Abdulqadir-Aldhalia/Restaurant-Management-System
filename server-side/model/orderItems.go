package model

import "github.com/google/uuid"

type OrderItems struct {
	Id       uuid.UUID `db:"id" json:"id"`
	Order_id uuid.UUID `db:"order_id" json:"order_id"`
	Item_id  uuid.UUID `db:"item_id" json:"item_id"`
	Quantity int32     `db:"quantity" json:"quantity"`
	Price    float32   `db:"price" json:"price"`
}
