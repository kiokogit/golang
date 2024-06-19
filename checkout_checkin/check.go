package checkout_checkin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CheckoutBook() bool {

	fmt.Println("Please enter the book Id that should be checked out:")
	var bookId string

	fmt.Scanln(&bookId)
	fmt.Println("Confirm. Are you sure? \n 1. Yes \n 2. No")
	var sure int

	fmt.Scanln(&sure)

	switch sure {
	case 1:
		fmt.Printf("Congratulations, %v has been checked out successfully.", bookId)
		return true
	case 2:
		fmt.Printf("Could not continue. Try again booking another")
		return false
	}
	return false
}

// CheckoutBookAPI Handler to get all users
func CheckoutBookAPI(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"details": "Ah, this has been successful"})
}
