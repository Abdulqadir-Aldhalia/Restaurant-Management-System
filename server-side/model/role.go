package model

type enum_role int

const (
	admin enum_role = iota + 1
	vendor
	customer
)

type Role struct {
	ID   int32  `db="id" json:"id"`
	Name string `db:"name" json:"name"`
}
