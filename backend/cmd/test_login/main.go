package main

import (
	"fmt"
	"log"

	"github.com/umess/backend/internal/model"
	"github.com/umess/backend/internal/repository"
	"github.com/umess/backend/internal/service"
	"github.com/umess/backend/pkg/config"
	"github.com/umess/backend/pkg/database"
	"github.com/umess/backend/pkg/utils"
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

	email := "aira@gmail.com"
	password := "upbjj@UT22"

	// Check if user exists - try case insensitive
	userRepo := repository.NewUserRepository()
	user, err := userRepo.GetByEmail(email)
	if err != nil {
		// Try with different case
		user, err = userRepo.GetByEmail("Aira@gmail.com")
		if err != nil {
		// Try direct database query
		var userID int
		query := "SELECT id FROM users WHERE LOWER(email) = LOWER($1)"
		err2 := database.DB.QueryRow(query, email).Scan(&userID)
			if err2 != nil {
				fmt.Printf("❌ User dengan email %s tidak ditemukan di database\n", email)
				fmt.Println("\nKemungkinan:")
				fmt.Println("1. User belum terdaftar")
				fmt.Println("2. Email salah")
				fmt.Println("\nSilakan register user terlebih dahulu melalui frontend")
				return
			}
			// Found via direct query, get full user
			user, err = userRepo.GetByID(userID)
			if err != nil {
				fmt.Printf("❌ Error getting user details: %v\n", err)
				return
			}
		}
	}

	fmt.Printf("✅ User ditemukan:\n")
	fmt.Printf("   ID: %d\n", user.ID)
	fmt.Printf("   Username: %s\n", user.Username)
	fmt.Printf("   Email: %s\n", user.Email)

	// Test password
	fmt.Printf("\n🔐 Testing password...\n")
	isValid := utils.CheckPasswordHash(password, user.Password)
	if isValid {
		fmt.Println("✅ Password BENAR - Login seharusnya berhasil")
	} else {
		fmt.Println("❌ Password SALAH - Password tidak cocok dengan yang tersimpan")
		fmt.Println("\nOpsi perbaikan:")
		fmt.Println("1. Reset password user ini")
		fmt.Println("2. Register ulang dengan email yang sama (akan error jika sudah ada)")
	}

	// Test login service
	fmt.Printf("\n🔍 Testing login service...\n")
	userService := service.NewUserService()
	loginReq := &model.LoginRequest{
		Email:    email,
		Password: password,
	}

	authResp, err := userService.Login(loginReq)

	if err != nil {
		fmt.Printf("❌ Login gagal: %s\n", err.Error())
		fmt.Println("\n💡 Solusi:")
		fmt.Println("   Jalankan: go run cmd/reset_password/main.go aira@gmail.com newpassword")
	} else {
		fmt.Println("✅ Login berhasil!")
		if len(authResp.Token) > 20 {
			fmt.Printf("   Token: %s...\n", authResp.Token[:20])
		} else {
			fmt.Printf("   Token: %s\n", authResp.Token)
		}
	}
}
