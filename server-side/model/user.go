package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `db:"id"         json:"id"`
	Name       string    `db:"name"       json:"name"`
	Email      string    `db:"email"      json:"email"`
	Phone      string    `db:"phone"      json:"phone"`
	Img        *string   `db:"img"        json:"img"`
	Password   string    `db:"password"   json:"-"`
	Created_at time.Time `db:"created_at" json:"created_at"`
	Updated_at time.Time `db:"updated_at" json:"updated_at"`
}
