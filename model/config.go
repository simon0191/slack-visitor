package model

type Config struct {
	SlackSettings     SlackSettings
	DebugEnabled      bool
	WebServerSettings WebServerSettings
	DBSettings        DBSettings
}

type SlackSettings struct {
	BotAPIKey string
	AppAPIKey string
}

type WebServerSettings struct {
	Port int
}

type DBSettings struct {
	Driver     string
	Connection string
}
