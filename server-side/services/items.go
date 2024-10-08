package services

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"server-side/model"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func CreateNewItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userDetails, ok := ctx.Value(userContextKey).(UserDetails)
	if !ok {
		SendCustomeErrorResponse(w, http.StatusUnauthorized, "You are not Authorized to do this action")
		return
	}

	user := userDetails.UserData

	if ValidateIsEmptyOrNil(r.FormValue("vendor_id"), r.FormValue("name"), r.FormValue("price")) {
		HandelError(w, http.StatusBadRequest, "name or price can't be empty!")
		return
	}

	price, err := strconv.ParseFloat(r.FormValue("price"), 64)
	if err != nil {
		log.Println("Error parsing price:", err)
		HandelError(w, http.StatusInternalServerError, "Error parsin price")
		return
	}

	if price <= 0 {
		SendCustomeErrorResponse(w, http.StatusBadRequest, "Price must be greater than zero")
		return
	}

	if !ValidUUID(r.FormValue("vendor_id")) {
		HandelError(w, http.StatusBadRequest, "vendorId is not Valid")
		return
	}

	item := model.Item{
		Id:        uuid.New(),
		Name:      r.FormValue("name"),
		Price:     float64(price),
		Vendor_id: uuid.MustParse(r.FormValue("vendor_id")),
	}

	exist, err := RowExists("vendor_admins", map[string]interface{}{"vendor_id": item.Vendor_id.String(), "user_id": user.ID.String()})
	if err != nil {
		SendErrorResponse(w, err)
		return
	}
	if !exist {
		SendCustomeErrorResponse(w, http.StatusUnauthorized, "You are not Authorized to do this action")
		return
	}

	exist, err = RowExists("vendors", map[string]interface{}{"id": item.Vendor_id.String()})
	if err != nil {
		SendErrorResponse(w, err)
		return
	}
	if !exist {
		SendCustomeErrorResponse(w, http.StatusNotFound, "Vendor not found")
		return
	}

	file, fileHeader, err := r.FormFile("img")
	if err != nil && err != http.ErrMissingFile {
		HandelError(w, http.StatusBadRequest, "Invalid file")
		time.Sleep(2 * time.Second)
		return
	} else if err == nil {
		defer file.Close()
		imageName, err := SaveImageFile(file, "items", fileHeader.Filename)
		if err != nil {
			HandelError(w, http.StatusInternalServerError, err.Error())
			time.Sleep(2 * time.Second)
			return
		}
		item.Img = &imageName
	}

	query, args, err := statement.Insert("items").
		Columns("id, name, price, img, vendor_id").
		Values(item.Id, item.Name, item.Price, item.Img, item.Vendor_id).
		Suffix(fmt.Sprintf("RETURNING %s", strings.Join(item_columns, ", "))).
		ToSql()

	if err := db.QueryRowx(query, args...).StructScan(&item); err != nil {
		HandelError(w, http.StatusInternalServerError, "Error creating item: "+err.Error())
		return
	}

	SendJsonResponse(w, http.StatusCreated, item)
}

func GetItems(w http.ResponseWriter, r *http.Request) {
	var items []model.Item

	meta, err := GetData(r, &items, "items", item_columns)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No Rows found")
			HandelError(w, http.StatusNotFound, "No items found")
			return

		}
		log.Println("Error retrieving items => ", err)
		HandelError(w, http.StatusInternalServerError, "Error retrieving items ")
		return
	}

	result := model.Response{
		Meta: meta,
		Data: items,
	}

	SendJsonResponse(w, http.StatusOK, result)
}

