package model

import "github.com/google/uuid"

type VendorAdmin struct {
	UserId   uuid.UUID `db:"user_id" json:"user_id"`
	VendorId uuid.UUID `db:"vendor_id" json:"vendor_id"`
}
