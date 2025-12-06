package main

import (
	"fmt"
	"log"
	"os"

	"github.com/umess/backend/internal/repository"
	"github.com/umess/backend/pkg/config"
	"github.com/umess/backend/pkg/database"
	"github.com/umess/backend/pkg/utils"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run reset_password/main.go <email> <new_password>")
		fmt.Println("Example: go run reset_password/main.go aira@gmail.com newpassword123")
		os.Exit(1)
	}

	email := os.Args[1]
	newPassword := os.Args[2]

	if len(newPassword) < 6 {
		fmt.Println("❌ Password harus minimal 6 karakter")
		os.Exit(1)
	}

	// Load configuration
	if err := config.Load(); err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	// Get user
	userRepo := repository.NewUserRepository()
	user, err := userRepo.GetByEmail(email)
	if err != nil {
		fmt.Printf("❌ User dengan email %s tidak ditemukan\n", email)
		os.Exit(1)
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		log.Fatal("Failed to hash password:", err)
	}

	// Update password in database
	query := `UPDATE users SET password = $1, updated_at = NOW() WHERE id = $2`
	_, err = database.DB.Exec(query, hashedPassword, user.ID)
	if err != nil {
		log.Fatal("Failed to update password:", err)
	}

	fmt.Printf("✅ Password berhasil direset untuk user: %s (%s)\n", user.Username, email)
	fmt.Printf("   Password baru: %s\n", newPassword)
	fmt.Println("\n💡 Sekarang Anda bisa login dengan password baru ini")
}
