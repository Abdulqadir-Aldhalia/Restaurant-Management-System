package services

// TODO: complete implementing Table
import (
	"fmt"
	"net/http"
	"server-side/model"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

var availableMap = map[string]bool{
	"true": true,
	"TRUE": true,
	"True": true,
	"T":    true,
	"t":    true,
	"1":    true,
}

func CreateTable(w http.ResponseWriter, r *http.Request) {
	var table model.Tables

	table.Id = uuid.New()
	table.Vendor_id = uuid.MustParse(r.FormValue("vendor_id"))
	table.Name = r.FormValue("name")
	table.Is_available = availableMap[r.FormValue("is_available")]
	table.Customer_id = uuid.MustParse(r.FormValue("customer_id"))
	table.Is_needs_service = r.FormValue("is_needs_service") == "true"

	query, args, err := statement.
		Insert("tables").
		Columns(table_columns...).
		Values(table.Id, table.Vendor_id, table.Name, table.Is_available, table.Customer_id, table.Is_needs_service).
		Suffix(fmt.Sprintf("RETURNING %s", strings.Join(table_columns, ", "))).
		ToSql()
	if err != nil {
		HandelError(w, http.StatusInternalServerError, "Error building query")
		return
	}

	if err := db.QueryRowx(query, args...).StructScan(&table); err != nil {
		HandelError(w, http.StatusInternalServerError, "Error creating table: "+err.Error())
		return
	}

	SendJsonResponse(w, http.StatusCreated, table)
}

func GetAllTables(w http.ResponseWriter, r *http.Request) {
	var tables []model.Tables

	query, args, err := statement.Select(strings.Join(table_columns, ", ")).
		From("tables").
		ToSql()
	if err != nil {
		HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := db.Select(&tables, query, args...); err != nil {
		HandelError(w, http.StatusInternalServerError, err.Error())
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
		HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := db.Get(&table, query, args...); err != nil {
		HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}

	SendJsonResponse(w, http.StatusOK, table)
}

func UpdateTable(w http.ResponseWriter, r *http.Request) {
	var table model.Tables
	id := r.PathValue("id")
	if !ValidUUID(id) {
		HandelError(w, http.StatusBadRequest, "Table Id is not Valid")
		return
	}

	query, args, err := statement.Select(table_columns...).
		From("tables").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		HandelError(w, http.StatusInternalServerError, err.Error())
		time.Sleep(2 * time.Second)
		return
	}
	if err := db.Get(&table, query, args...); err != nil {
		HandelError(w, http.StatusInternalServerError, err.Error())
		time.Sleep(2 * time.Second)
		return
	}

	if r.FormValue("name") != "" {
		table.Name = r.FormValue("name")
	}
	if r.FormValue("is_available") != "" {
		table.Is_available = r.FormValue("is_available") == "true"
	}
	if r.FormValue("is_needs_service") != "" {
		table.Is_needs_service = r.FormValue("is_needs_service") == "true"
	}

	query, args, err = statement.
		Update("tables").
		Set("name", table.Name).
		Set("is_available", table.Is_available).
		Set("is_needs_service", table.Is_needs_service).
		Where(squirrel.Eq{"id": table.Id}).
		Suffix(fmt.Sprintf("RETURNING %s", strings.Join(table_columns, ", "))).
		ToSql()
	if err != nil {
		HandelError(w, http.StatusInternalServerError, "Error building query")
		time.Sleep(2 * time.Second)
		return
	}

	if err := db.QueryRowx(query, args...).StructScan(&table); err != nil {
		HandelError(w, http.StatusInternalServerError, "Error updating table: "+err.Error())
		time.Sleep(2 * time.Second)
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
		HandelError(w, http.StatusInternalServerError, "Error building query: "+err.Error())
		time.Sleep(2 * time.Second)
		return
	}

	if _, err := db.Exec(query, args...); err != nil {
		HandelError(w, http.StatusInternalServerError, "Error deleting table: "+err.Error())
		time.Sleep(2 * time.Second)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
