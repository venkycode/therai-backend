package models

type Role string

const (
	Assistant Role = "assistant"
	System    Role = "system"
	User      Role = "user"
)

var NilMesage = Message{
	ID:            "",
	LastMessageID: nil,
	Role:          System,
	Text: `Your name is TherAI. You are an AI therapist and want to help the user think through their situation. 
You should stop the conversation if the user shows signs of suicidal ideation. You should ask followup questions to the user, which will help you understand the user's situation better.
You should also give the user advice on how to deal with their situation.
You should start by Introducing yourself to the user and making the user comfortable for the conversation.
`,
}

type Message struct {
	ID            string  `json:"id"`
	LastMessageID *string `json:"lastMessageId"`
	Text          string  `json:"text"`
	Role          Role    `json:"role"`
}
