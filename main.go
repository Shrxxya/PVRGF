package main

import (
	"bufio"
	"os"

	"PVRGF/internal/auth"
	"PVRGF/internal/db"
)

func main() {
	database := db.InitDB()
	defer database.Close()

	scanner := bufio.NewScanner(os.Stdin)
	auth.StartAuth(scanner, database)
}
