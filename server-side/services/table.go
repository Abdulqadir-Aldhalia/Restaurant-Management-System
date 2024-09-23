package services

// TODO: complete implementing Table
import (
	"fmt"
	"log"
	"net/http"
	"server-side/model"
	"strconv"
	"strings"

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

func CreateNewTable(w http.ResponseWriter, r *http.Request) {
	var table model.Tables

	vendorId, err := uuid.Parse(r.FormValue("vendor_id"))
	if err != nil {
		HandleError(w, http.StatusBadRequest, "Invalide vendor_id")
		return
	}
	fmt.Printf("vendor_id => %s", vendorId)
	customerId, err := uuid.Parse(r.FormValue("customer_id"))
	if err != nil {
		HandleError(w, http.StatusBadRequest, "Invalide customer_id")
		return
	}

	name := r.FormValue("name")

	if ValidateIsEmptyOrNil(name) {
		HandleError(w, http.StatusBadRequest, "Name can't be empty")
		return
	}

	table.Id = uuid.New()
	table.Name = name
	table.Vendor_id = vendorId
	table.Is_available = availableMap[r.FormValue("is_available")]
	table.Customer_id = customerId
	table.Is_needs_service = r.FormValue("is_needs_service") == "true"

	query, args, err := statement.Select("id").
		From("vendors").
		Where("id = ?", table.Vendor_id).
		ToSql()

	result, err := db.Exec(query, args...)
	if err != nil {
		log.Println("Error: vendor query result -> ", result)
		log.Println("Error: vendor query error -> ", err)
		HandleError(w, http.StatusBadRequest, "vendor not found")
		return
	}

	query, args, err = statement.
		Insert("tables").
		Columns(table_columns...).
		Values(table.Id, table.Vendor_id, table.Name, table.Is_available, table.Customer_id, table.Is_needs_service).
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

	q := r.FormValue("query")
	f := r.FormValue("filter")
	s := r.FormValue("sort")
	isAvailable := r.FormValue("is_available")
	needsService := r.FormValue("is_needs_service")

	statementBuilder := statement.Select(strings.Join(table_columns, ", ")).From("tables")

	if q != "" {
		statementBuilder = statementBuilder.Where("name ILIKE ?", "%"+q+"%")
	}

	if isAvailable != "" {
		available, err := strconv.ParseBool(isAvailable)
		if err == nil {
			statementBuilder = statementBuilder.Where("is_available = ?", available)
		}
	}

	if f != "" {
		vendorId, err := uuid.Parse(f)
		if err != nil {
			log.Printf("invalid vendor uuid => %s", f)
		} else {
			query, args, err := statement.Select("id").
				From("vendors").
				Where("id = ?", vendorId).
				ToSql()
			if err != nil {
				log.Println("error while creating query -> ", err)
			}

			result, err := db.Exec(query, args...)
			if err != nil {
				log.Println("Error: vendor query result -> ", result)
				log.Println("Error: vendor query error -> ", err)
			} else {
				statementBuilder = statementBuilder.Where("vendor_id = ?", f)
			}
		}
	}

	if needsService != "" {
		service, err := strconv.ParseBool(needsService)
		if err == nil {
			statementBuilder = statementBuilder.Where("is_needs_service = ?", service)
		}
	}

	if order, exists := sortOptions[strings.ToLower(s)]; exists {
		statementBuilder = statementBuilder.OrderBy("name " + order)
	} else {
		if s == "" {
			log.Println("No Sort option were passed")
		}
		log.Println("Invalid sort type => ", s)
	}

	query, args, err := statementBuilder.ToSql()
	if err != nil {
		log.Println("Error while building query -> ", err)
		HandleError(w, http.StatusInternalServerError, "Error building query")
		return
	}

	if err := db.Select(&tables, query, args...); err != nil {
		log.Println("Error while executing query -> ", err)
		HandleError(w, http.StatusInternalServerError, "Error fetching tables")
		return
	}

	SendJsonResponse(w, http.StatusOK, tables)
}

func GetAllTables(w http.ResponseWriter, r *http.Request) {
	var tables []model.Tables

	query, args, err := statement.Select(strings.Join(table_columns, ", ")).
		From("tables").
		ToSql()
	if err != nil {
		log.Println("Error while building query -> ", err)
		HandleError(w, http.StatusInternalServerError, "Error building query")
		return
	}

	if err := db.Select(&tables, query, args...); err != nil {
		log.Println("Error while excuting query -> ", err)
		HandleError(w, http.StatusInternalServerError, "Error fetching tables")
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
		HandleError(w, http.StatusInternalServerError, "Error building query")
		return
	}

	if err := db.Get(&table, query, args...); err != nil {
		log.Println("Error while excuting query -> ", err)
		HandleError(w, http.StatusNotFound, "Table not found")
		return
	}

	SendJsonResponse(w, http.StatusOK, table)
}

func UpdateTable(w http.ResponseWriter, r *http.Request) {
	var table model.Tables
	id := r.PathValue("id")

	if !ValidUUID(id) {
		HandleError(w, http.StatusBadRequest, "Invalid table ID")
		return
	}

	// Retrieve the existing table
	query, args, err := statement.Select(strings.Join(table_columns, ", ")).
		From("tables").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		log.Println("Error while building query -> ", err)
		HandleError(w, http.StatusInternalServerError, "Error building query")
		return
	}
	if err := db.Get(&table, query, args...); err != nil {
		log.Println("Error while excuting query -> ", err)
		HandleError(w, http.StatusNotFound, "Table not found")
		return
	}

	// Update fields if provided
	if name := r.FormValue("name"); name != "" {
		table.Name = name
	}
	if isAvailable := r.FormValue("is_available"); isAvailable != "" {
		table.Is_available = isAvailable == "true"
	}
	if isNeedsService := r.FormValue("is_needs_service"); isNeedsService != "" {
		table.Is_needs_service = isNeedsService == "true"
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
		HandleError(w, http.StatusInternalServerError, "Error building query: "+err.Error())
		return
	}

	if err := db.QueryRowx(query, args...).StructScan(&table); err != nil {
		log.Println("Error while excuting query -> ", err)
		HandleError(w, http.StatusInternalServerError, "Error updating table: "+err.Error())
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
		HandleError(w, http.StatusInternalServerError, "Error building query")
		return
	}

	if _, err := db.Exec(query, args...); err != nil {
		log.Println("Error while excuting query -> ", err)
		HandleError(w, http.StatusInternalServerError, "Error deleting table")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
