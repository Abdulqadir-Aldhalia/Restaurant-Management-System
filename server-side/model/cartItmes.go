package model

import "github.com/google/uuid"

type Cart_items struct {
	cart_id  uuid.UUID `db:"cart_id json:"cart_id"`
	item_id  uuid.UUID `db:"item_id" json:"item_id"`
	quantity int       `db:"quantity" json:"quantity"`
}
