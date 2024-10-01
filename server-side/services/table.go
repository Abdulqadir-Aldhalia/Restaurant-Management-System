package services

// TODO: complete implementing Table
import (
	"fmt"
	"log"
	"net/http"
	"server-side/model"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

var TrueMap = map[string]bool{
	"true": true,
	"TRUE": true,
	"True": true,
	"T":    true,
	"t":    true,
	"1":    true,
}

func CreateNewTable(w http.ResponseWriter, r *http.Request) {
	var table model.Tables

	vendorId, err := uuid.Parse(r.FormValue("vendor_id"))
	if err != nil {
		SendCustomeErrorResponse(w, http.StatusBadRequest, "Invalide vendor_id")
		return
	}
	fmt.Printf("vendor_id => %s", vendorId)

	name := r.FormValue("name")

	if ValidateIsEmptyOrNil(name) {
		SendCustomeErrorResponse(w, http.StatusBadRequest, "Name can't be empty")
		return
	}

	table.Id = uuid.New()
	table.Name = name
	table.Vendor_id = vendorId
	table.Is_available = TrueMap[r.FormValue("is_available")]
	table.Is_needs_service = TrueMap[r.FormValue("is_needs_service")]

	query, args, err := statement.Select("id").
		From("vendors").
		Where("id = ?", table.Vendor_id).
		ToSql()

	result, err := db.Exec(query, args...)
	if err != nil {
		log.Println("Error: vendor query result -> ", result)
		log.Println("Error: vendor query error -> ", err)
		SendCustomeErrorResponse(w, http.StatusBadRequest, "vendor not found")
		return
	}

	query, args, err = statement.
		Insert("tables").
		Columns(table_columns...).
		Values(table.Id, table.Vendor_id, table.Name, table.Is_available, nil, table.Is_needs_service).
		Suffix(fmt.Sprintf("RETURNING %s", strings.Join(table_columns, ", "))).
		ToSql()
	if err != nil {
		log.Printf("Error in creating query with qury = %s and args = %s", query, args)
		HandelError(w, http.StatusInternalServerError, "Error building query")
		return
	}

	if err := db.QueryRowx(query, args...).StructScan(&table); err != nil {
		log.Printf("Error in excuting query with query = %s and args = %s", query, args)
		log.Printf("Error Info -> %s", err)
		HandelError(w, http.StatusInternalServerError, "Error excuting query: ")
		return
	}

	SendJsonResponse(w, http.StatusCreated, table)
}

func GetTables(w http.ResponseWriter, r *http.Request) {
	var tables []model.Tables

	meta, err := GetData(r, &tables, "tables", table_columns)
	if err != nil {
		SendErrorResponse(w, err)
		return
	}

	result := model.Response{
		Meta: meta,
		Data: tables,
	}

	SendJsonResponse(w, http.StatusOK, result)
}

func GetAllTables(w http.ResponseWriter, r *http.Request) {
	var tables []model.Tables

	query, args, err := statement.Select(strings.Join(table_columns, ", ")).
		From("tables").
		ToSql()
	if err != nil {
		log.Println("Error while building query -> ", err)
		SendCustomeErrorResponse(w, http.StatusInternalServerError, "Error building query")
		return
	}

	if err := db.Select(&tables, query, args...); err != nil {
		log.Println("Error while excuting query -> ", err)
		SendCustomeErrorResponse(w, http.StatusInternalServerError, "Error fetching tables")
		return
	}

	SendJsonResponse(w, http.StatusOK, tables)
}

func GetTableById(w http.ResponseWriter, r *http.Request) {
	var table model.Tables
	id := r.PathValue("id")

	query, args, err := statement.Select(strings.Join(table_columns, ", ")).
		From("tables").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		log.Println("Error while building query -> ", err)
		SendCustomeErrorResponse(w, http.StatusInternalServerError, "Error building query")
		return
	}

	if err := db.Get(&table, query, args...); err != nil {
		log.Println("Error while excuting query -> ", err)
		SendCustomeErrorResponse(w, http.StatusNotFound, "Table not found")
		return
	}

	SendJsonResponse(w, http.StatusOK, table)
}

