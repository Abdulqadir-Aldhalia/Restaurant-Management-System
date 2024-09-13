package services

import (
	"fmt"
	"log"
	"net/http"
	"server-side/model"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/lib/pq"
)

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	var users []model.User

	query, args, err := statement.Select(strings.Join(user_columns, ", ")).
		From("users").
		ToSql()
	if err != nil {
		HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := db.Select(&users, query, args...); err != nil {
		HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}

	SendJsonResponse(w, http.StatusOK, users)
}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	var user model.User
	id := r.PathValue("id")
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

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	id := r.PathValue("id")
	if !ValidUUID(id) {
		HandelError(w, http.StatusBadRequest, "User Id is not Valid")
		return
	}

	// Fetch existing user
	query, args, err := statement.Select(user_columns...).
		From("users").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		HandelError(w, http.StatusInternalServerError, err.Error())
		time.Sleep(2 * time.Second)
		return
	}
	if err := db.Get(&user, query, args...); err != nil {
		HandelError(w, http.StatusInternalServerError, err.Error())
		time.Sleep(2 * time.Second)
		return
	}

	// Update user fields based on form input
	if r.FormValue("name") != "" {
		user.Name = r.FormValue("name")
	}
	if r.FormValue("phone") != "" {
		user.Phone = r.FormValue("phone")
	}
	if r.FormValue("email") != "" {
		user.Email = r.FormValue("email")
	}
	if r.FormValue("password") != "" {
		log.Printf("password: %s", r.FormValue("password"))
		hashedPassword, err := HashPassword(r.FormValue("password"))
		if err != nil {
			HandelError(w, http.StatusInternalServerError, "Error hashing password")
			return
		}
		user.Password = hashedPassword
	}

	errUserValidation := ValidateUser(user)
	if errUserValidation != nil {
		HandelError(w, http.StatusBadRequest, errUserValidation.Error())
		return
	}

	var oldImg *string
	var newImg *string
	file, fileHeader, err := r.FormFile("img")
	if err != nil && err != http.ErrMissingFile {
		HandelError(w, http.StatusBadRequest, "Error retrieving file: "+err.Error())
		time.Sleep(2 * time.Second)
		return
	} else if err == nil {
		defer file.Close()
		if user.Img != nil {
			oldImg = user.Img
		}
		imageName, err := SaveImageFile(file, "users", fileHeader.Filename)
		if err != nil {
			HandelError(w, http.StatusInternalServerError, "Error saving image file: "+err.Error())
			time.Sleep(2 * time.Second)
			return
		}
		user.Img = &imageName
		newImg = &imageName
	}

	if user.Img != nil && *user.Img != "" {
		*user.Img = strings.TrimPrefix(*user.Img, Domain+"/")
	}

	// Prepare the update query
	query, args, err = statement.
		Update("users").
		Set("img", user.Img).
		Set("name", user.Name).
		Set("email", user.Email).
		Set("phone", user.Phone).
		Set("password", user.Password).
		Set("updated_at", time.Now()).
		Where(squirrel.Eq{"id": user.ID}).
		Suffix(fmt.Sprintf("RETURNING %s", strings.Join(user_columns, ", "))).
		ToSql()
	if err != nil {
		if newImg != nil {
			DeleteImageFile(*newImg)
		}
		HandelError(w, http.StatusInternalServerError, "Error building query")
		time.Sleep(2 * time.Second)
		return
	}

	// Execute the update query and handle unique constraint violation
	if err := db.QueryRowx(query, args...).StructScan(&user); err != nil {
		if newImg != nil {
			DeleteImageFile(*newImg)
		}

		// Check for duplicate key (unique constraint) violation
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				HandelError(w, http.StatusConflict, "Email already exists")
				return
			}
		}

		HandelError(w, http.StatusInternalServerError, "Error updating user: "+err.Error())
		time.Sleep(2 * time.Second)
		return
	}

	// Delete old image if it was updated
	if oldImg != nil && *oldImg != "" {
		if err := DeleteImageFile(*oldImg); err != nil {
			log.Println(err)
			time.Sleep(2 * time.Second)
		}
	}

	// Send updated user as response
	SendJsonResponse(w, http.StatusOK, user)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	query, args, err := statement.Delete("users").
		Where("id = ?", id).
		Suffix("RETURNING img").
		ToSql()
	if err != nil {
		HandelError(w, http.StatusInternalServerError, "Error building query: "+err.Error())
		time.Sleep(2 * time.Second)
		return
	}

	var img *string

	if err := db.QueryRow(query, args...).Scan(&img); err != nil {
		HandelError(w, http.StatusInternalServerError, "Error deleting user: "+err.Error())
		time.Sleep(2 * time.Second)
		return
	}
	// If the user has an image, delete it
	if img != nil {
		if err := DeleteImageFile(*img); err != nil {
			log.Println(err)
			time.Sleep(2 * time.Second)
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
