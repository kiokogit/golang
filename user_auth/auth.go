package user_auth

import (
	"fmt"
	"log"
	"time"
)

func Authenticate() bool {
	fmt.Println("Please enter your username:")
	var userName string
	var userPass string
	fmt.Scanln(&userName)
	fmt.Println("Please enter your password:")
	fmt.Scanln(&userPass)

	fmt.Println("Authenticating user...Please wait")
	time.Sleep(2000000000)
	if userName == "kiokogit" {
		if userPass == "smart" {
			fmt.Println("User is authenticated! Proceeding...")
			return true
		}
	}

	log.Fatalln("Invalid username/password")
	return false
}
