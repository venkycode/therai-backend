package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/venkycode/therai-backend/pkg/models"
	therapistchat "github.com/venkycode/therai-backend/pkg/services/therapist-chat"
	"github.com/venkycode/therai-backend/pkg/utils/env"
	openai "github.com/venkycode/therai-backend/pkg/utils/open-ai"
)

func main() {
	env := env.NewEnv()
	client := openai.NewOpenAIClient(env)
	therapistChat := therapistchat.NewTherapistChat(client)

	fmt.Println("TherAI is here to help you")
	message, err := therapistChat.Begin()
	if err != nil {
		panic(err)
	}
	if message.Role != models.Assistant {
		panic("role is not assistant, role:" + string(message.Role))
	}
	fmt.Println("TherAI: " + message.Text)
	lastMessageID := &message.ID
	in := bufio.NewReader(os.Stdin)

	for lastMessageID != nil {
		in.Discard(in.Buffered())
		fmt.Print("\n You: ")
		userResponse, err := in.ReadString('\n')
		if err != nil {
			panic(err)
		}
		message, err := therapistChat.Chat(lastMessageID, userResponse)
		if err != nil {
			panic(err)
		}
		fmt.Println("TherAI: " + message.Text)
		lastMessageID = &message.ID

	}

}
