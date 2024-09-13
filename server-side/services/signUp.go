package services

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"server-side/model"
	"strings"
	"time"

	"github.com/google/uuid"
)

func SignUpNewUser(w http.ResponseWriter, r *http.Request) {
	user := model.User{
		ID:       uuid.New(),
		Name:     r.FormValue("name"),
		Phone:    r.FormValue("phone"),
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}
	if user.Password == "" {
		HandelError(w, http.StatusBadRequest, "Password is required")
		time.Sleep(2 * time.Second)
		return
	}

	errUserValidation := ValidateUser(user)
	if errUserValidation != nil {
		HandelError(w, http.StatusBadRequest, errUserValidation.Error())
		return
	}

	queryEmail, argsEmail, errEmail := statement.Select("email").
		From("users").
		Where("email = ?", user.Email).
		ToSql()

	if errEmail != nil {
		HandleError(w, http.StatusInternalServerError, "Error generating query")
		return
	}

	var storedEmail string
	err := db.Get(&storedEmail, queryEmail, argsEmail...)
	if err != nil {
		if err == sql.ErrNoRows {
		} else {
			HandelError(w, http.StatusInternalServerError, "Error checking email: "+err.Error())
			return
		}
	} else {
		HandelError(w, http.StatusConflict, "Email already exists")
		return
	}

	file, fileHeader, err := r.FormFile("img")
	if err != nil && err != http.ErrMissingFile {
		HandelError(w, http.StatusBadRequest, "Invalid file")
		time.Sleep(2 * time.Second)
		return
	} else if err == nil {
		defer file.Close()
		imageName, err := SaveImageFile(file, "users", fileHeader.Filename)
		if err != nil {
			HandelError(w, http.StatusInternalServerError, err.Error())
			time.Sleep(2 * time.Second)
			return
		}
		user.Img = &imageName
	}

	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		HandelError(w, http.StatusInternalServerError, "Error hashing password")
		time.Sleep(2 * time.Second)
		return
	}
	user.Password = hashedPassword

	var usersCount int
	countQuery, args, err := statement.Select("count(*)").
		From("users").
		ToSql()
	err = db.Get(&usersCount, countQuery, args...)
	if err != nil {
		HandleError(w, http.StatusInternalServerError, "Error counting users")
		fmt.Printf("count = %d erro:%s", usersCount, err)
		time.Sleep(2 * time.Second)
		return
	}

	query, args, err := statement.
		Insert("users").
		Columns("id", "img", "name", "phone", "email", "password").
		Values(user.ID, user.Img, user.Name, user.Phone, user.Email, user.Password).
		Suffix(fmt.Sprintf("RETURNING %s", strings.Join(user_columns, ", "))).
		ToSql()
	if err != nil {
		HandelError(w, http.StatusInternalServerError, "Error generating query")
		time.Sleep(2 * time.Second)
		return
	}

	if err := db.QueryRowx(query, args...).StructScan(&user); err != nil {
		HandelError(w, http.StatusInternalServerError, "Error creating user: "+err.Error())
		time.Sleep(2 * time.Second)
		return
	}
	if usersCount == 0 {
		assignSystemAdmin(user.ID)
	}

	SendJsonResponse(w, http.StatusCreated, user)
}

func assignSystemAdmin(userId uuid.UUID) {
	query, args, err := statement.Insert("user_roles").
		Values(userId, 1).
		ToSql()
	if err != nil {
		log.Printf("error in query generation: %s", err)
		return
	}
	result, err := db.Exec(query, args...)
	if err != nil {
		log.Printf("eror while excuting query: %s", err)
		return
	}
	log.Printf("%s", result)
}
