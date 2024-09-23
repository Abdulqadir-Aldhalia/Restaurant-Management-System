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
	if ValidateIsEmptyOrNil(r.FormValue("vendor_id"), r.FormValue("name"), r.FormValue("price")) {
		HandelError(w, http.StatusBadGateway, "name or price can't be empty!")
		return
	}

	price, err := strconv.ParseFloat(r.FormValue("price"), 64)
	if err != nil {
		log.Println("Error parsing price:", err)
		HandelError(w, http.StatusInternalServerError, "Error parsin price")
		return
	}

	if price <= 0 {
		HandleError(w, http.StatusBadRequest, "Price must be greater than zero")
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

	query, args, err := statement.Select("id").
		From("vendors").
		Where("id = ?", item.Vendor_id).
		ToSql()

	var id string
	err = db.Get(&id, query, args...)
	fmt.Printf("query = %s and args = %s", query, args)
	if err != nil {
		log.Printf("vendor is not exist %s", err)
		HandelError(w, http.StatusInternalServerError, "error while excuting qeury !")
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

	query, args, err = statement.Insert("items").
		Columns("id, name, price, img, vendor_id").
		Values(item.Id, item.Name, item.Price, item.Img, item.Vendor_id).
		ToSql()

	result, err := db.Exec(query, args...)
	if err != nil {
		log.Printf("error while excuting query, %s", err)
		HandelError(w, http.StatusInternalServerError, "error while excuting query")
		return
	}
	log.Printf("query result: %s", result)

	SendJsonResponse(w, http.StatusCreated, item)
}

func GetAllItems(w http.ResponseWriter, r *http.Request) {
	var items []model.Item

	query, args, err := statement.Select(strings.Join(item_columns, ",")).
		From("items").
		ToSql()
	if err != nil {
		HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := db.Select(&items, query, args...); err != nil {
		HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if items == nil {
		HandelError(w, http.StatusNotFound, "There is no vendors!")
		return
	}

	SendJsonResponse(w, http.StatusOK, items)
}

func GetItemById(w http.ResponseWriter, r *http.Request) {
	itemID := r.PathValue("id")

	if itemID == "" {
		HandleError(w, http.StatusBadRequest, "There is not id in the path!")
	}

	query, args, err := statement.Select(strings.Join(item_columns, ",")).
		From("items").
		Where("id = ?", itemID).
		ToSql()
	if err != nil {
		HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var item model.Item
	if err := db.Get(&item, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			HandleError(w, http.StatusNotFound, "Item not found")
		} else {
			HandleError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	SendJsonResponse(w, http.StatusOK, item)
}

// TODO edit your shitty sql builder !

func UpdateItem(w http.ResponseWriter, r *http.Request) {
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
	id := r.PathValue("id")

	if !ValidUUID(id) {
		HandelError(w, http.StatusBadRequest, "Invalid User Id")
	}

	query, args, err := statement.Delete("*").
		From("items").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		HandleError(w, http.StatusInternalServerError, "Error generating delete query")
		return
	}

	result, err := db.Exec(query, args...)
	if err != nil {
		HandleError(w, http.StatusInternalServerError, "Error executing delete query")
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		HandleError(w, http.StatusInternalServerError, "Error getting rows affected")
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
		HandleError(w, http.StatusInternalServerError, "Error generating delete query")
		return
	}

	result, err := db.Exec(query, args...)
	if err != nil {
		HandleError(w, http.StatusInternalServerError, "Error executing delete query")
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		HandleError(w, http.StatusInternalServerError, "Error getting rows affected")
		return
	}

	log.Printf("Deleted %d items", rowsAffected)

	SendJsonResponse(w, http.StatusAccepted, "Deleted Successfuly")
}
