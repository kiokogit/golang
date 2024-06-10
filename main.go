package main

import (
	"bookingapp/shared_utils"
	"bookingapp/user_auth"
)

var WelcomeMessage string = "Welcome to Unatum solutions."

func main() {
	shared_utils.WelcomeCustomer("kiokogit", user_auth.Authenticate())
	shared_utils.ChoicesForCustomer()
}
