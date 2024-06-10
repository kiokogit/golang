package shared_utils

import (
	"fmt"
	"bookingapp/checkout_checkin"
)


var WelcomeMessage string = "Welcome to Unatum solutions."

func WelcomeCustomer(user_name string, authed bool) {
	if authed {
		fmt.Println(WelcomeMessage, user_name)
	} else {
		fmt.Println(WelcomeMessage)
	}
}


func ChoicesForCustomer() bool {

	var no_repeat bool = false
	for !no_repeat {
		fmt.Println("Choose your service from the choices below:")
		fmt.Println("1. Borrow a book\n 2. Return a book\n 3. Check account status\n 4. Exit peacefully")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			no_repeat = checkout_checkin.CheckoutBook()
		default:
			no_repeat = false
		}
	}
	return true

}
