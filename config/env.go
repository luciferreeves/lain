package config

import "time"

type server struct {
	Host           string   `env:"SERVER_HOST" default:"localhost"`
	Port           int      `env:"SERVER_PORT" default:"8080"`
	AppSecret      string   `env:"APP_SECRET" default:"mysecret"`
	AppName        string   `env:"APP_NAME" default:"Lain Mail"`
	AppDescription string   `env:"APP_DESCRIPTION" default:"Present day, present time!"`
	AppEngine      string   `env:"APP_ENGINE" default:"Lain"`
	AllowedDomains []string `env:"ALLOWED_DOMAINS" default:"localhost"`
	DevMode        bool     `env:"DEV_MODE" default:"true"`
}

type mail struct {
	IMAPHost string `env:"IMAP_HOST" default:""`
	IMAPPort int    `env:"IMAP_PORT" default:"993"`
	IMAPTLS  bool   `env:"IMAP_TLS" default:"true"`
	SMTPHost string `env:"SMTP_HOST" default:""`
	SMTPPort int    `env:"SMTP_PORT" default:"587"`
	SMTPTLS  bool   `env:"SMTP_TLS" default:"true"`
}

type database struct {
	Host     string `env:"DB_HOST" default:"localhost"`
	Port     int    `env:"DB_PORT" default:"5432"`
	Username string `env:"DB_USER" default:"postgres"`
	Password string `env:"DB_PASS" default:""`
	Name     string `env:"DB_NAME" default:"lain"`
	SSLMode  string `env:"DB_SSLMODE" default:"disable"`
}

type minio struct {
	Endpoint   string `env:"MINIO_ENDPOINT" default:"localhost:9000"`
	AccessKey  string `env:"MINIO_ACCESS_KEY" default:""`
	SecretKey  string `env:"MINIO_SECRET_KEY" default:""`
	BucketName string `env:"MINIO_BUCKET_NAME" default:"lain"`
	UseSSL     bool   `env:"MINIO_USE_SSL" default:"false"`
}

type ai struct {
	URL     string `env:"AI_SERVER_URL" default:""`
	AuthKey string `env:"AI_SERVER_AUTH_KEY" default:""`
}

type session struct {
	CookieDomain   string        `env:"SESSION_COOKIE_DOMAIN" default:"localhost"`
	CookieName     string        `env:"SESSION_COOKIE_NAME" default:"lain_session"`
	CookiePath     string        `env:"SESSION_COOKIE_PATH" default:"/"`
	CookieSameSite string        `env:"SESSION_COOKIE_SAME_SITE" default:"Lax"`
	CookieSecure   bool          `env:"SESSION_SECURE_COOKIE" default:"false"`
	CookieTimeout  time.Duration `env:"SESSION_TIMEOUT" default:"24h"`
}
