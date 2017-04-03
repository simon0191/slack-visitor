package model

type Config struct {
	SlackSettings     SlackSettings
	DebugEnabled      bool
	WebServerSettings WebServerSettings
}

type SlackSettings struct {
	BotAPIKey string
	AppAPIKey string
}

type WebServerSettings struct {
	Port int
}
