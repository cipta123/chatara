package main

import (
	"fmt"
	"log"

	"github.com/umess/backend/internal/repository"
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

	// Test search
	userRepo := repository.NewUserRepository()
	
	// Search for "aira" excluding user ID 5 (ciptas)
	query := "aira"
	excludeUserID := 5
	
	fmt.Printf("🔍 Testing search for: '%s' (excluding user ID: %d)\n\n", query, excludeUserID)
	
	users, err := userRepo.SearchUsers(query, excludeUserID, 20)
	if err != nil {
		log.Fatal("Search error:", err)
	}

	if len(users) == 0 {
		fmt.Println("❌ No users found")
		fmt.Println("\n💡 Checking all users in database...")
		
		// Get all users
		allUsers, _ := userRepo.GetByIDs([]int{1, 2, 5})
		for _, u := range allUsers {
			fmt.Printf("   - ID: %d, Username: %s, Email: %s\n", u.ID, u.Username, u.Email)
		}
	} else {
		fmt.Printf("✅ Found %d user(s):\n\n", len(users))
		for _, user := range users {
			fmt.Printf("   ID: %d\n", user.ID)
			fmt.Printf("   Username: %s\n", user.Username)
			fmt.Printf("   Email: %s\n", user.Email)
			fmt.Println()
		}
	}
}
