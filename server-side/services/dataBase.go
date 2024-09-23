package services

import (
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

var (
	db        *sqlx.DB
	statement = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	user_columns = []string{
		"id",
		"name",
		"email",
		"phone",
		"created_at",
		"updated_at",
		fmt.Sprintf("CASE WHEN NULLIF(img, '') IS NOT NULL THEN FORMAT('%s/%%s', img) ELSE NULL END AS img", Domain),
	}

	vendor_columns = []string{
		"id",
		"name",
		"description",
		"created_at",
		"updated_at",
		fmt.Sprintf("CASE WHEN NULLIF(img, '') IS NOT NULL THEN FORMAT('%s/%%s', img) ELSE NULL END AS img", Domain),
	}

	role_columns = []string{
		"id",
		"name",
	}

	userRole_columns = []string{
		"user_id",
		"role_id",
	}

	vendorAdmins_columns = []string{
		"user_id",
		"vendor_id",
	}

	item_columns = []string{
		"id",
		"name",
		"price",
		"vendor_id",
		"img",
		"created_at",
		"updated_at",
	}

	table_columns = []string{
		"id",
		"vendor_id",
		"name",
		"is_available",
		"customer_id",
		"is_needs_service",
	}
)

func SetDB(database *sqlx.DB) {
	db = database
}
