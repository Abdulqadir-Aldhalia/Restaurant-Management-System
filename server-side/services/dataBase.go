package services

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

var (
	db        *sqlx.DB
	statement = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
)

func SetDB(database *sqlx.DB) { db = database }

var (
	user_columns = []string{
		"id", "name", "email", "phone", "created_at", "updated_at",
		fmt.Sprintf("CASE WHEN NULLIF(img, '') IS NOT NULL THEN FORMAT('%s/%%s', img) ELSE NULL END AS img", Domain),
	}

	vendor_columns = []string{
		"id", "name", "description", "created_at", "updated_at",
		fmt.Sprintf("CASE WHEN NULLIF(img, '') IS NOT NULL THEN FORMAT('%s/%%s', img) ELSE NULL END AS img", Domain),
	}

	role_columns = []string{"id", "name"}

	userRole_columns = []string{"user_id", "role_id"}

	vendorAdmins_columns = []string{"user_id", "vendor_id"}

	item_columns = []string{"id", "name", "price", "vendor_id", "img", "created_at", "updated_at"}

	table_columns = []string{"id", "vendor_id", "name", "is_available", "customer_id", "is_needs_service"}

	cart_columns = []string{"id", "vendor_id", "created_at", "updated_at"}

	cartItems_columns = []string{"cart_id", "item_id", "quantity"}

	order_columns = []string{"id", "total_order_cost", "customer_id", "vendor_id", "status", "created_at", "updated_at"}

	orderItems_columns = []string{"id", "order_id", "items_id", "quantity", "price"}
)

type MetaData struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	TotalRows  int `json:"total_rows"`
	TotalPages int `json:"total_pages"`
	From       int `json:"from"`
	To         int `json:"to"`
}

var stringColumns = map[string]bool{
	"name":        true,
	"description": true,
	"email":       true,
	"title":       true,
}

