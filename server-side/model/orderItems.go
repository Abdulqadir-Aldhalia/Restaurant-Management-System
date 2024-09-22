package model

import "github.com/google/uuid"

type Order_items struct {
	id       uuid.UUID
	order_id uuid.UUID
	item_id  uuid.UUID
	quantity int32
	price    float32
}
