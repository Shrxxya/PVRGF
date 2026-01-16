package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type User struct {
	username string
	password string
}

var users []User

func main() {

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("\n\t\tWELCOME TO VAULTSEC")
		fmt.Print("\t    ---------------------------")
		fmt.Print("\n1.Register")
		fmt.Println("\n2.Login")
		fmt.Println("3.Exit")
		fmt.Print("Please Enter your choice: ")

		scanner.Scan()
		input := scanner.Text()

		choice, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Please enter a valid number!")
			continue
		}

		switch choice {
		case 1:
			register(scanner)
		case 2:
			login()
		case 3:
			fmt.Println("Thank you!")
			return
		default:
			fmt.Println("Please enter a valid option!")
		}
	}
}

// Function to Register
func register(scanner *bufio.Scanner) {

	fmt.Print("Enter the Username: ")
	scanner.Scan()
	username := scanner.Text()

	fmt.Print("Enter the Password: ")
	scanner.Scan()
	password := scanner.Text()

	if username == "" || password == "" {
		fmt.Println("Username or password should not be empty!")
		return
	}

	for _, u := range users {
		if u.username == username {
			fmt.Println("User already exists!")
			return
		}
	}
	newUser := User{
		username: username,
		password: password,
	}

	users = append(users, newUser)

	fmt.Println("\nRegistration successful!")
}

// Function to Login
func login() {
	exists := false
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("\nEnter your username: ")
	scanner.Scan()
	username := scanner.Text()
	fmt.Println("\nEnter your password: ")
	scanner.Scan()
	pwd := scanner.Text()

	for _, value := range users {
		if value.username == username && value.password == pwd {
			exists = true
			break
		}
	}
	if exists {
		fmt.Printf("\nWelcome %s to your Vault!\n", username)
	} else {
		fmt.Println("\nUser with this username does not exist. Please register yourself.")
	}
}
