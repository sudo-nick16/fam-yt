package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/sudo-nick16/env"
)

type Config struct {
	MongoUri     string
	Port         string
	YtApiKey     string
	DbName       string
	PollInterval int
}

func GetConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Printf("[ERROR] %v", err)
	}

	config := &Config{}
	config.MongoUri = env.GetEnv("MONGO_URI", "")
	config.Port = env.GetEnv("PORT", ":5000")
	config.YtApiKey = env.GetEnv("YT_API_KEY", "")
	config.DbName = env.GetEnv("DB_NAME", "fam-yt-dev")
	config.PollInterval = env.GetEnvAsInt("POLL_INTERVAL", 10)

	return config
}
