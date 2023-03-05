package env

import (
	"os"

	"github.com/joho/godotenv"
)

type Env interface {
	OpenAIApiKey() string
}

type env struct {
}

func NewEnv() Env {
	godotenv.Load()
	return &env{}
}

func (e *env) OpenAIApiKey() string { return os.Getenv("OPENAI_API_KEY") }
