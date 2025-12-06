package service

import (
	"errors"

	"github.com/umess/backend/internal/model"
	"github.com/umess/backend/internal/repository"
)

type ConversationService struct {
	convRepo  *repository.ConversationRepository
	userRepo  *repository.UserRepository
	messageRepo *repository.MessageRepository
}

func NewConversationService() *ConversationService {
	return &ConversationService{
		convRepo:    repository.NewConversationRepository(),
		userRepo:    repository.NewUserRepository(),
		messageRepo: repository.NewMessageRepository(),
	}
}

func (s *ConversationService) CreateDirectConversation(userID1, userID2 int) (*model.ConversationWithParticipant, error) {
	// Check if conversation already exists
	existingConv, err := s.convRepo.FindDirectConversation(userID1, userID2)
	if err != nil {
		return nil, err
	}

	if existingConv != nil {
		// Populate participants for existing conversation
		participantIDs, err := s.convRepo.GetParticipants(existingConv.ID)
		if err == nil {
			users, err := s.userRepo.GetByIDs(participantIDs)
			if err == nil {
				convWithParticipants := &model.ConversationWithParticipant{
					Conversation: *existingConv,
				}
				convWithParticipants.Participants = make([]model.User, len(users))
				for i, u := range users {
					u.Password = ""
					convWithParticipants.Participants[i] = *u
				}
				return convWithParticipants, nil
			}
		}
		// Fallback: return basic conversation
		return &model.ConversationWithParticipant{
			Conversation: *existingConv,
			Participants: []model.User{},
		}, nil
	}

	// Create new conversation
	conv := &model.Conversation{
		Type: "direct",
	}

	err = s.convRepo.Create(conv)
	if err != nil {
		return nil, errors.New("failed to create conversation")
	}

	// Add participants
	err = s.convRepo.AddParticipant(conv.ID, userID1)
	if err != nil {
		return nil, errors.New("failed to add participant")
	}

	err = s.convRepo.AddParticipant(conv.ID, userID2)
	if err != nil {
		return nil, errors.New("failed to add participant")
	}

	// Populate participants before returning
	participantIDs, err := s.convRepo.GetParticipants(conv.ID)
	if err != nil {
		return nil, errors.New("failed to get participants: " + err.Error())
	}
	
	if len(participantIDs) == 0 {
		return nil, errors.New("no participants found for conversation")
	}
	
	users, err := s.userRepo.GetByIDs(participantIDs)
	if err != nil {
		return nil, errors.New("failed to get user details: " + err.Error())
	}
	
	if len(users) == 0 {
		return nil, errors.New("no users found for participant IDs")
	}
	
	// Convert to ConversationWithParticipant format
	convWithParticipants := &model.ConversationWithParticipant{
		Conversation: *conv,
	}
	convWithParticipants.Participants = make([]model.User, len(users))
	for i, u := range users {
		u.Password = "" // Clear password
		convWithParticipants.Participants[i] = *u
	}
	
	return convWithParticipants, nil
}

func (s *ConversationService) GetUserConversations(userID int) ([]*model.ConversationWithParticipant, error) {
	conversations, err := s.convRepo.GetUserConversations(userID)
	if err != nil {
		return []*model.ConversationWithParticipant{}, err
	}
	
	// Ensure we always return at least an empty slice
	if conversations == nil {
		conversations = []*model.ConversationWithParticipant{}
	}

	// Populate participants and last message for each conversation
	for _, conv := range conversations {
		participantIDs, err := s.convRepo.GetParticipants(conv.ID)
		if err != nil {
			continue
		}

		users, err := s.userRepo.GetByIDs(participantIDs)
		if err == nil {
			conv.Participants = make([]model.User, len(users))
			for i, u := range users {
				u.Password = "" // Clear password
				conv.Participants[i] = *u
			}
		}

		lastMessage, err := s.messageRepo.GetLastMessage(conv.ID)
		if err == nil && lastMessage != nil {
			sender, err := s.userRepo.GetByID(lastMessage.SenderID)
			if err == nil {
				sender.Password = ""
				lastMessage.Sender = sender
			}
			conv.LastMessage = lastMessage
		}
	}

	return conversations, nil
}

func (s *ConversationService) GetConversation(conversationID, userID int) (*model.ConversationWithParticipant, error) {
	conv, err := s.convRepo.GetByID(conversationID)
	if err != nil {
		return nil, err
	}

	// Check if user is participant
	participantIDs, err := s.convRepo.GetParticipants(conversationID)
	if err != nil {
		return nil, errors.New("failed to get participants")
	}

	isParticipant := false
	for _, pid := range participantIDs {
		if pid == userID {
			isParticipant = true
			break
		}
	}

	if !isParticipant {
		return nil, errors.New("unauthorized access to conversation")
	}

	convWithParticipant := &model.ConversationWithParticipant{
		Conversation: *conv,
	}

	users, err := s.userRepo.GetByIDs(participantIDs)
	if err == nil {
		convWithParticipant.Participants = make([]model.User, len(users))
		for i, u := range users {
			u.Password = ""
			convWithParticipant.Participants[i] = *u
		}
	}

	return convWithParticipant, nil
}


