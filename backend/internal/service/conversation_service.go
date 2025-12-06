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
	err = s.convRepo.AddParticipant(conv.ID, userID1, "member")
	if err != nil {
		return nil, errors.New("failed to add participant")
	}

	err = s.convRepo.AddParticipant(conv.ID, userID2, "member")
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

// GetParticipantIDs gets all participant IDs for a conversation (for broadcasting)
func (s *ConversationService) GetParticipantIDs(conversationID int) ([]int, error) {
	return s.convRepo.GetParticipants(conversationID)
}

func (s *ConversationService) CreateGroup(userID int, name string, description string, participantIDs []int) (*model.ConversationWithParticipant, error) {
	// Check permission
	canCreate, err := s.userRepo.CheckCanCreateGroup(userID)
	if err != nil {
		return nil, errors.New("failed to check permission")
	}
	if !canCreate {
		return nil, errors.New("user does not have permission to create groups")
	}

	// Validate: at least 2 participants (excluding creator)
	if len(participantIDs) < 1 {
		return nil, errors.New("group must have at least 2 members (including creator)")
	}

	// Create group conversation
	conv := &model.Conversation{
		Type:        "group",
		Name:        name,
		Description: description,
	}

	err = s.convRepo.Create(conv)
	if err != nil {
		return nil, errors.New("failed to create group")
	}

	// Add creator as admin
	err = s.convRepo.AddParticipant(conv.ID, userID, "admin")
	if err != nil {
		return nil, errors.New("failed to add creator as admin")
	}

	// Add other participants as members
	for _, pid := range participantIDs {
		if pid != userID { // Don't add creator twice
			err = s.convRepo.AddParticipant(conv.ID, pid, "member")
			if err != nil {
				return nil, errors.New("failed to add participant")
			}
		}
	}

	// Populate participants before returning
	participantIDs, err = s.convRepo.GetParticipants(conv.ID)
	if err != nil {
		return nil, errors.New("failed to get participants")
	}

	users, err := s.userRepo.GetByIDs(participantIDs)
	if err != nil {
		return nil, errors.New("failed to get user details")
	}

	convWithParticipants := &model.ConversationWithParticipant{
		Conversation: *conv,
	}
	convWithParticipants.Participants = make([]model.User, len(users))
	for i, u := range users {
		u.Password = ""
		convWithParticipants.Participants[i] = *u
	}

	return convWithParticipants, nil
}

func (s *ConversationService) AddGroupMember(conversationID, userID, addedBy int) error {
	// Verify conversation is a group
	conv, err := s.convRepo.GetByID(conversationID)
	if err != nil {
		return errors.New("conversation not found")
	}
	if conv.Type != "group" {
		return errors.New("conversation is not a group")
	}

	// Check if addedBy is admin
	role, err := s.convRepo.GetParticipantRole(conversationID, addedBy)
	if err != nil || role != "admin" {
		return errors.New("only admins can add members")
	}

	// Check if user is already a participant
	participants, err := s.convRepo.GetParticipants(conversationID)
	if err != nil {
		return errors.New("failed to get participants")
	}
	for _, pid := range participants {
		if pid == userID {
			return errors.New("user is already a member")
		}
	}

	// Add user as member
	err = s.convRepo.AddParticipant(conversationID, userID, "member")
	if err != nil {
		return errors.New("failed to add member")
	}

	return nil
}

func (s *ConversationService) RemoveGroupMember(conversationID, userID, removedBy int) error {
	// Verify conversation is a group
	conv, err := s.convRepo.GetByID(conversationID)
	if err != nil {
		return errors.New("conversation not found")
	}
	if conv.Type != "group" {
		return errors.New("conversation is not a group")
	}

	// Check if removedBy is admin
	role, err := s.convRepo.GetParticipantRole(conversationID, removedBy)
	if err != nil || role != "admin" {
		return errors.New("only admins can remove members")
	}

	// Prevent removing yourself
	if userID == removedBy {
		return errors.New("admins cannot remove themselves")
	}

	// Check if user is a participant
	participants, err := s.convRepo.GetParticipants(conversationID)
	if err != nil {
		return errors.New("failed to get participants")
	}
	isParticipant := false
	for _, pid := range participants {
		if pid == userID {
			isParticipant = true
			break
		}
	}
	if !isParticipant {
		return errors.New("user is not a member")
	}

	// Remove user
	err = s.convRepo.RemoveParticipant(conversationID, userID)
	if err != nil {
		return errors.New("failed to remove member")
	}

	return nil
}

func (s *ConversationService) LeaveGroup(conversationID, userID int) error {
	// Verify conversation is a group
	conv, err := s.convRepo.GetByID(conversationID)
	if err != nil {
		return errors.New("conversation not found")
	}
	if conv.Type != "group" {
		return errors.New("conversation is not a group")
	}

	// Check if user is admin
	role, err := s.convRepo.GetParticipantRole(conversationID, userID)
	if err != nil {
		return errors.New("user is not a member")
	}

	// If admin, check if there are other admins
	if role == "admin" {
		adminCount, err := s.convRepo.GetAdminCount(conversationID)
		if err != nil {
			return errors.New("failed to check admin count")
		}
		if adminCount <= 1 {
			return errors.New("cannot leave group: you are the only admin")
		}
	}

	// Remove user
	err = s.convRepo.RemoveParticipant(conversationID, userID)
	if err != nil {
		return errors.New("failed to leave group")
	}

	return nil
}

func (s *ConversationService) UpdateGroupInfo(conversationID, userID int, name, description string) error {
	// Verify conversation is a group
	conv, err := s.convRepo.GetByID(conversationID)
	if err != nil {
		return errors.New("conversation not found")
	}
	if conv.Type != "group" {
		return errors.New("conversation is not a group")
	}

	// Check if user is admin
	role, err := s.convRepo.GetParticipantRole(conversationID, userID)
	if err != nil || role != "admin" {
		return errors.New("only admins can update group info")
	}

	// Update group info
	err = s.convRepo.UpdateGroupInfo(conversationID, name, description)
	if err != nil {
		return errors.New("failed to update group info")
	}

	return nil
}


