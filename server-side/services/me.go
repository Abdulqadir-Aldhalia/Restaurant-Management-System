package services

import (
	"fmt"
	"net/http"
	"server-side/model"
	"strings"
)

func GetUserData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userDetails, ok := ctx.Value(userContextKey).(UserDetails)
	if !ok {
		SendCustomeErrorResponse(w, http.StatusUnauthorized, "You are not authorized to do this action")
		return
	}

	id := userDetails.UserData.ID
	var user model.User
	query, args, err := statement.Select(strings.Join(user_columns, ", ")).
		From("users").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := db.Get(&user, query, args...); err != nil {
		HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}

	SendJsonResponse(w, http.StatusOK, user)
}

func GetUserVendors(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userDetails, ok := ctx.Value(userContextKey).(UserDetails)
	if !ok {
		SendCustomeErrorResponse(w, http.StatusUnauthorized, "You are not authorized to do this action")
		return
	}

	user := userDetails.UserData

	query := fmt.Sprintf("SELECT * FROM vendors WHERE id IN (SELECT vendor_id FROM vendor_admins WHERE user_id = '%s')", user.ID.String())
	var vendors []model.Vendor
	err := db.Select(&vendors, query)
	if err != nil {
		SendErrorResponse(w, err)
		return
	}

	SendJsonResponse(w, http.StatusOK, vendors)
}
