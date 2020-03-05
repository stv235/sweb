package pages

type SessionConfig struct {
	CookieName string `toml:"cookie-name"`
	SignKey string `toml:"sign-key"`
	EncryptKey string `toml:"encrypt-key"`
	Timeout int64 `toml:"timeout"`
}
