package main

import (
	"fmt"
	"log"

	"github.com/umess/backend/pkg/config"
	"github.com/umess/backend/pkg/database"
)

func main() {
	// Load configuration
	if err := config.Load(); err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	// Get all users
	rows, err := database.DB.Query("SELECT id, username, email, created_at FROM users ORDER BY id")
	if err != nil {
		log.Fatal("Failed to query users:", err)
	}
	defer rows.Close()

	fmt.Println("📋 Daftar User yang Terdaftar:")
	fmt.Println("=====================================")
	
	userCount := 0
	for rows.Next() {
		var id int
		var username, email string
		var createdAt string
		
		if err := rows.Scan(&id, &username, &email, &createdAt); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		
		userCount++
		fmt.Printf("\n%d. ID: %d\n", userCount, id)
		fmt.Printf("   Username: %s\n", username)
		fmt.Printf("   Email: %s\n", email)
		fmt.Printf("   Created: %s\n", createdAt)
	}

	if userCount == 0 {
		fmt.Println("\n❌ Tidak ada user yang terdaftar")
		fmt.Println("\n💡 Silakan register user baru melalui frontend di http://localhost:3000/register")
	} else {
		fmt.Printf("\n✅ Total: %d user(s)\n", userCount)
	}
}
