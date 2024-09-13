package services

import (
	"database/sql"
	"errors"
	"net/http"
	"server-side/model"
	"strings"
)

func GetAllRoles(w http.ResponseWriter, r *http.Request) {
	var roles []model.Role

	query, args, err := statement.Select(strings.Join(role_columns, ",")).
		From("roles").
		ToSql()
	if err != nil {
		HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := db.Select(&roles, query, args...); err != nil {
		HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if roles == nil {
		HandelError(w, http.StatusNotFound, "There is no roles!")
		return
	}

	SendJsonResponse(w, http.StatusOK, roles)
}

func GetRoleById(w http.ResponseWriter, r *http.Request) {
	roleId := r.PathValue("id")

	if roleId == "" {
		HandleError(w, http.StatusBadRequest, "There is not id in the path!")
		return
	}

	query, args, err := statement.Select(strings.Join(role_columns, ",")).
		From("roles").
		Where("id = ?", roleId).
		ToSql()
	if err != nil {
		HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var role model.Role
	if err := db.Get(&role, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			HandleError(w, http.StatusNotFound, "role not found")
		} else {
			HandleError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	SendJsonResponse(w, http.StatusOK, role)
}
