package services

import (
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

func GrantRole(w http.ResponseWriter, r *http.Request) {
	userId := r.FormValue("user_id")
	roleId := r.FormValue("role_id")

	if userId == "" || roleId == "" {
		HandelError(w, http.StatusBadRequest, "userId or roleId can't be empty!")
		return
	}

	if !ValidRole(roleId) {
		HandelError(w, http.StatusBadRequest, "roleId is not a valid role")
		return
	}

	user_id, err := uuid.Parse(userId)
	if err != nil {
		HandelError(w, http.StatusBadRequest, "Invalid UUID format")
		log.Fatal(err)
		return
	}

	role_id, err := strconv.Atoi(roleId)

	query, args, err := statement.
		Insert("user_roles").
		Columns(userRole_columns...).
		Values(user_id, role_id).
		ToSql()
	if err != nil {
		HandelError(w, http.StatusInternalServerError, "Error while creating Query")
		log.Println(err)
		return
	}

	result, err := db.Exec(query, args...)
	if err != nil {
		SendCustomeErrorResponse(w, http.StatusInternalServerError, "Error while excuting query")
		log.Fatal(err)
		return
	}

	log.Printf("Granted Role Successfully the effected columns are: \n %s", result)

	SendJsonResponse(w, http.StatusAccepted, "Granted Role Successfully")
}

func RevokeRole(w http.ResponseWriter, r *http.Request) {
	userId := r.FormValue("user_id")
	roleId := r.FormValue("role_id")

	if userId == "" || roleId == "" {
		HandelError(w, http.StatusBadRequest, "userId or roleId can't be empty!")
		return
	}

	if !ValidRole(roleId) {
		HandelError(w, http.StatusBadRequest, "roleId is not a valid role")
		return
	}

	user_id, err := uuid.Parse(userId)
	if err != nil {
		HandelError(w, http.StatusBadRequest, "Invalid UUID format")
		return
	}

	roleIdInt, err := strconv.Atoi(roleId)
	role_id := int32(roleIdInt)

	query, args, err := statement.
		Delete("user_roles").
		Where("user_id = ?", user_id).
		Where("role_id = ?", role_id).
		ToSql()
	if err != nil {
		HandelError(w, http.StatusInternalServerError, "Error while creating Query")
		return
	}

	result, err := db.Exec(query, args...)
	if err != nil {
		SendCustomeErrorResponse(w, http.StatusInternalServerError, "Error while excuting query")
		return
	}

	log.Printf("Role Deleted Successfully the effected columns are: \n %s", result)

	SendJsonResponse(w, http.StatusAccepted, "Role Deleted Successfully")
}
