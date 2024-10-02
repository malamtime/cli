package internal

type MalamTimeConfig struct {
	Token       string
	APIEndpoint string
}

var DefaultConfig = MalamTimeConfig{
	Token:       "",
	APIEndpoint: "https://malamtime.com",
}
