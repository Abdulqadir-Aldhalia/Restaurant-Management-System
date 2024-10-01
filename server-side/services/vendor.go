package services

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"server-side/model"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func CreateNewVendor(w http.ResponseWriter, r *http.Request) {
	vendor := model.Vendor{
		ID:          uuid.New(),
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
	}

	file, fileHeader, err := r.FormFile("img")

	if err != nil && err != http.ErrMissingFile {
		HandelError(w, http.StatusBadRequest, "Invalid file")
		return
	} else if err == nil {
		defer file.Close()
		imageName, err := SaveImageFile(file, "vendors", fileHeader.Filename)
		if err != nil {
			HandelError(w, http.StatusInternalServerError, "Error saving image")
		}

		vendor.Img = &imageName
	}

	query, args, err := statement.
		Insert("vendors").
		Columns("id", "img", "name", "description").
		Values(vendor.ID, vendor.Img, vendor.Name, vendor.Description).
		Suffix(fmt.Sprintf("RETURNING %s", strings.Join(vendor_columns, ", "))).
		ToSql()
	if err != nil {
		HandelError(w, http.StatusInternalServerError, "Error generate query")
		return
	}

	if err := db.QueryRowx(query, args...).StructScan(&vendor); err != nil {
		HandelError(w, http.StatusInternalServerError, "Error creating vendor: "+err.Error())
		return
	}

	SendJsonResponse(w, http.StatusCreated, vendor)
}

func GetVendors(w http.ResponseWriter, r *http.Request) {
	var vendors []model.Vendor

	meta, err := GetData(r, &vendors, "vendors", vendor_columns)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No Rows found")
			HandelError(w, http.StatusNotFound, "No vendor found")
			return

		}
		log.Println("Error retrieving vendros => ", err)
		HandelError(w, http.StatusInternalServerError, "Error retrieving vendors")
		return
	}

	result := model.Response{
		Meta: meta,
		Data: vendors,
	}

	SendJsonResponse(w, http.StatusOK, result)
}

func GetVendorById(w http.ResponseWriter, r *http.Request) {
	vendorID := r.PathValue("id")

	if vendorID == "" {
		SendCustomeErrorResponse(w, http.StatusBadRequest, "There is not id in the path!")
	}

	query, args, err := statement.Select(strings.Join(vendor_columns, ",")).
		From("vendors").
		Where("id = ?", vendorID).
		ToSql()
	if err != nil {
		SendCustomeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	var vendor model.Vendor
	if err := db.Get(&vendor, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			SendCustomeErrorResponse(w, http.StatusNotFound, "Vendor not found")
		} else {
			SendCustomeErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	SendJsonResponse(w, http.StatusOK, vendor)
}

func UpdateVendor(w http.ResponseWriter, r *http.Request) {
	var vendor model.Vendor
	id := r.PathValue("id")
	query, args, err := statement.Select(vendor_columns...).
		From("vendors").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := db.Get(&vendor, query, args...); err != nil {
		HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// update user
	if r.FormValue("name") != "" {
		vendor.Name = r.FormValue("name")
	}

	if r.FormValue("description") != "" {
		vendor.Description = r.FormValue("description")
	}

	var oldImg *string
	var newImg *string
	// Handle image file upload
	file, fileHeader, err := r.FormFile("img")
	if err != nil && err != http.ErrMissingFile {
		HandelError(w, http.StatusBadRequest, "Error retrieving file: "+err.Error())
		return
	} else if err == nil {
		defer file.Close()
		if vendor.Img != nil {
			oldImg = vendor.Img
			fmt.Printf("the old image path is -> %s \n", *oldImg)
		}

		imageName, err := SaveImageFile(file, "vendors", fileHeader.Filename)
		if err != nil {
			HandelError(w, http.StatusInternalServerError, "Error saving image file: "+err.Error())
			return
		}
		vendor.Img = &imageName
		newImg = &imageName
	}

	if vendor.Img != nil {
		*vendor.Img = strings.TrimPrefix(*vendor.Img, DOMAIN+"/")
	}

	query, args, err = statement.
		Update("vendors").
		Set("img", vendor.Img).
		Set("name", vendor.Name).
		Set("updated_at", time.Now()).
		Set("description", vendor.Description).
		Where(squirrel.Eq{"id": vendor.ID}).
		Suffix(fmt.Sprintf("RETURNING %s", strings.Join(vendor_columns, ", "))).
		ToSql()
	if err != nil {
		DeleteImageFile(*newImg)
		HandelError(w, http.StatusInternalServerError, "Error building query")
		return
	}

	if err := db.QueryRowx(query, args...).StructScan(&vendor); err != nil {
		DeleteImageFile(*newImg)
		HandelError(w, http.StatusInternalServerError, "Error creating user"+err.Error())
		return
	}

	if oldImg != nil {
		if err := DeleteImageFile(*oldImg); err != nil {
			log.Println(err)
		}
	}

	SendJsonResponse(w, http.StatusOK, vendor)
}

func DeleteAllVendors(w http.ResponseWriter, r *http.Request) {
	query, args, err := statement.
		Delete("*").
		From("vendors").
		ToSql()
	if err != nil {
		SendCustomeErrorResponse(w, http.StatusInternalServerError, "Error generating delete query")
		return
	}

	result, err := db.Exec(query, args...)
	if err != nil {
		SendCustomeErrorResponse(w, http.StatusInternalServerError, "Error executing delete query")
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		SendCustomeErrorResponse(w, http.StatusInternalServerError, "Error getting rows affected")
		return
	}

	log.Printf("Deleted %d vendors", rowsAffected)

	SendJsonResponse(w, http.StatusAccepted, "Deleted Successfuly")
}

func DeleteVendorById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if !ValidUUID(id) {
		HandelError(w, http.StatusBadRequest, "Invalid User Id")
	}

	query, args, err := statement.Delete("*").
		From("vendors").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		SendCustomeErrorResponse(w, http.StatusInternalServerError, "Error generating delete query")
		return
	}

	result, err := db.Exec(query, args...)
	if err != nil {
		SendCustomeErrorResponse(w, http.StatusInternalServerError, "Error executing delete query")
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		SendCustomeErrorResponse(w, http.StatusInternalServerError, "Error getting rows affected")
		return
	}

	log.Printf("Deleted %d vendors", rowsAffected)

	SendJsonResponse(w, http.StatusAccepted, "Deleted Successfuly")
}
