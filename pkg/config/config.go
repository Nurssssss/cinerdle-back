package config

import "os"

type Config struct {
	Port         string
	DatabaseURL  string
	JWTSecret    string
	TMDbAPIKey   string
	AppBaseURL   string
	SMTPAddress  string
	SMTPPort     string
	SMTPLogin    string
	SMTPPassword string
}

func LoadConfig() *Config {
	return &Config{
		Port:         os.Getenv("PORT"),
		DatabaseURL:  os.Getenv("DATABASE_URL"),
		JWTSecret:    os.Getenv("JWT_SECRET"),
		TMDbAPIKey:   os.Getenv("TMDB_API_KEY"),
		AppBaseURL:   os.Getenv("APPBASE_URL"),
		SMTPAddress:  os.Getenv("SMTP_ADDRESS"),
		SMTPPort:     os.Getenv("SMTP_PORT"),
		SMTPLogin:    os.Getenv("SMTP_LOGIN"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
	}
}
