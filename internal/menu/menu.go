package menu

import (
	"bufio"
	"crypto/rand"
	"database/sql"
	"fmt"
	"math/big"
	"os"
	"time"
)

var (
	lowerCharSet   = "abcdefghijklmnopqrstuvwxyz"
	upperCharSet   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	specialCharSet = "!@#$%^&*()-"
	numberCharSet  = "1234567890"

	minSpecialChar = 2
	minUpperChar   = 2
	minNumberChar  = 2
	passwordLength = 10
)

func ShowPostLoginMenu(scanner *bufio.Scanner, db *sql.DB, userID int) {
	for {
		fmt.Println("\n\tVAULT MENU")
		fmt.Println("\t----------")
		fmt.Println("1. Add Password")
		fmt.Println("2. View Passwords")
		fmt.Println("3. Logout")
		fmt.Print("Enter choice: ")

		scanner.Scan()
		choice := scanner.Text()

		switch choice {
		case "1":
			showAddPasswordMenu(scanner, db, userID)
		case "2":
			ViewPasswords(db, userID)
		case "3":
			fmt.Println("Logged out successfully.")
			return
		default:
			fmt.Println("Invalid option")
		}
	}
}

func showAddPasswordMenu(scanner *bufio.Scanner, db *sql.DB, userID int) {
	for {
		fmt.Println("\n\tADD PASSWORD")
		fmt.Println("\t------------")
		fmt.Println("1. Generate Random Password")
		fmt.Println("2. Save My Own Password")
		fmt.Println("3. Back")
		fmt.Print("Enter choice: ")

		scanner.Scan()
		choice := scanner.Text()

		fmt.Print("Enter label (Gmail/ Instagram/ Linkedin/ Facebook): ")
		scanner.Scan()
		label := scanner.Text()

		switch choice {
		case "1":
			password := generateRandomPassword()
			fmt.Println("Generated Password:", password)
			savePassword(db, userID, password, label)
		case "2":
			fmt.Print("Enter your password: ")
			scanner.Scan()
			password := scanner.Text()
			savePassword(db, userID, password, label)
		case "3":
			return
		default:
			fmt.Println("Invalid option")
		}
	}
}

func savePassword(db *sql.DB, userID int, password string, label string) {
	_, err := db.Exec(
		"INSERT INTO passwords(user_id, label, password, created_at) VALUES(?, ?, ?, ?)",
		userID, label, password, time.Now(),
	)

	if err != nil {
		fmt.Println("Error saving password:", err)
		return
	}

	fmt.Println("Password saved successfully!")
}

func ViewPasswords(db *sql.DB, userID int) {
	rows, err := db.Query(
		"SELECT id, label, password, created_at FROM passwords WHERE user_id=?",
		userID,
	)

	if err != nil {
		fmt.Println("Error retrieving passwords:", err)
		return
	}
	defer rows.Close()

	fmt.Println("\n\tYOUR SAVED PASSWORDS")
	fmt.Println("\t---------------------")

	for rows.Next() {
		var id int
		var label, password string
		var createdAt string

		rows.Scan(&id, &label, &password, &createdAt)

		fmt.Printf("ID: %d\nLabel: %s\nPassword: %s\nCreated: %s\n\n",
			id, label, password, createdAt)
	}
}

func generateRandomPassword() string {

	totalCharLenWithoutLowerChar := minUpperChar + minSpecialChar + minNumberChar

	if totalCharLenWithoutLowerChar >= passwordLength {
		fmt.Println("Invalid password configuration")
		os.Exit(1)
	}

	password := ""

	for i := 0; i < minSpecialChar; i++ {
		password += randomChar(specialCharSet)
	}

	for i := 0; i < minUpperChar; i++ {
		password += randomChar(upperCharSet)
	}

	for i := 0; i < minNumberChar; i++ {
		password += randomChar(numberCharSet)
	}

	remainingCharLen := passwordLength - totalCharLenWithoutLowerChar

	for i := 0; i < remainingCharLen; i++ {
		password += randomChar(lowerCharSet)
	}

	return shuffleString(password)
}

func randomChar(charset string) string {
	max := big.NewInt(int64(len(charset)))
	n, _ := rand.Int(rand.Reader, max)
	return string(charset[n.Int64()])
}

func shuffleString(s string) string {
	runes := []rune(s)

	for i := range runes {
		jBig, _ := rand.Int(rand.Reader, big.NewInt(int64(len(runes))))
		j := int(jBig.Int64())
		runes[i], runes[j] = runes[j], runes[i]
	}
	fmt.Println("-=-==-=--=-=-=--=-==")
	fmt.Println(string(runes))

	return string(runes)
}