func GetData(r *http.Request, dest interface{}, tableName string, columns []string) (MetaData, error) {
	page, err := strconv.Atoi(r.FormValue("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	perPage, err := strconv.Atoi(r.FormValue("per_page"))
	if err != nil || perPage <= 0 {
		perPage = 10
	}

	meta := MetaData{
		Page:    page,
		PerPage: perPage,
	}

	sb := squirrel.Select(columns...).PlaceholderFormat(squirrel.Dollar).From(tableName)
	countSb := squirrel.Select("COUNT(*)").PlaceholderFormat(squirrel.Dollar).From(tableName)

	filters := r.FormValue("filters")
	query := r.FormValue("query")
	sort := r.FormValue("sort")

	if filters != "" {
		filterMap := parseFilters(filters)
		sb = applyFilters(sb, filterMap)
		countSb = applyFilters(countSb, filterMap)
	}

	if query != "" {
		orConditions := squirrel.Or{}
		for _, col := range columns {
			if isStringColumn(col) {
				orConditions = append(orConditions, squirrel.ILike{col: "%" + query + "%"})
			}
		}
		sb = sb.Where(orConditions)
		countSb = countSb.Where(orConditions)
	}

	if sort != "" {
		sortColumn, sortOrder := parseSortOption(sort)
		sb = sb.OrderBy(fmt.Sprintf("%s %s", sortColumn, sortOrder))
	}

	countSQL, countArgs, err := countSb.ToSql()
	if err != nil {
		return MetaData{}, err
	}

	var totalRows int
	if err := db.QueryRow(countSQL, countArgs...).Scan(&totalRows); err != nil {
		return MetaData{}, err
	}

	totalPages := int(math.Ceil(float64(totalRows) / float64(meta.PerPage)))
	meta.TotalRows = totalRows
	meta.TotalPages = totalPages

	offset := (meta.Page - 1) * meta.PerPage
	sb = sb.Limit(uint64(meta.PerPage)).Offset(uint64(offset))

	meta.From = offset + 1
	meta.To = offset + meta.PerPage
	if meta.To > totalRows {
		meta.To = totalRows
	}

	querySQL, args, err := sb.ToSql()
	if err != nil {
		return MetaData{}, err
	}

	log.Printf("Query: %s, Args: %v", querySQL, args)

	if err := db.Select(dest, querySQL, args...); err != nil {
		return MetaData{}, err
	}

	return meta, nil
}

func Create(tableName string, data map[string]interface{}, dest interface{}, tx *sqlx.Tx) error {
	sb := squirrel.Insert(tableName).
		SetMap(data).
		PlaceholderFormat(squirrel.Dollar).
		Suffix("RETURNING *")

	query, args, err := sb.ToSql()
	if err != nil {
		return err
	}

	log.Printf("Generated Query: %s", query)
	log.Printf("Arguments: %v", args)

	if tx != nil {
		return tx.QueryRowx(query, args...).StructScan(dest)
	}

	return db.QueryRowx(query, args...).StructScan(dest)
}

func ReadByID(dest interface{}, tableName string, columns []string, id string) error {
	sb := squirrel.Select(columns...).PlaceholderFormat(squirrel.Dollar).From(tableName).Where("id = ?", id)

	query, args, err := sb.ToSql()
	log.Println("query -> ", query)
	log.Println("args -> ", args)
	if err != nil {
		return err
	}

	return db.Get(dest, query, args...)
}

func ReadByColumns(dest interface{}, tableName string, columns []string, filters map[string]interface{}) error {
	sb := squirrel.Select(columns...).PlaceholderFormat(squirrel.Dollar).From(tableName)

	for column, value := range filters {
		sb = sb.Where(fmt.Sprintf("%s = ?", column), value)
		log.Printf("%s = %s", column, value)
	}

	query, args, err := sb.ToSql()
	if err != nil {
		return err
	}
	log.Println("query -> ", query)
	log.Println("args -> ", args)

	return db.Select(dest, query, args...)
}

func RowExists(tableName string, filters map[string]interface{}) (bool, error) {
	sb := squirrel.Select("1").PlaceholderFormat(squirrel.Dollar).From(tableName).Limit(1)

	for column, value := range filters {
		sb = sb.Where(fmt.Sprintf("%s = ?", column), value)
	}

	query, args, err := sb.ToSql()

	log.Println("query -> ", query)
	log.Println("args -> ", args)
	if err != nil {
		return false, err
	}

	var exists int
	err = db.Get(&exists, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return exists > 0, nil
}

func UpdateById(tableName string, id interface{}, data map[string]interface{}, dest interface{}, tx *sqlx.Tx) error {
	filters := map[string]interface{}{
		"id": id, // Assuming the ID column is named "id"
	}

	return UpdateByColumns(tableName, data, filters, dest, tx)
}

func UpdateByColumns(tableName string, data map[string]interface{}, filters map[string]interface{}, dest interface{}, tx *sqlx.Tx) error {
	sb := squirrel.Update(tableName).
		SetMap(data).
		PlaceholderFormat(squirrel.Dollar).
		Suffix("RETURNING *") // Returning the updated row

	for column, value := range filters {
		sb = sb.Where(squirrel.Eq{column: value})
		log.Printf("Filter: %s = %v", column, value)
	}

	query, args, err := sb.ToSql()
	if err != nil {
		log.Println("Error creating query:", err)
		return err
	}

	log.Println("Generated Update Query:", query)
	log.Println("Arguments:", args)

	if tx != nil {
		err = tx.QueryRowx(query, args...).StructScan(dest)
		if err != nil {
			return err
		}
	} else {
		err = db.QueryRowx(query, args...).StructScan(dest)
		if err != nil {
			return err
		}
	}

	return nil
}

func DeleteById(tableName string, id string) error {
	sb := squirrel.Delete(tableName).Where("id = ?", id)

	query, args, err := sb.ToSql()
	if err != nil {
		return err
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

func DeleteByColumns(tableName string, filters map[string]interface{}, tx *sqlx.Tx) error {
	tableName = strings.TrimSpace(tableName)

	sb := squirrel.Delete(tableName).PlaceholderFormat(squirrel.Dollar)

	for column, value := range filters {
		sb = sb.Where(squirrel.Eq{column: value})
		log.Printf("Filter: %s = %v", column, value)
	}

	query, args, err := sb.ToSql()
	if err != nil {
		log.Println("Error creating query:", err)
		return err
	}

	log.Println("Generated Delete Query ->", query)
	log.Println("Arguments ->", args)

	if tx != nil {
		log.Println("Executing within transaction")
		_, err := tx.Exec(query, args...)
		if err != nil {
			log.Println("Error executing query within transaction:", err)
		}
		return err
	}

	log.Println("Executing without transaction")
	_, err = db.Exec(query, args...)
	if err != nil {
		log.Println("Error executing query:", err)
	}
	return err
}

func parseFilters(filters string) map[string]string {
	filterMap := make(map[string]string)
	filterPairs := strings.Split(filters, ",")

	for _, pair := range filterPairs {
		parts := strings.Split(pair, ":")
		if len(parts) == 2 {
			filterMap[parts[0]] = parts[1]
		}
	}
	return filterMap
}

func applyFilters(statementBuilder squirrel.SelectBuilder, filters map[string]string) squirrel.SelectBuilder {
	for column, value := range filters {
		switch column {
		case "-created_at":
			t, err := time.Parse("2006-01-02", value)
			if err == nil {
				statementBuilder = statementBuilder.Where("DATE(created_at) < ?", t)
			} else {
				log.Printf("Invalid date filter: %s", value)
			}
		case "created_at":
			t, err := time.Parse("2006-01-02", value)
			if err == nil {
				statementBuilder = statementBuilder.Where("DATE(created_at) > ?", t)
			} else {
				log.Printf("Invalid date filter: %s", value)
			}
		default:
			statementBuilder = statementBuilder.Where(fmt.Sprintf("%s = ?", column), value)
		}
	}
	return statementBuilder
}

func isStringColumn(column string) bool {
	return stringColumns[column]
}

func parseSortOption(sort string) (string, string) {
	if strings.HasPrefix(sort, "-") {
		return strings.TrimPrefix(sort, "-"), "DESC"
	}
	return sort, "ASC"
}
