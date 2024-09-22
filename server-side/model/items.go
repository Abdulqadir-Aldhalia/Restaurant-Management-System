package model

import (
	"time"

	"github.com/google/uuid"
)

type Item struct {
	Id        uuid.UUID `db:"id" json:"id"`
	Vendor_id uuid.UUID `db:"vendor_id" json:"vendor_id"`
	Name      string    `db:"name" json:"name"`
	Price     float64   `db:"price" json:"price"`
	Img       *string   `db:"img" json:"img"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
