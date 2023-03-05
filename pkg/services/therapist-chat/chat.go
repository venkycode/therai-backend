package therapistchat

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/venkycode/therai-backend/pkg/models"
	openai "github.com/venkycode/therai-backend/pkg/utils/open-ai"
)

var messageDB = map[string]models.Message{}

type TherapistChat interface {
	Chat(lastMessageID *string, message string) (*models.Message, error)
	Begin() (*models.Message, error)
}

type therapistChat struct {
	aiClient openai.OpenAIClient
}

func NewTherapistChat(aiClient openai.OpenAIClient) TherapistChat {
	return &therapistChat{
		aiClient: aiClient,
	}
}

func (t *therapistChat) Begin() (*models.Message, error) {
	therapistInitialResponse, err := t.getTherapistResponse(t.getConversation(nil))
	if err != nil {
		return nil, err
	}
	therapistMessage := t.createNewMessage(nil, therapistInitialResponse, models.Assistant)
	return &therapistMessage, nil
}

func (t *therapistChat) Chat(lastMessageID *string, userResponse string) (*models.Message, error) {
	newUserMessage := t.createNewMessage(lastMessageID, userResponse, models.User)
	conversation := t.getConversation(&newUserMessage.ID)
	therapistResponse, err := t.getTherapistResponse(conversation)
	if err != nil {
		return nil, err
	}
	therapistMessage := t.createNewMessage(&newUserMessage.ID, therapistResponse, models.Assistant)
	return &therapistMessage, nil
}

func (t *therapistChat) getTherapistResponse(conversation *models.Conversation) (string, error) {
	response, err := t.aiClient.ChatCompletion(conversation)
	if err != nil {
		return "", fmt.Errorf("cannot aiClient.ChatCompletion: %w", err)
	}
	return response, nil
}
func (t *therapistChat) getConversation(lastMessageID *string) *models.Conversation {
	conversation := models.NewConversation()
	cur := lastMessageID
	for cur != nil {
		lastMessage := messageDB[*cur]
		conversation.AddOldest(lastMessage)
		cur = lastMessage.LastMessageID
	}
	if cur == nil {
		conversation.AddOldest(models.NilMesage)
	}
	return conversation
}

func (t *therapistChat) createNewMessage(lastMessageID *string, response string, role models.Role) models.Message {
	userMessageID := uuid.New().String()
	newUserMessage := models.Message{
		ID:            userMessageID,
		LastMessageID: lastMessageID,
		Text:          response,
		Role:          role,
	}
	messageDB[userMessageID] = newUserMessage
	return newUserMessage
}
