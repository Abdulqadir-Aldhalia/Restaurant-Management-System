package model

import (
	"time"

	"github.com/google/uuid"
)

type Vendor struct {
	ID          uuid.UUID `db:"id"         json:"id"`
	Name        string    `db:"name"       json:"name"`
	Description string    `db:"description" json:"description"`
	Img         *string   `db:"img"        json:"img"`
	Created_at  time.Time `db:"created_at" json:"created_at"`
	Updated_at  time.Time `db:"updated_at" json:"updated_at"`
}
