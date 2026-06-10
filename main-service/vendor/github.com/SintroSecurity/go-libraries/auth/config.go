package auth

type Config struct {
	Key        string `mapstructure:"KEY"`
	CookieName string `mapstructure:"COOKIENAME"`
}
