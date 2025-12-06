package service

import (
	"errors"

	"github.com/umess/backend/internal/model"
	"github.com/umess/backend/internal/repository"
)

type MessageService struct {
	messageRepo      *repository.MessageRepository
	conversationRepo *repository.ConversationRepository
	userRepo         *repository.UserRepository
}

func NewMessageService() *MessageService {
	return &MessageService{
		messageRepo:      repository.NewMessageRepository(),
		conversationRepo: repository.NewConversationRepository(),
		userRepo:         repository.NewUserRepository(),
	}
}

func (s *MessageService) SendMessage(userID int, req *model.SendMessageRequest) (*model.Message, error) {
	// Verify user is participant
	participantIDs, err := s.conversationRepo.GetParticipants(req.ConversationID)
	if err != nil {
		return nil, errors.New("conversation not found")
	}

	isParticipant := false
	for _, pid := range participantIDs {
		if pid == userID {
			isParticipant = true
			break
		}
	}

	if !isParticipant {
		return nil, errors.New("unauthorized to send message")
	}

	// Create message
	message := &model.Message{
		ConversationID: req.ConversationID,
		SenderID:       userID,
		Content:        req.Content,
		Type:           req.Type,
		MediaURL:       req.MediaURL,
		Status:         model.MessageStatusSent,
	}

	err = s.messageRepo.Create(message)
	if err != nil {
		return nil, errors.New("failed to create message")
	}

	// Update conversation timestamp
	err = s.conversationRepo.UpdateTimestamp(req.ConversationID)
	if err != nil {
		// Log error but don't fail the request
	}

	// Load sender info
	sender, err := s.userRepo.GetByID(userID)
	if err == nil {
		sender.Password = ""
		message.Sender = sender
	}

	return message, nil
}

func (s *MessageService) GetMessages(conversationID, userID, limit, offset int) ([]*model.Message, error) {
	// Verify user is participant
	participantIDs, err := s.conversationRepo.GetParticipants(conversationID)
	if err != nil {
		return nil, errors.New("conversation not found")
	}

	isParticipant := false
	for _, pid := range participantIDs {
		if pid == userID {
			isParticipant = true
			break
		}
	}

	if !isParticipant {
		return nil, errors.New("unauthorized access")
	}

	// Get messages
	messages, err := s.messageRepo.GetByConversation(conversationID, limit, offset)
	if err != nil {
		return nil, err
	}

	// Load sender info for each message
	for _, msg := range messages {
		sender, err := s.userRepo.GetByID(msg.SenderID)
		if err == nil {
			sender.Password = ""
			msg.Sender = sender
		}
	}

	return messages, nil
}

func (s *MessageService) MarkAsDelivered(conversationID, userID int) error {
	return s.messageRepo.MarkAsDelivered(conversationID, userID)
}

func (s *MessageService) MarkAsRead(conversationID, userID int) error {
	return s.messageRepo.MarkAsRead(conversationID, userID)
}


