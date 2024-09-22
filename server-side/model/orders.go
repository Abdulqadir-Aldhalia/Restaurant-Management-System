package model

import "github.com/google/uuid"

type order_status int

const (
	PENDING order_status = iota + 1
	PREPARING
	READY
)

type Orders struct {
	id               uuid.UUID    `db:"id" json:"id"`
	total_order_cost float32      `db:"total_order_cost" json:"total_order_cost"`
	customer_id      uuid.UUID    `db:"customer_id" json:"customer_id"`
	vendor_id        uuid.UUID    `db:"vendor_id" json:"vendor_id"`
	status           order_status `db:"status" json:"status"`
	created_at       uuid.Time    `db:"created_at" json:"created_at"`
	updated_at       uuid.Time    `db:"updated_at" json:"updated_at"`
}
