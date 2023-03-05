package models

type Conversation struct {
	ChronologicalMessages []Message `json:"chronologicalMessages"`
}

func (c *Conversation) AddOldest(message Message) {
	c.ChronologicalMessages = append([]Message{message}, c.ChronologicalMessages...)
}

func (c *Conversation) AddNewest(message Message) {
	c.ChronologicalMessages = append(c.ChronologicalMessages, message)
}

func NewConversation() *Conversation {
	return &Conversation{
		ChronologicalMessages: []Message{},
	}
}
