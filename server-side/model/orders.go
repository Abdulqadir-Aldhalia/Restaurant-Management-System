package model

import (
	"time"

	"github.com/google/uuid"
)

type order_status string

const (
	PENDING   order_status = "PENDING"
	PREPARING order_status = "PREPEARING"
	READY     order_status = "READY"
)

type Orders struct {
	Id               uuid.UUID    `db:"id" json:"id"`
	Total_order_cost float32      `db:"total_order_cost" json:"total_order_cost"`
	Customer_id      uuid.UUID    `db:"customer_id" json:"customer_id"`
	Vendor_id        uuid.UUID    `db:"vendor_id" json:"vendor_id"`
	Status           order_status `db:"status" json:"status"`
	Created_at       time.Time    `db:"created_at" json:"created_at"`
	Updated_at       time.Time    `db:"updated_at" json:"updated_at"`
}
