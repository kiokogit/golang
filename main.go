package main

import (
	"bookingapp/auth"
	"bookingapp/checkout_checkin"
	"bookingapp/scraper_api"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize a new Gin router
	router := gin.Default()
	// Define a route group for user-related routes
	userGroup := router.Group("/books")
	{
		userGroup.GET("/checkout", checkout_checkin.CheckoutBookAPI)
		userGroup.GET("/scrap", scraper_api.ScrapeDataView)

	}
	authRouter := router.Group("/auth")
	{
		authRouter.GET("/login", auth.AuthenticateAPI)
	}

	// Start the server on port 8080
	router.Run(":8080")
}
