package controllers

import (
	"log"
	"net/http"
	"server-side/services"

	"github.com/go-michi/michi"
)

var serverPort string

func SetServerPort(port string) {
	serverPort = port
}

func Controllers() {
	r := michi.NewRouter()
	if r == nil {
		log.Fatalf("Router failed to initialize")
	}

	r.Use(services.CORS)

	r.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	r.Route("/", func(sub *michi.Router) {
		sub.HandleFunc("POST signup", services.SignUpNewUser)
		sub.HandleFunc("POST signin", services.SignIn)
		sub.HandleFunc("POST adminSignin", services.AdminSignin)

		sub.With(services.AuthenticateJWT).Group(func(auth *michi.Router) {
			handleUserRoutes(auth)
		})

		sub.With(services.AuthenticateJWT, services.Authorize("admin")).Group(func(admin *michi.Router) {
			handleAdminRoutes(admin)
			handleUserRoleRoutes(admin)
			handleVendorRoutes(admin)
			handleRoleRoutes(admin)
		})

		sub.With(services.AuthenticateJWT, services.Authorize("admin", "vendor")).Group(func(vendor *michi.Router) {
			handleVendorAdminRoutes(vendor)
			handelItemVendorRoutes(vendor)
		})
	})

	log.Printf("Starting server on port %s", serverPort)

	err := http.ListenAndServe(serverPort, r)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func handleUserRoutes(sub *michi.Router) {
	// users
	sub.HandleFunc("GET users", services.GetUsers)
	sub.HandleFunc("GET users/{id}", services.GetUserById)
	sub.HandleFunc("PUT users/{id}", services.UpdateUser)

	// vendors
	sub.HandleFunc("GET vendors", services.GetVendors)

	// items
	sub.HandleFunc("GET items", services.GetItems)
	sub.HandleFunc("GET items/{id}", services.GetItemById)

	// tables
	sub.HandleFunc("GET tables", services.GetTables)
	sub.HandleFunc("GET tables/{id}", services.GetTableById)

	sub.HandleFunc("PUT tables/{id}", services.UpdateTable)
	sub.HandleFunc("PUT tables/{id}/reserve", services.ReserveTable)

	// carts
	sub.HandleFunc("GET cart", services.GetUserCart)
	sub.HandleFunc("POST cart", services.CreateCart)
	sub.HandleFunc("POST cart/add", services.AddItemToCart)

	sub.HandleFunc("DELETE cart/empty", services.EmptyTheCart)
	sub.HandleFunc("DELETE cart/item", services.RemoveItemFromCart)

	// order
	sub.HandleFunc("POST cart/checkout", services.CreateOrder)
	sub.HandleFunc("GET orders/vendors/{id}", services.GetVendorOrders)
	sub.HandleFunc("GET orders", services.GetUserOrders)
	sub.HandleFunc("GET orders/{id}", services.GetOrderItems)
	sub.HandleFunc("PUT orders/{id}", services.UpdateOrder)
}

func handleVendorRoutes(sub *michi.Router) {
	sub.HandleFunc("GET vendors/{id}", services.GetVendorById)
	sub.HandleFunc("POST vendors", services.CreateNewVendor)
	sub.HandleFunc("PUT vendors/{id}", services.UpdateVendor)

	// tables
	sub.HandleFunc("DELETE tables/{id}/empty", services.EmptyTheTable)
}

func handleRoleRoutes(sub *michi.Router) {
	sub.HandleFunc("GET roles", services.GetAllRoles)
	sub.HandleFunc("GET roles/{id}", services.GetRoleById)
}

func handleUserRoleRoutes(sub *michi.Router) {
	sub.HandleFunc("POST users/grant-role", services.GrantRole)
	sub.HandleFunc("POST users/revoke-role", services.RevokeRole)
}

func handleVendorAdminRoutes(sub *michi.Router) {
	sub.HandleFunc("GET vendors/admins", services.GetAllVendorAdmins)
	sub.HandleFunc("GET vendors/{vendor_id}/admins", services.GetAllAdminsForVendor)
}

func handleAdminRoutes(sub *michi.Router) {
	sub.HandleFunc("DELETE users/{id}", services.DeleteUser)
	sub.HandleFunc("POST vendors/assign-admin", services.AssignAdminToVendor)
	sub.HandleFunc("POST vendors/revoke-admin", services.RevokeAdminFromVendor)
	sub.HandleFunc("DELETE vendors", services.DeleteAllVendors)
	sub.HandleFunc("DELETE vendors/{id}", services.DeleteVendorById)
}

func handelItemVendorRoutes(sub *michi.Router) {
	sub.HandleFunc("POST items", services.CreateNewItem)
	sub.HandleFunc("PUT items/{id}", services.UpdateItem)
	sub.HandleFunc("DELETE items/{id}", services.DeleteItemById)

	sub.HandleFunc("POST tables", services.CreateNewTable)
}