func ReserveTable(w http.ResponseWriter, r *http.Request) {
	var table model.Tables
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		SendCustomeErrorResponse(w, http.StatusBadRequest, "Invalid table ID")
		return
	}

	customer_id, err := uuid.Parse(r.FormValue("customer_id"))
	if err != nil {
		SendCustomeErrorResponse(w, http.StatusBadRequest, "Invalid customer_id")
		return
	}

	searchedColumns := map[string]interface{}{
		"customer_id": customer_id.String(),
	}

	var userTable []model.Tables

	err = ReadByColumns(&userTable, "tables", table_columns, searchedColumns)
	if err != nil {
		SendErrorResponse(w, err)
		return
	}

	if !userTable[0].Is_available {
		SendErrorResponse(w, ErrConflict)
		return
	}

	query, args, err := statement.Select(strings.Join(table_columns, ", ")).
		From("tables").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		log.Println("Error while building query -> ", err)
		SendCustomeErrorResponse(w, http.StatusInternalServerError, "Error building query")
		return
	}
	if err := db.Get(&table, query, args...); err != nil {
		log.Println("Error while excuting query -> ", err)
		SendCustomeErrorResponse(w, http.StatusNotFound, "Table not found")
		return
	}

	if !table.Is_available {
		log.Println("Table is not avaliable")
		HandelError(w, http.StatusConflict, "Table is not available")
		return
	}

	query, args, err = statement.Select(strings.Join(user_columns, ", ")).
		From("users").
		Where("id = ?", customer_id).
		ToSql()
	if err != nil {
		log.Println("Error building user query -> ", err)
		HandelError(w, http.StatusInternalServerError, "error building query")
		return
	}

	table.Customer_id = customer_id
	table.Is_available = false
	table.Is_needs_service = false

	if _, err := db.Exec(query, args...); err != nil {
		HandelError(w, http.StatusNotFound, "User not found")
		return
	}

	query, args, err = statement.
		Update("tables").
		Set("is_available", table.Is_available).
		Set("customer_id", table.Customer_id).
		Where(squirrel.Eq{"id": table.Id}).
		Suffix(fmt.Sprintf("RETURNING %s", strings.Join(table_columns, ", "))).
		ToSql()
	if err != nil {
		log.Println("Error while building query -> ", err)
		SendCustomeErrorResponse(w, http.StatusInternalServerError, "Error building query: ")
		return
	}

	if err := db.QueryRowx(query, args...).StructScan(&table); err != nil {
		log.Println("Error while excuting query -> ", err)
		SendCustomeErrorResponse(w, http.StatusInternalServerError, "Error updating table: ")
		return
	}

	SendJsonResponse(w, http.StatusOK, table)
}

func UpdateTable(w http.ResponseWriter, r *http.Request) {
	var table model.Tables
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		SendCustomeErrorResponse(w, http.StatusBadRequest, "Invalid table ID")
		return
	}

	customer_id := r.FormValue("customer_id")
	if customer_id != "" {
		customer_id, err := uuid.Parse(customer_id)
		if err != nil {
			SendCustomeErrorResponse(w, http.StatusBadRequest, "Invalid customer_id")
			return
		}

		query, args, err := statement.Select(strings.Join(user_columns, ", ")).
			From("users").
			Where("id = ?", customer_id).
			ToSql()
		if err != nil {
			log.Println("Error building user query -> ", err)
			HandelError(w, http.StatusInternalServerError, "error building query")
			return
		}
		_, err = db.Exec(query, args...)
		if err != nil {
			SendErrorResponse(w, err)
		}
	}

	// Retrieve the existing table
	query, args, err := statement.Select(strings.Join(table_columns, ", ")).
		From("tables").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		log.Println("Error while building query -> ", err)
		SendCustomeErrorResponse(w, http.StatusInternalServerError, "Error building query")
		return
	}
	if err := db.Get(&table, query, args...); err != nil {
		log.Println("Error while excuting query -> ", err)
		SendCustomeErrorResponse(w, http.StatusNotFound, "Table not found")
		return
	}

	// Update fields if provided
	if name := r.FormValue("name"); name != "" {
		table.Name = name
	}
	if isAvailable := r.FormValue("is_available"); isAvailable != "" {
		table.Is_available = TrueMap[r.FormValue("is_available")]
	}
	if isNeedsService := r.FormValue("is_needs_service"); isNeedsService != "" {
		table.Is_needs_service = TrueMap[isNeedsService]
	}
	if customer_id != "" {
		table.Customer_id = uuid.MustParse(customer_id)
	}

	// Update the table in the database
	query, args, err = statement.
		Update("tables").
		Set("name", table.Name).
		Set("is_available", table.Is_available).
		Set("is_needs_service", table.Is_needs_service).
		Where(squirrel.Eq{"id": table.Id}).
		Suffix(fmt.Sprintf("RETURNING %s", strings.Join(table_columns, ", "))).
		ToSql()
	if err != nil {
		log.Println("Error while building query -> ", err)
		SendCustomeErrorResponse(w, http.StatusInternalServerError, "Error building query: "+err.Error())
		return
	}

	if err := db.QueryRowx(query, args...).StructScan(&table); err != nil {
		log.Println("Error while excuting query -> ", err)
		SendCustomeErrorResponse(w, http.StatusInternalServerError, "Error updating table: "+err.Error())
		return
	}

	SendJsonResponse(w, http.StatusOK, table)
}

func DeleteTable(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	query, args, err := statement.Delete("tables").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		log.Println("Error while building query -> ", err)
		SendCustomeErrorResponse(w, http.StatusInternalServerError, "Error building query")
		return
	}

	if _, err := db.Exec(query, args...); err != nil {
		log.Println("Error while excuting query -> ", err)
		SendCustomeErrorResponse(w, http.StatusInternalServerError, "Error deleting table")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func EmptyTheTable(w http.ResponseWriter, r *http.Request) {
	var table model.Tables
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		SendCustomeErrorResponse(w, http.StatusBadRequest, "Invalid table ID")
		return
	}

	err = ReadByID(&table, "tables", table_columns, id.String())
	if err != nil {
		SendErrorResponse(w, err)
		return
	}

	updatedData := map[string]interface{}{
		"customer_id":      nil,
		"is_available":     true,
		"is_needs_service": false,
	}

	err = UpdateById("tables", id, updatedData, &table, nil)
	if err != nil {
		SendErrorResponse(w, err)
		return
	}

	SendJsonResponse(w, http.StatusAccepted, table)
}
