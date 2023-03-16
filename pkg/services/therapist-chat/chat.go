package therapistchat

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/venkycode/therai-backend/pkg/models"
	openai "github.com/venkycode/therai-backend/pkg/utils/open-ai"
)

type MessageDB map[string]models.Message

var conversationDB = map[string]MessageDB{}

type TherapistChat interface {
	Chat(conversationID string, lastMessageID *string, message string) (*models.Message, error)
	Begin(conversationID string) (*models.Message, error)
}

type therapistChat struct {
	aiClient openai.OpenAIClient
}

func NewTherapistChat(aiClient openai.OpenAIClient) TherapistChat {
	return &therapistChat{
		aiClient: aiClient,
	}
}

func (t *therapistChat) Begin(conversationID string) (*models.Message, error) {
	log.Println(len(conversationDB))
	therapistInitialResponse, err := t.getTherapistResponse(t.getConversation(conversationID, nil))
	if err != nil {
		return nil, err
	}
	therapistMessage := t.createNewMessage(conversationID, nil, therapistInitialResponse, models.Assistant)
	return &therapistMessage, nil
}

func (t *therapistChat) Chat(conversationID string, lastMessageID *string, userResponse string) (*models.Message, error) {
	newUserMessage := t.createNewMessage(conversationID, lastMessageID, userResponse, models.User)
	conversation := t.getConversation(conversationID, &newUserMessage.ID)
	therapistResponse, err := t.getTherapistResponse(conversation)
	if err != nil {
		return nil, err
	}
	therapistMessage := t.createNewMessage(conversationID, &newUserMessage.ID, therapistResponse, models.Assistant)
	return &therapistMessage, nil
}

func (t *therapistChat) getTherapistResponse(conversation *models.Conversation) (string, error) {
	response, err := t.aiClient.ChatCompletion(conversation)
	if err != nil {
		return "", fmt.Errorf("cannot aiClient.ChatCompletion: %w", err)
	}
	return response, nil
}
func (t *therapistChat) getConversation(conversationID string, lastMessageID *string) *models.Conversation {
	conversation := models.NewConversation()
	cur := lastMessageID
	for cur != nil {
		messageDB := getMessageDB(conversationID)
		lastMessage := messageDB[*cur]
		conversation.AddOldest(lastMessage)
		cur = lastMessage.LastMessageID
	}
	if cur == nil {
		conversation.AddOldest(models.NilMesage)
	}
	return conversation
}

func (t *therapistChat) createNewMessage(conversationID string, lastMessageID *string, response string, role models.Role) models.Message {
	userMessageID := uuid.New().String()
	newUserMessage := models.Message{
		ID:            userMessageID,
		LastMessageID: lastMessageID,
		Text:          response,
		Role:          role,
	}
	messageDB := getMessageDB(conversationID)
	messageDB[userMessageID] = newUserMessage
	return newUserMessage
}

func getMessageDB(conversationID string) MessageDB {
	_, ok := conversationDB[conversationID]
	if !ok {
		conversationDB[conversationID] = MessageDB{}
	}
	return conversationDB[conversationID]
}
