package auth

import (
	"bufio"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"PVRGF/internal/menu"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("your-256-bit-secret") // In production, use env var

func StartAuth(scanner *bufio.Scanner, db *sql.DB) {
	for {
		fmt.Println("\n\tWELCOME TO VAULTSEC")
		fmt.Println("\t---------------------")
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

// GenerateToken creates a new JWT token for a user
func GenerateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateToken parses and validates a JWT token, returning the user_id
func ValidateToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := int(claims["user_id"].(float64))
		return userID, nil
	}

	return 0, fmt.Errorf("invalid token")
}
