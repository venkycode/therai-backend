package env

import (
	"os"

	"github.com/joho/godotenv"
)

type Env interface {
	OpenAIApiKey() string
	TelegramBotKey() string
}

type env struct {
}

func NewEnv() Env {
	godotenv.Load()
	return &env{}
}

func (e *env) OpenAIApiKey() string   { return os.Getenv("OPENAI_API_KEY") }
func (e *env) TelegramBotKey() string { return os.Getenv("TELEGRAM_BOT_API_KEY") }
