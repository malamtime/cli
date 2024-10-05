package model

type MalamTimeConfig struct {
	Token       string
	APIEndpoint string
	// how often sync to server
	FlushCount int
	// how long the synced data would keep in db:
	// unit is days
	GCTime int
}

var DefaultConfig = MalamTimeConfig{
	Token:       "",
	APIEndpoint: "https://malamtime.com",
	FlushCount:  10,
	// 2 weeks by default
	GCTime: 14,
}
