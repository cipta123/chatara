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

	// Validate reply_to_id if provided
	if req.ReplyToID != nil {
		replyToMsg, err := s.messageRepo.GetByID(*req.ReplyToID)
		if err != nil {
			return nil, errors.New("reply to message not found")
		}
		// Verify reply_to message is in the same conversation
		if replyToMsg.ConversationID != req.ConversationID {
			return nil, errors.New("reply to message must be in the same conversation")
		}
	}

	// Create message
	message := &model.Message{
		ConversationID: req.ConversationID,
		SenderID:       userID,
		Content:        req.Content,
		Type:           req.Type,
		MediaURL:       req.MediaURL,
		ReplyToID:      req.ReplyToID,
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

	// Load reply_to message if exists
	if message.ReplyToID != nil {
		replyTo, err := s.messageRepo.GetByID(*message.ReplyToID)
		if err == nil {
			// Load sender of reply_to message
			replyToSender, err := s.userRepo.GetByID(replyTo.SenderID)
			if err == nil {
				replyToSender.Password = ""
				replyTo.Sender = replyToSender
			}
			message.ReplyTo = replyTo
		}
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

	// Load sender info and reply_to for each message
	for _, msg := range messages {
		sender, err := s.userRepo.GetByID(msg.SenderID)
		if err == nil {
			sender.Password = ""
			msg.Sender = sender
		}

		// Load reply_to message if exists
		if msg.ReplyToID != nil {
			replyTo, err := s.messageRepo.GetByID(*msg.ReplyToID)
			if err == nil {
				// Load sender of reply_to message
				replyToSender, err := s.userRepo.GetByID(replyTo.SenderID)
				if err == nil {
					replyToSender.Password = ""
					replyTo.Sender = replyToSender
				}
				msg.ReplyTo = replyTo
			}
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

func (s *MessageService) DeleteMessage(messageID int, userID int) error {
	// Get message to verify it exists and user is sender
	message, err := s.messageRepo.GetByID(messageID)
	if err != nil {
		return errors.New("message not found")
	}

	// Only sender can delete their own message
	if message.SenderID != userID {
		return errors.New("unauthorized to delete message")
	}

	// Verify user is still participant in conversation
	participantIDs, err := s.conversationRepo.GetParticipants(message.ConversationID)
	if err != nil {
		return errors.New("conversation not found")
	}

	isParticipant := false
	for _, pid := range participantIDs {
		if pid == userID {
			isParticipant = true
			break
		}
	}

	if !isParticipant {
		return errors.New("unauthorized access")
	}

	// Perform soft delete
	return s.messageRepo.Delete(messageID, userID)
}


