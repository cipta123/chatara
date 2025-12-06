package main

import (
	"fmt"
	"log"

	"github.com/umess/backend/internal/repository"
	"github.com/umess/backend/internal/service"
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

	// Test creating conversation
	userID1 := 2 // Aira
	userID2 := 1 // cipta

	fmt.Printf("🔍 Testing conversation creation:\n")
	fmt.Printf("   User 1: %d\n", userID1)
	fmt.Printf("   User 2: %d\n\n", userID2)

	convService := service.NewConversationService()
	conversation, err := convService.CreateDirectConversation(userID1, userID2)
	if err != nil {
		log.Fatal("Error creating conversation:", err)
	}

	fmt.Printf("✅ Conversation created:\n")
	fmt.Printf("   ID: %d\n", conversation.ID)
	fmt.Printf("   Type: %s\n", conversation.Type)
	fmt.Printf("   Participants count: %d\n", len(conversation.Participants))
	
	if len(conversation.Participants) == 0 {
		fmt.Println("\n❌ ERROR: Participants array is empty!")
		
		// Check participants in database
		convRepo := repository.NewConversationRepository()
		participantIDs, err := convRepo.GetParticipants(conversation.ID)
		if err != nil {
			fmt.Printf("   Error getting participants: %v\n", err)
		} else {
			fmt.Printf("   Participants in DB: %v\n", participantIDs)
		}
	} else {
		fmt.Println("\n✅ Participants:")
		for i, p := range conversation.Participants {
			fmt.Printf("   %d. ID: %d, Username: %s, Email: %s\n", i+1, p.ID, p.Username, p.Email)
		}
	}
}

