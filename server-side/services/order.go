package services

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"server-side/model"
	"time"

	"github.com/google/uuid"
)

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userDetails, ok := ctx.Value(userContextKey).(UserDetails)
	if !ok {
		SendCustomeErrorResponse(w, http.StatusUnauthorized, "You are not authorized to do this action")
		return
	}
	user := userDetails.UserData

	var userCart model.Carts
	err := ReadByID(&userCart, "carts", cart_columns, user.ID.String())
	if err != nil {
		SendErrorResponse(w, err)
		return
	}

	var cartItems []model.Cart_items
	err = ReadByColumns(&cartItems, "cart_items", cartItems_columns, map[string]interface{}{"cart_id": userCart.User_id.String()})
	if err != nil {
		SendErrorResponse(w, err)
		return
	}

	if len(cartItems) == 0 {
		SendCustomeErrorResponse(w, http.StatusBadRequest, "Cart is empty")
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

	var totalOrderCost float64
	itemPrices := make(map[uuid.UUID]float64)

	for _, cartItem := range cartItems {
		var item model.Item
		err = ReadByID(&item, "items", item_columns, cartItem.Item_id.String())
		if err != nil {
			SendErrorResponse(w, err)
			return
		}

		itemPrices[cartItem.Item_id] = item.Price
		totalOrderCost += float64(cartItem.Quantity) * item.Price
	}

	orderMap := map[string]interface{}{
		"id":               uuid.New().String(),
		"total_order_cost": totalOrderCost,
		"customer_id":      user.ID.String(),
		"vendor_id":        userCart.Vendor_id.String(),
		"status":           model.PENDING,
		"created_at":       time.Now(),
		"updated_at":       time.Now(),
	}

	var userOrder model.Orders
	err = Create("orders", orderMap, &userOrder, tx)
	if err != nil {
		SendErrorResponse(w, err)
		return
	}

	for _, cartItem := range cartItems {
		var userOrderItem model.OrderItems
		price := itemPrices[cartItem.Item_id]
		orderItem := map[string]interface{}{
			"id":       uuid.New(),
			"order_id": userOrder.Id.String(),
			"item_id":  cartItem.Item_id,
			"quantity": int32(cartItem.Quantity),
			"price":    price,
		}
		err = Create("order_items", orderItem, &userOrderItem, tx)
		if err != nil {
			log.Println("Error on creating order_items!")
			SendErrorResponse(w, err)
			return
		}
	}

	err = DeleteByColumns("cart_items", map[string]interface{}{"cart_id": userCart.User_id.String()}, tx)
	if err != nil {
		SendErrorResponse(w, err)
		return
	}

	var updatedCart model.Carts
	err = UpdateById("carts", userCart.User_id.String(), map[string]interface{}{"vendor_id": nil}, &updatedCart, tx)
	if err != nil {
		SendErrorResponse(w, err)
		return
	}

	orderMap["cart"] = userCart

	SendJsonResponse(w, http.StatusCreated, orderMap)
}

func GetVendorOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userDetails, ok := ctx.Value(userContextKey).(UserDetails)
	if !ok {
		SendCustomeErrorResponse(w, http.StatusUnauthorized, "You are not authorized to do this action")
		return
	}

	user := userDetails.UserData
	userRoles := userDetails.UserRoles

	vendor_id := r.PathValue("id")
	_, err := uuid.Parse(vendor_id)
	if err != nil {
		SendCustomeErrorResponse(w, http.StatusBadRequest, "Not a valid vendor_id")
		return
	}

	hasRole := false
	for _, role := range userRoles {
		if role == "vendor" {
			hasRole = true
			break
		}
		if !hasRole {
			exist, err := RowExists("vendor_admins", map[string]interface{}{"user_id": user.ID.String(), "vendor_id": vendor_id})
			if err != nil {
				SendErrorResponse(w, err)
				return
			}
			if !exist {
				SendCustomeErrorResponse(w, http.StatusUnauthorized, "You are not authorized to do this action: You should be a vendor")
				return

			}
		}
	}

	exist, err := RowExists("vendor_admins", map[string]interface{}{"user_id": user.ID.String(), "vendor_id": vendor_id})
	if !exist {
		SendCustomeErrorResponse(w, http.StatusUnauthorized, "You are not authorized to do this action")
		return
	}

	type OrderResult struct {
		OrderID        uuid.UUID `db:"order_id" json:"order_id"`
		TableName      string    `db:"table_name" json:"table_name"`
		Status         string    `db:"status" json:"status"`
		TotalOrderCost float64   `db:"total_order_cost" json:"total_order_cost"`
		OrderNumber    int       `db:"order_number" json:"order_number"`
	}
	var result []OrderResult
	query := fmt.Sprintf(`
    SELECT 
        orders.id AS order_id, 
        tables.name AS table_name, 
        status, 
        total_order_cost, 
        ROW_NUMBER() OVER (ORDER BY orders.created_at) AS order_number 
    FROM orders 
    JOIN tables ON orders.customer_id = tables.customer_id 
    WHERE orders.vendor_id = '%s' 
      AND (orders.status = '%s' OR orders.status = '%s')
`, vendor_id, model.PENDING, model.PREPARING)
	err = db.Select(&result, query)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No Rows found")
			SendCustomeErrorResponse(w, http.StatusNotFound, "There is no orders!")
			return

		}
		log.Println("Error retrieving orders=> ", err)
		SendCustomeErrorResponse(w, http.StatusInternalServerError, "Error retrieving orders")
		return
	}

	SendJsonResponse(w, http.StatusOK, result)
}

