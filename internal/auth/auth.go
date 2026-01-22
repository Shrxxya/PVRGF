package auth

import (
	"bufio"
	"database/sql"
	"fmt"
	"strconv"
)

func StartAuth(scanner *bufio.Scanner, db *sql.DB) {
	for {
		fmt.Println("\n\t\tWELCOME TO VAULTSEC")
		fmt.Println("\t--------------------")
		fmt.Println("1. Register")
		fmt.Println("2. Login")
		fmt.Println("3. Exit")
		fmt.Print("Enter choice: ")

		scanner.Scan()
		choice, _ := strconv.Atoi(scanner.Text())

		switch choice {
		case 1:
			register(scanner, db)
		case 2:
			login(scanner, db)
		case 3:
			fmt.Println("Thank you!")
			return
		default:
			fmt.Println("Invalid option")
		}
	}
}

func register(scanner *bufio.Scanner, db *sql.DB) {
	fmt.Print("Username: ")
	scanner.Scan()
	username := scanner.Text()

	fmt.Print("Password: ")
	scanner.Scan()
	password := scanner.Text()

	_, err := db.Exec(
		"INSERT INTO users(username, password) VALUES(?, ?)",
		username, password,
	)

	if err != nil {
		fmt.Println("User already exists!")
		return
	}
	fmt.Println("Registration successful!")
}

func login(scanner *bufio.Scanner, db *sql.DB) {
	fmt.Print("Username: ")
	scanner.Scan()
	username := scanner.Text()

	fmt.Print("Password: ")
	scanner.Scan()
	password := scanner.Text()

	row := db.QueryRow(
		"SELECT id FROM users WHERE username=? AND password=?",
		username, password,
	)

	var id int
	err := row.Scan(&id)
	if err != nil {
		fmt.Println("Invalid credentials")
		return
	}

	fmt.Println("Login successful. Welcome to your vault!")
}