func GetItemById(w http.ResponseWriter, r *http.Request) {
	itemID := r.PathValue("id")

	if itemID == "" {
		SendCustomeErrorResponse(w, http.StatusBadRequest, "There is not id in the path!")
	}

	query, args, err := statement.Select(strings.Join(item_columns, ",")).
		From("items").
		Where("id = ?", itemID).
		ToSql()
	if err != nil {
		SendCustomeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	var item model.Item
	if err := db.Get(&item, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			SendCustomeErrorResponse(w, http.StatusNotFound, "Item not found")
		} else {
			SendCustomeErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	SendJsonResponse(w, http.StatusOK, item)
}

// TODO edit your shitty sql builder !

func UpdateItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userDetails, ok := ctx.Value(userContextKey).(UserDetails)
	if !ok {
		SendCustomeErrorResponse(w, http.StatusUnauthorized, "You are not Authorized to do this action")
		return
	}

	user := userDetails.UserData

	var item model.Item
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		HandelError(w, http.StatusBadRequest, "Not Valid item id")
		return
	}
	query, args, err := statement.Select(item_columns...).
		From("items").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := db.Get(&item, query, args...); err != nil {
		HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}

	hasRole := false
	for _, role := range userDetails.UserRoles {
		if role == "admin" {
			hasRole = true
			break
		}
	}

	if !hasRole {
		exist, err := RowExists("vendor_admins", map[string]interface{}{"vendor_id": item.Vendor_id.String(), "user_id": user.ID.String()})
		if err != nil {
			SendErrorResponse(w, err)
			return
		}
		if !exist {
			SendCustomeErrorResponse(w, http.StatusUnauthorized, "You are not Authorized to do this action")
			return
		}
	}

	if r.FormValue("name") != "" {
		item.Name = r.FormValue("name")
	}

	if r.FormValue("price") != "" {
		price, err := strconv.ParseFloat(r.FormValue("price"), 64)
		if err != nil {
			HandelError(w, http.StatusInternalServerError, "not valide price type")
		}
		if price < 0 {
			HandelError(w, http.StatusBadRequest, "price can't be less than 0")
		}
		item.Price = price
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
		if item.Img != nil {
			oldImg = item.Img
			fmt.Printf("the old image path is -> %s \n", *oldImg)
		}

		imageName, err := SaveImageFile(file, "items", fileHeader.Filename)
		if err != nil {
			HandelError(w, http.StatusInternalServerError, "Error saving image file: "+err.Error())
			return
		}
		item.Img = &imageName
		newImg = &imageName
	}

	if item.Img != nil {
		*item.Img = strings.TrimPrefix(*item.Img, DOMAIN+"/")
	}

	query, args, err = statement.
		Update("items").
		Set("img", item.Img).
		Set("name", item.Name).
		Set("updated_at", time.Now()).
		Set("price", item.Price).
		Where(squirrel.Eq{"id": item.Id}).
		Suffix(fmt.Sprintf("RETURNING %s", strings.Join(item_columns, ", "))).
		ToSql()
	if err != nil {
		DeleteImageFile(*newImg)
		HandelError(w, http.StatusInternalServerError, "Error building query")
		return
	}

	if err := db.QueryRowx(query, args...).StructScan(&item); err != nil {
		DeleteImageFile(*newImg)
		HandelError(w, http.StatusInternalServerError, "Error creating user"+err.Error())
		return
	}

	if oldImg != nil {
		if err := DeleteImageFile(*oldImg); err != nil {
			log.Println(err)
		}
	}

	SendJsonResponse(w, http.StatusOK, item)
}

func DeleteItemById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userDetails, ok := ctx.Value(userContextKey).(UserDetails)
	if !ok {
		SendCustomeErrorResponse(w, http.StatusUnauthorized, "You are not Authorized to do this action")
		return
	}

	user := userDetails.UserData
	hasRole := false
	for _, role := range userDetails.UserRoles {
		if role == "admin" {
			hasRole = true
			break
		}
	}

	id := r.PathValue("id")
	var item model.Item
	err := ReadByID(&item, "items", item_columns, id)
	if err != nil {
		SendErrorResponse(w, err)
		return
	}

	if !hasRole {
		exist, err := RowExists("vendor_admins", map[string]interface{}{"vendor_id": item.Vendor_id.String(), "user_id": user.ID.String()})
		if err != nil {
			SendErrorResponse(w, err)
			return
		}
		if !exist {
			SendCustomeErrorResponse(w, http.StatusUnauthorized, "You are not Authorized to do this action")
			return
		}
	}

	if !ValidUUID(id) {
		HandelError(w, http.StatusBadRequest, "Invalid item Id")
	}

	query, args, err := statement.Delete("items").
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

	log.Printf("Deleted %d items", rowsAffected)

	SendJsonResponse(w, http.StatusAccepted, "Deleted Successfuly")
}

func DeleteAllItems(w http.ResponseWriter, r *http.Request) {
	query, args, err := statement.
		Delete("*").
		From("items").
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

	log.Printf("Deleted %d items", rowsAffected)

	SendJsonResponse(w, http.StatusAccepted, "Deleted Successfuly")
}
