package model

import "github.com/google/uuid"

type Tables struct {
	Id               uuid.UUID `db:"id" json:"id"`
	Vendor_id        uuid.UUID `db:"vendor_id" json:"vendor_id"`
	Name             string    `db:"name" json:"name"`
	Is_available     bool      `db:"is_available" json:"is_available"`
	Customer_id      uuid.UUID `db:"customer_id" json:"customer_id"`
	Is_needs_service bool      `db:"is_needs_service" json:"is_needs_service`
}
