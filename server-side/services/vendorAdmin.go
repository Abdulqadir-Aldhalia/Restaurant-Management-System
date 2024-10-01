package services

import (
	"fmt"
	"log"
	"net/http"
	"server-side/model"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func GetAllAdminsForVendor(w http.ResponseWriter, r *http.Request) {
	var vendorAdmins []model.VendorAdmin
	var admins []model.User

	vendorId := r.PathValue("vendor_id")

	vendor_id, err := uuid.Parse(vendorId)
	if err != nil {
		SendCustomeErrorResponse(w, http.StatusBadRequest, "Invalid UUID format for vendor_id")
		log.Println(err)
		return
	}

	query, args, err := statement.Select(strings.Join(vendorAdmins_columns, ", ")).
		From("vendor_admins").
		Where("vendor_id = ?", vendor_id).
		ToSql()
	if err != nil {
		SendCustomeErrorResponse(w, http.StatusInternalServerError, "Error while creating query")
		log.Println(err)
		return
	}

	if err := db.Select(&vendorAdmins, query, args...); err != nil {
		SendCustomeErrorResponse(w, http.StatusInternalServerError, "Error while executing query")
		log.Println(err)
		return
	}

	var userIds []uuid.UUID
	for _, va := range vendorAdmins {
		userIds = append(userIds, va.UserId)
	}

	if len(userIds) == 0 {
		SendCustomeErrorResponse(w, http.StatusNotFound, "There is no admins for the provided vendor !")
		return
	}

	query, args, err = statement.
		Select(strings.Join(user_columns, ", ")).
		From("users").
		Where(squirrel.Eq{"id": userIds}).
		ToSql()
	if err != nil {
		SendCustomeErrorResponse(w, http.StatusInternalServerError, "Error while creating user query")
		log.Println(err)
		return
	}

	if err := db.Select(&admins, query, args...); err != nil {
		SendCustomeErrorResponse(w, http.StatusInternalServerError, "Error while executing user query")
		log.Println(err)
		return
	}

	SendJsonResponse(w, http.StatusOK, admins)
}

func GetAllVendorAdmins(w http.ResponseWriter, r *http.Request) {
	var vendorAdmins []model.VendorAdmin

	query, args, err := statement.Select(strings.Join(vendorAdmins_columns, ", ")).
		From("vendor_admins").
		ToSql()
	if err != nil {
		HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := db.Select(&vendorAdmins, query, args...); err != nil {
		HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}

	SendJsonResponse(w, http.StatusOK, vendorAdmins)
}

func AssignAdminToVendor(w http.ResponseWriter, r *http.Request) {
	userId := r.FormValue("user_id")
	vendorId := r.FormValue("vendor_id")

	if userId == "" || vendorId == "" {
		HandelError(w, http.StatusBadRequest, "userId or vendor_id can't be empty!")
		return
	}

	user_id, err := uuid.Parse(userId)
	if err != nil {
		log.Print(user_id)
		HandelError(w, http.StatusBadRequest, "Invalid UUID for user_id format")
		log.Println(err)
		return
	}

	vendor_id, err := uuid.Parse(vendorId)
	if err != nil {
		log.Print(vendorId)
		HandelError(w, http.StatusBadRequest, "Invalid UUID vendor_id format")
		log.Println(err)
		return
	}

	query, args, err := statement.
		Insert("vendor_admins").
		Columns(vendorAdmins_columns...).
		Values(user_id, vendor_id).
		ToSql()
	if err != nil {
		HandelError(w, http.StatusInternalServerError, "Error while creating Query")
		log.Println(err)
		return
	}

	result, err := db.Exec(query, args...)
	if err != nil {
		SendCustomeErrorResponse(w, http.StatusInternalServerError, "Error while excuting query")
		log.Println(err)
		return
	}

	query = fmt.Sprintf("INSERT INTO user_roles VALUES ('%s', '%d')", userId, 2)
	_, err = db.Exec(query)
	if err != nil {
		SendErrorResponse(w, err)
		return
	}

	log.Printf("Admin Assigned Successfully the effected columns are: \n %s", result)

	SendJsonResponse(w, http.StatusAccepted, "Admin Assigned Successfully")
}

func RevokeAdminFromVendor(w http.ResponseWriter, r *http.Request) {
	userId := r.FormValue("user_id")
	vendorId := r.FormValue("vendor_id")

	if userId == "" || vendorId == "" {
		SendCustomeErrorResponse(w, http.StatusBadRequest, "userId or vendorId can't be empty!")
		return
	}

	user_id, err := uuid.Parse(userId)
	if err != nil {
		SendCustomeErrorResponse(w, http.StatusBadRequest, "Invalid UUID format for user_id")
		log.Println(err)
		return
	}

	vendor_id, err := uuid.Parse(vendorId)
	if err != nil {
		SendCustomeErrorResponse(w, http.StatusBadRequest, "Invalid UUID format for vendor_id")
		log.Println(err)
		return
	}

	query, args, err := statement.
		Delete("vendor_admins").
		Where("user_id = ?", user_id).
		Where("vendor_id = ?", vendor_id).
		ToSql()
	if err != nil {
		SendCustomeErrorResponse(w, http.StatusInternalServerError, "Error while creating query")
		log.Println(err)
		return
	}

	result, err := db.Exec(query, args...)
	if err != nil {
		SendCustomeErrorResponse(w, http.StatusInternalServerError, "Error while executing query")
		log.Println(err)
		return
	}

	log.Printf("Admin role revoked successfully. Affected rows: %d\n", result)

	SendJsonResponse(w, http.StatusAccepted, "Admin role revoked successfully")
}
