package services

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"server-side/model"
	"strconv"

	"github.com/google/uuid"
)

func CreateCart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userDetails, ok := ctx.Value(userContextKey).(UserDetails)
	if !ok {
		log.Println("User not found in context or not of type model.User")
		SendCustomeErrorResponse(w, http.StatusUnauthorized, "Not allowed to do this action!")
		return
	}
	user := userDetails.UserData

	cartId := map[string]interface{}{
		"id": user.ID,
	}
	exist, err := RowExists("carts", cartId)
	if exist {
		log.Println("User Already has a cart")
		SendErrorResponse(w, ErrConflict)
		return
	}

	tx, err := db.Beginx()
	if err != nil {
		SendErrorResponse(w, err)
		return
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	searchColumns := map[string]interface{}{
		"customer_id": user.ID,
	}

	var tables []model.Tables

	err = ReadByColumns(&tables, "tables", table_columns, searchColumns)
	if err != nil {
		SendErrorResponse(w, err)
		return
	}

	table := tables[0]

	if table.Customer_id != user.ID {
		SendCustomeErrorResponse(w, http.StatusBadRequest, "You need to reserve a table first")
		return
	}

	data := map[string]interface{}{
		"id":        user.ID,
		"vendor_id": table.Vendor_id,
	}

	var userCart model.Carts

	err = Create("carts", data, &userCart, tx)
	if err != nil {
		log.Println("Error came from cart")
		SendErrorResponse(w, err)
		return
	}

	SendJsonResponse(w, http.StatusCreated, userCart)
}

func GetUserCart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userDetails, ok := ctx.Value(userContextKey).(UserDetails)
	if !ok {
		SendCustomeErrorResponse(w, http.StatusUnauthorized, "You are not Authorized to do this action")
		return
	}

	user := userDetails.UserData

	var userCart model.Carts

	err := ReadByID(&userCart, "carts", cart_columns, user.ID.String())
	if err != nil {
		log.Println("user id -> ", user.ID.String())
		SendErrorResponse(w, err)
		return
	}

	var userItems []model.Item

	subquery := "SELECT item_id FROM cart_items WHERE cart_id = '%s'"
	mainQuery := "SELECT * FROM items WHERE id IN (" + subquery + ")"

	query := fmt.Sprintf(mainQuery, user.ID.String())

	err = db.Select(&userItems, query)
	if err != nil {
		SendErrorResponse(w, err)
		return
	}

	result := map[string]interface{}{
		"cart":  userCart,
		"items": userItems,
	}
	SendJsonResponse(w, http.StatusOK, result)
}

func AddItemToCart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userDetails, ok := ctx.Value(userContextKey).(UserDetails)
	if !ok {
		log.Println("User not found in context or not of type model.User")
		SendCustomeErrorResponse(w, http.StatusUnauthorized, "Not allowed to do this action!")
		return
	}

	user := userDetails.UserData

	item_id, err := uuid.Parse(r.FormValue("item_id"))
	if err != nil {
		SendErrorResponse(w, err)
		return
	}

	quantity, err := strconv.Atoi(r.FormValue("quantity"))
	if err != nil {
		SendErrorResponse(w, err)
		return
	}
	if quantity <= 0 {
		SendErrorResponse(w, err)
		return
	}

	var userCart model.Carts

	err = ReadByID(&userCart, "carts", cart_columns, user.ID.String())
	if err != nil {
		if err != sql.ErrNoRows {
			SendErrorResponse(w, err)
			return
		}

		createdData := map[string]interface{}{
			"id": user.ID,
		}

		err = Create("carts", createdData, &userCart, nil)
		if err != nil {
			SendErrorResponse(w, err)
			return
		}

	}

	exist, err := RowExists("cart_items", map[string]interface{}{"cart_id": user.ID.String(), "item_id": item_id.String()})
	if exist {
		SendErrorResponse(w, ErrConflict)
		return
	}

	var item model.Item

	err = ReadByID(&item, "items", item_columns, item_id.String())
	if err != nil {
		SendErrorResponse(w, err)
		return
	}

	if item.Vendor_id != userCart.Vendor_id {
		if userCart.Vendor_id == uuid.Nil {
			updatedData := map[string]interface{}{
				"vendor_id": item.Vendor_id,
			}
			err = UpdateById("carts", user.ID, updatedData, &userCart, nil)
			if err != nil {
				SendErrorResponse(w, err)
				return
			}
		} else {
			log.Println("You can only purches from the same vendor!")
			SendErrorResponse(w, ErrInvalidArgument)
			return

		}
	}

	tx, err := db.Beginx()
	if err != nil {
		SendErrorResponse(w, err)
		return
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	var userCartItem model.Cart_items
	cartItemMap := map[string]interface{}{
		"cart_id":  user.ID,
		"item_id":  item_id,
		"quantity": quantity,
	}

	err = Create("cart_items", cartItemMap, &userCartItem, tx)
	if err != nil {
		SendErrorResponse(w, err)
		return
	}

	var userItem model.Item
	err = ReadByID(&userItem, "items", item_columns, item_id.String())
	if err != nil {
		SendErrorResponse(w, err)
		return
	}

	result := map[string]interface{}{
		"item":        userItem,
		"total_price": userItem.Price * float64(quantity),
	}

	SendJsonResponse(w, http.StatusOK, result)
}

func RemoveItemFromCart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userDetails, ok := ctx.Value(userContextKey).(UserDetails)
	if !ok {
		log.Println("User not found in context or not of type model.User")
		SendCustomeErrorResponse(w, http.StatusUnauthorized, "Not allowed to do this action!")
		return
	}

	item_id, err := uuid.Parse(r.FormValue("item_id"))
	if err != nil {
		SendErrorResponse(w, err)
		return
	}

	err = DeleteByColumns("cart_items", map[string]interface{}{"item_id": item_id, "cart_id": userDetails.UserData.ID.String()}, nil)
	if err != nil {
		SendErrorResponse(w, err)
		return
	}

	SendJsonResponse(w, http.StatusAccepted, "Deleted Successfully")
}

func EmptyTheCart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userDetails, ok := ctx.Value(userContextKey).(UserDetails)
	if !ok {
		log.Println("User not found in context or not of type model.User")
		SendCustomeErrorResponse(w, http.StatusUnauthorized, "Not allowed to do this action!")
		return
	}

	user := userDetails.UserData

	var userCart model.Carts

	err := ReadByID(&userCart, "carts", cart_columns, user.ID.String())
	if err != nil {
		SendErrorResponse(w, err)
		return
	}

	tx, err := db.Beginx()
	if err != nil {
		SendErrorResponse(w, err)
		return
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	searchedColumns := map[string]interface{}{
		"cart_id": user.ID.String(),
	}

	err = DeleteByColumns("cart_items", searchedColumns, tx)
	if err != nil {
		SendErrorResponse(w, err)
		return
	}

	updatedData := map[string]interface{}{
		"vendor_id": nil,
	}

	searchedColumns = map[string]interface{}{
		"id": user.ID,
	}

	err = UpdateByColumns("carts", updatedData, searchedColumns, &userCart, tx)
	if err != nil {
		SendErrorResponse(w, err)
		return
	}

	SendJsonResponse(w, http.StatusAccepted, userCart)
}
