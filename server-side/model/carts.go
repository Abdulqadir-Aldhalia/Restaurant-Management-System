package model

import (
	"time"

	"github.com/google/uuid"
)

type Carts struct {
	User_id    uuid.UUID `db:"id" json:"id"`
	Vendor_id  uuid.UUID `db:"vendor_id" json:"vendor_id"`
	Created_at time.Time `db:"created_at" json:"created_at"`
	Updated_at time.Time `db:"updated_at" json:"updated_at"`
}
