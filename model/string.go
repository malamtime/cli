package model

import (
	"regexp"
	"strings"
)

// MaskSensitiveTokens masks sensitive tokens (like JWT tokens) in shell commands.
func MaskSensitiveTokens(command string) string {
	// Define a regex to identify JWT tokens
	// A typical JWT token has 3 parts separated by dots and includes base64 characters
	jwtRegex := regexp.MustCompile(`ey[a-zA-Z0-9_-]+\.[a-zA-Z0-9_-]+\.[a-zA-Z0-9_-]+`)

	// Replace all matches with a masked value
	maskedCommand := jwtRegex.ReplaceAllStringFunc(command, func(token string) string {
		return maskToken(token)
	})

	// If there are more than 3 consecutive asterisks, reduce them to 3
	maskedCommand = regexp.MustCompile(`\*{4,}`).ReplaceAllString(maskedCommand, "***")

	return maskedCommand
}

// maskToken masks all but the first and last 4 characters of a token.
func maskToken(token string) string {
	if len(token) <= 8 {
		return strings.Repeat("*", len(token))
	}
	return token[:4] + strings.Repeat("*", len(token)-8) + token[len(token)-4:]
}
