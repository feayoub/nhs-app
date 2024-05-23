package config

import "os"

type AppConfig struct {
	TrelloAPIKey   string
	TrelloAPIToken string
}

func LoadConfig() AppConfig {
	key := os.Getenv("TRELLO_API_KEY")
	token := os.Getenv("TRELLO_API_TOKEN")

	return AppConfig{
		TrelloAPIKey:   key,
		TrelloAPIToken: token,
	}
}
