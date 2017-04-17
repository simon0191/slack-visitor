package model

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
