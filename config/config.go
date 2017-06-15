package config

import (
	_ "github.com/joho/godotenv/autoload"
	"os"
	"strconv"
)

type Config struct {
	DebugEnabled      bool
	VisitorChannelID  string
	SlackSettings     SlackSettings
	WebServerSettings WebServerSettings
	DBSettings        DBSettings
}

type SlackSettings struct {
	BotAPIKey         string
	AppAPIKey         string
	VerificationToken string
}

type WebServerSettings struct {
	Port int
}

type DBSettings struct {
	Driver     string
	Connection string
}

func Load() *Config {
	debugEnabled := false
	port := 8000

	if os.Getenv("DEBUG_ENABLED") != "" {
		debugEnabled = true
	}

	if tmpPort, err := strconv.Atoi(os.Getenv("PORT")); err == nil {
		port = tmpPort
	}

	return &Config{
		DebugEnabled:     debugEnabled,
		VisitorChannelID: os.Getenv("VISITOR_CHANNEL_ID"),
		SlackSettings: SlackSettings{
			BotAPIKey:         os.Getenv("SLACK_BOT_API_KEY"),
			AppAPIKey:         os.Getenv("SLACK_APP_API_KEY"),
			VerificationToken: os.Getenv("SLACK_VERIFICATION_TOKEN"),
		},
		WebServerSettings: WebServerSettings{
			Port: port,
		},
		DBSettings: DBSettings{
			Driver:     "postgres",
			Connection: os.Getenv("DATABASE_URL"),
		},
	}
}
