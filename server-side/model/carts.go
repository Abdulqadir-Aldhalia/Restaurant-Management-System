package model

import (
	"time"

	"github.com/google/uuid"
)

type carts struct {
	user_id     uuid.UUID `db:"id" json:"id"`
	total_price float32   `db:"total_price" json:"total_price"`
	quantity    int       `db:"quantity" json:"total_price"`
	vendor_id   uuid.UUID `db:"vendor_id" json:"vendor_id"`
	created_at  time.Time `db:"created_at" json:"created_at"`
	updated_at  time.Time `db:"updated_at" json:"updated_at"`
}
