package menu

import (
	"bufio"
	"crypto/rand"
	"database/sql"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"PVRGF/internal/rules"
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

		if choice == "3" {
			return
		}

		if choice != "1" && choice != "2" {
			fmt.Println("Invalid option")
			continue
		}

		fmt.Print("Enter Label (Gmail/Linkedin/Facebook/Instagram): ")
		scanner.Scan()
		label := strings.ToLower(scanner.Text())

		criteria, err := rules.LoadCriteria(label)
		if err != nil {
			fmt.Println(err)
			continue
		}

		switch choice {

		case "1":
			password := generateWithCriteria(criteria)

			if !validatePassword(password, criteria) {
				fmt.Println("Generated password does not meet criteria. Try again.")
				continue
			}

			fmt.Println("Generated Password:", password)
			savePassword(db, userID, password, label)

		case "2":
			fmt.Println("\nWebsite Password Requirements:")
			fmt.Println("---------------------------------")
			fmt.Println("Minimum Length      :", criteria.MinLength)
			fmt.Println("Minimum Uppercase   :", criteria.MinUppercase)
			fmt.Println("Minimum Lowercase   :", criteria.MinLowercase)
			fmt.Println("Minimum Numbers     :", criteria.MinNumbers)
			fmt.Println("Minimum Special     :", criteria.MinSpecial)
			fmt.Println("Allowed Special Chars:", criteria.AllowedSpecial)
			fmt.Println("---------------------------------")

			fmt.Print("Enter your password: ")
			scanner.Scan()
			password := scanner.Text()

			if !validatePassword(password, criteria) {
				fmt.Println("\nPassword does NOT match website criteria!")
				continue
			}

			savePassword(db, userID, password, label)
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

	return string(runes)
}

func generateWithCriteria(c *rules.Criteria) string {
	password := ""

	for i := 0; i < c.MinUppercase; i++ {
		password += randomChar(upperCharSet)
	}

	for i := 0; i < c.MinLowercase; i++ {
		password += randomChar(lowerCharSet)
	}

	for i := 0; i < c.MinNumbers; i++ {
		password += randomChar(numberCharSet)
	}

	for i := 0; i < c.MinSpecial; i++ {
		password += randomChar(c.AllowedSpecial)
	}

	allChars := lowerCharSet + upperCharSet + numberCharSet + c.AllowedSpecial
	remaining := c.MinLength - len(password)

	for i := 0; i < remaining; i++ {
		password += randomChar(allChars)
	}

	return shuffleString(password)
}

func validatePassword(pass string, c *rules.Criteria) bool {

	if len(pass) < c.MinLength {
		return false
	}

	var upper, lower, number, special int

	for _, ch := range pass {
		switch {
		case strings.ContainsRune(upperCharSet, ch):
			upper++
		case strings.ContainsRune(lowerCharSet, ch):
			lower++
		case strings.ContainsRune(numberCharSet, ch):
			number++
		case strings.ContainsRune(c.AllowedSpecial, ch):
			special++
		}
	}

	return upper >= c.MinUppercase &&
		lower >= c.MinLowercase &&
		number >= c.MinNumbers &&
		special >= c.MinSpecial
}
