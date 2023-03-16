package main

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	therapistchat "github.com/venkycode/therai-backend/pkg/services/therapist-chat"
	"github.com/venkycode/therai-backend/pkg/utils/env"
	openai "github.com/venkycode/therai-backend/pkg/utils/open-ai"
)

func main() {
	// Create a new bot API instance
	env := env.NewEnv()
	log.Println(env.TelegramBotKey())
	bot, err := tgbotapi.NewBotAPI(env.TelegramBotKey())
	if err != nil {
		log.Panic(err)
	}

	// Set up an update configuration
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	// Get updates from Telegram
	updates, err := bot.GetUpdatesChan(updateConfig)
	client := openai.NewOpenAIClient(env)

	therapistChat := therapistchat.NewTherapistChat(client)
	lastMSGDB := map[string]*string{}
	// Handle updates

	for update := range updates {

		if update.Message == nil {
			continue
		}
		// Check if the update is a command
		if update.Message.IsCommand() {
			// Create a new message
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

			// Switch on the command name
			switch update.Message.Command() {
			case "start":
				tmsg, err := therapistChat.Begin(fmt.Sprint(update.Message.Chat.ID))
				if err != nil {
					log.Panic(err)
				}
				msg.Text = tmsg.Text + " \n\n" + "You can use  /start to start a new conversation, /help to get a list of commands"
				lastMSGDB[fmt.Sprint(update.Message.Chat.ID)] = &tmsg.ID
			case "help":
				msg.Text = "You can use the following commands: /start, /help"
			default:
				msg.Text = "I don't know that command."
			}

			// Send the message
			_, err = bot.Send(msg)
			if err != nil {
				log.Println(err)
			}
		} else {
			// Create a new message
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

			// Check if the update is a text message
			if update.Message.Text != "" {
				// Get the last message
				lastMSG := lastMSGDB[fmt.Sprint(update.Message.Chat.ID)]
				if lastMSG == nil {
					tmsg, err := therapistChat.Begin(fmt.Sprint(update.Message.Chat.ID))
					if err != nil {
						log.Panic(err)
					}
					msg.Text = tmsg.Text
					lastMSGDB[fmt.Sprint(update.Message.Chat.ID)] = &tmsg.ID

				} else {
					tmsg, err := therapistChat.Chat(fmt.Sprint(update.Message.Chat.ID), lastMSG, update.Message.Text)
					if err != nil {
						log.Panic(err)
					}
					msg.Text = tmsg.Text
					lastMSGDB[fmt.Sprint(update.Message.Chat.ID)] = &tmsg.ID
				}
			}

			// Send the message
			_, err = bot.Send(msg)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
