package auth

import (
	"bufio"
	"database/sql"
	"fmt"
	"strconv"

	"PVRGF/internal/menu"
	"golang.org/x/crypto/bcrypt"
)

func StartAuth(scanner *bufio.Scanner, db *sql.DB) {
	for {
		fmt.Println("\n\tWELCOME TO VAULTSEC")
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
			userID, success := login(scanner, db)
			if success {
				menu.ShowPostLoginMenu(scanner, db, userID)
			}
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

	err := RegisterUser(db, username, password)
	if err != nil {
		fmt.Printf("Registration failed: %v\n", err)
		return
	}
	fmt.Println("Registration successful!")
}

func RegisterUser(db *sql.DB, username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	_, err = db.Exec(
		"INSERT INTO users(username, password) VALUES(?, ?)",
		username, string(hashedPassword),
	)

	if err != nil {
		return fmt.Errorf("user already exists or database error")
	}
	return nil
}

func login(scanner *bufio.Scanner, db *sql.DB) (int, bool) {
	fmt.Print("Username: ")
	scanner.Scan()
	username := scanner.Text()

	fmt.Print("Password: ")
	scanner.Scan()
	password := scanner.Text()

	id, err := LoginUser(db, username, password)
	if err != nil {
		fmt.Printf("Login failed: %v\n", err)
		return 0, false
	}

	fmt.Println("Login successful. Welcome to your vault!")
	return id, true
}

func LoginUser(db *sql.DB, username, password string) (int, error) {
	var id int
	var hashedPassword string
	err := db.QueryRow(
		"SELECT id, password FROM users WHERE username=?",
		username,
	).Scan(&id, &hashedPassword)

	if err != nil {
		return 0, fmt.Errorf("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return 0, fmt.Errorf("invalid credentials")
	}

	return id, nil
}
