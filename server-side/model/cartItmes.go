package model

import "github.com/google/uuid"

type Cart_items struct {
	Cart_id  uuid.UUID `db:"cart_id" json:"cart_id"`
	Item_id  uuid.UUID `db:"item_id" json:"item_id"`
	Quantity int       `db:"quantity" json:"quantity"`
}
