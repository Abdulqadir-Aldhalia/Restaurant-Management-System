package model

import "github.com/google/uuid"

type UserRole struct {
	UserId uuid.UUID `db="user_id" json:"user_id"`
	RoleId int32     `db="role_id" json:"role_id"`
}
