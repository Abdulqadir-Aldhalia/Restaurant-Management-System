package services

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func AdminSignin(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	log.Printf("%s %s", username, password)
	if username == "" || password == "" {
		HandleError(w, http.StatusBadRequest, "Username and password can't be empty!")
		return
	}

	if len(password) < 8 {
		HandelError(w, http.StatusBadRequest, "password is too short")
		return
	}

	if len(username) < 3 || len(username) > 30 {
		HandelError(w, http.StatusBadRequest, "username must be between 3 and 30 in length")
		return
	}

	var userId, storedEmail, storedHashedPassword string

	query, args, err := statement.Select("id, email, password").
		From("users").
		Where("email = ?", username).
		ToSql()
	if err != nil {
		HandleError(w, http.StatusInternalServerError, "Error while creating query")
		log.Println("SQL generation error:", err)
		return
	}

	err = db.QueryRow(query, args...).Scan(&userId, &storedEmail, &storedHashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			HandleError(w, http.StatusBadRequest, "Username or password is incorrect")
		} else {
			HandleError(w, http.StatusInternalServerError, "Error while retrieving user")
			log.Println("Database query error:", err)
		}
		return
	}

	rQuery, rArgs, rErr := statement.Select("user_id, role_id").
		From("user_roles").
		Where("user_id = ? AND role_id = ?", userId, 1).
		ToSql()
	if rErr != nil {
		HandleError(w, http.StatusInternalServerError, "Error while creating role check query")
		log.Println("SQL generation error:", rErr)
		return
	}

	var userUUID uuid.UUID
	var roleId int32

	err = db.QueryRow(rQuery, rArgs...).Scan(&userUUID, &roleId)
	if err != nil {
		if err == sql.ErrNoRows {
			HandleError(w, http.StatusNotFound, "User is not an admin!")
		} else {
			HandleError(w, http.StatusInternalServerError, "Error while retrieving user role")
			log.Println("Database query error:", err)
		}
		return
	}

	if !matchPassword(password, storedHashedPassword) {
		HandleError(w, http.StatusBadRequest, "Username or password is incorrect")
		return
	}

	token, err := GenerateJWT(storedEmail)
	if err != nil {
		HandleError(w, http.StatusInternalServerError, "Error while generating token")
		log.Println("JWT generation error:", err)
		return
	}

	jsonToken := map[string]string{
		"token": token,
	}
	SendJsonResponse(w, http.StatusOK, jsonToken)
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		HandleError(w, http.StatusBadRequest, "Username and password can't be empty!")
		return
	}
	if len(password) < 8 {
		HandelError(w, http.StatusBadRequest, "password is too short")
		return
	}

	if len(username) < 3 || len(username) > 30 {
		HandelError(w, http.StatusBadRequest, "username must be between 3 and 30 in length")
		return
	}

	query, args, err := statement.Select("email, password").
		From("users").
		Where("email = ?", username).
		ToSql()
	if err != nil {
		HandleError(w, http.StatusInternalServerError, "Error while creating query")
		log.Println("SQL generation error:", err)
		return
	}

	var storedEmail, storedHashedPassword string
	err = db.QueryRow(query, args...).Scan(&storedEmail, &storedHashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			HandleError(w, http.StatusBadRequest, "Username or password is incorrect")
		} else {
			HandleError(w, http.StatusInternalServerError, "Error while retrieving user")
			log.Println("Database query error:", err)
		}
		return
	}

	if !matchPassword(password, storedHashedPassword) {
		HandleError(w, http.StatusBadRequest, "Username or password is incorrect")
		return
	}

	token, err := GenerateJWT(storedEmail)
	if err != nil {
		HandleError(w, http.StatusInternalServerError, "Error while generating token")
		log.Println("JWT generation error:", err)
		return
	}

	jsonToken := map[string]string{
		"token": token,
	}
	SendJsonResponse(w, http.StatusOK, jsonToken)
}

func matchPassword(password string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	log.Println(err)
	return err == nil
}
