package main

import (
	"bookingapp/checkout_checkin"
	"bookingapp/scraper_api"
	"net/http"

	"github.com/gin-gonic/gin"
)

var WelcomeMessage string = "Welcome to Unatum solutions."

func main() {
	// shared_utils.WelcomeCustomer("kiokogit", user_auth.Authenticate())
	// shared_utils.ChoicesForCustomer()
	// Initialize a new Gin router
	router := gin.Default()

	// Define a simple GET route
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Define a route group for user-related routes
	userGroup := router.Group("/books")
	{
		userGroup.GET("/checkout", checkout_checkin.CheckoutBookAPI)
		userGroup.GET("/scrap", scraper_api.ScrapeDataView)

	}

	// Start the server on port 8080
	router.Run(":8080")
}