func GetOrderItems(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_, ok := ctx.Value(userContextKey).(UserDetails)
	if !ok {
		SendCustomeErrorResponse(w, http.StatusUnauthorized, "You are not authorized to do this action")
		return
	}

	order_id := r.PathValue("id")
	// Validate UUID
	_, err := uuid.Parse(order_id)
	if err != nil {
		SendCustomeErrorResponse(w, http.StatusBadRequest, "Not a valid order_id")
		return
	}

	var order model.Orders
	err = ReadByID(&order, "orders", order_columns, order_id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No Rows found")
			SendCustomeErrorResponse(w, http.StatusNotFound, "There is no order with provided id!")
			return
		}
		log.Println("Error retrieving orders=> ", err)
		SendCustomeErrorResponse(w, http.StatusInternalServerError, "Error retrieving orders")
		return
	}

	type OrderItemDetails struct {
		Name           string  `db:"name" json:"name"`
		Quantity       int32   `db:"quantity" json:"quantity"`
		Price          float32 `db:"price" json:"price"`
		TotalItemPrice float64 `db:"total_item_price" json:"total_item_price"`
	}

	query := fmt.Sprintf(`
        SELECT items.name, order_items.quantity, order_items.price, 
               order_items.price * order_items.quantity AS total_item_price 
        FROM items
        JOIN order_items ON order_items.item_id = items.id
        WHERE order_items.order_id = '%s'`, order_id)

	var result []OrderItemDetails
	err = db.Select(&result, query)
	if err != nil {
		log.Println("Error retrieving order items =>", err)
		SendCustomeErrorResponse(w, http.StatusInternalServerError, "Error retrieving order items")
		return
	}

	SendJsonResponse(w, http.StatusOK, result)
}

func GetUserOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userDetails, ok := ctx.Value(userContextKey).(UserDetails)
	if !ok {
		SendCustomeErrorResponse(w, http.StatusUnauthorized, "You are not authorized to do this action")
		return
	}

	user := userDetails.UserData

	var orders []model.Orders

	query := fmt.Sprintf("SELECT * FROM orders WHERE customer_id = '%s' AND (status = '%s' OR status = '%s')", user.ID, model.PENDING, model.PREPARING)
	fmt.Println(query)
	err := db.Select(&orders, query)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No Rows found")
			HandelError(w, http.StatusNotFound, "There is no orders!")
			return

		}
		log.Println("Error retrieving orders => ", err)
		HandelError(w, http.StatusInternalServerError, "Error retrieving orders")
		return
	}

	SendJsonResponse(w, http.StatusOK, orders)
}

func UpdateOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userDetails, ok := ctx.Value(userContextKey).(UserDetails)
	if !ok {
		SendCustomeErrorResponse(w, http.StatusUnauthorized, "You are not authorized to do this action")
		return
	}

	user := userDetails.UserData
	userRoles := userDetails.UserRoles

	vendor_id := r.FormValue("vendor_id")
	_, err := uuid.Parse(vendor_id)
	if err != nil {
		log.Println("vendor_id = ", vendor_id)
		SendCustomeErrorResponse(w, http.StatusBadRequest, "Not a valid vendor_id")
		return
	}

	hasRole := false
	for _, role := range userRoles {
		if role == "vendor" {
			hasRole = true
			break
		}
		if !hasRole {
			exist, err := RowExists("vendor_admins", map[string]interface{}{"user_id": user.ID.String(), "vendor_id": vendor_id})
			if err != nil {
				SendErrorResponse(w, err)
				return
			}
			if !exist {
				SendCustomeErrorResponse(w, http.StatusUnauthorized, "You are not authorized to do this action: You should be a vendor ")
				return

			}
		}
	}
	order_id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		log.Println("Error in id -> ", err)
		SendErrorResponse(w, ErrInvalidArgument)
		return
	}

	exist, err := RowExists("vendor_admins", map[string]interface{}{"user_id": user.ID.String(), "vendor_id": vendor_id})
	if !exist {
		SendCustomeErrorResponse(w, http.StatusUnauthorized, "You are not authorized to do this action")
		return
	}

	var order model.Orders
	err = ReadByID(&order, "orders", order_columns, order_id.String())
	if err != nil {
		SendErrorResponse(w, err)
		return
	}

	status := r.FormValue("status")
	log.Println("comming status -> ", status)
	switch status {
	case "1":
		status = string(model.PENDING)
	case "2":
		status = string(model.PREPARING)
	case "3":
		status = string(model.READY)
	default:
		status = string(model.PENDING)
	}

	orderMap := map[string]interface{}{
		"id":               order.Id,
		"total_order_cost": order.Total_order_cost,
		"customer_id":      order.Customer_id,
		"vendor_id":        order.Vendor_id,
		"status":           status,
		"created_at":       order.Created_at,
		"updated_at":       time.Now(),
	}

	var updatedOrder model.Orders

	err = UpdateById("orders", order_id, orderMap, &updatedOrder, nil)
	if err != nil {
		SendErrorResponse(w, err)
		return
	}

	SendJsonResponse(w, http.StatusOK, "Successfully updated")
}
