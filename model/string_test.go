package model

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// commands := []string{
// 		`curl -X POST -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c" https://api.example.com/data`,
// 		`export JWT_TOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c`,
// 	}

// 	for _, command := range commands {
// 		fmt.Println("Original:", command)
// 		fmt.Println("Masked:  ", MaskSensitiveTokens(command))
// 	}

type stringTestSuite struct {
	suite.Suite
}

// before each test
func (s *stringTestSuite) SetupSuite() {
}

func (s *stringTestSuite) TestTrackWithSendData() {
	commands := []string{
		`curl -X POST -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c" https://api.shelltime.xyz/api/v1/track`,
		`export JWT_TOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c`,
	}

	result0 := MaskSensitiveTokens(commands[0])
	s.EqualValues(`curl -X POST -H "Authorization: Bearer eyJh***sw5c" https://api.shelltime.xyz/api/v1/track`, result0)
	result1 := MaskSensitiveTokens(commands[1])
	s.EqualValues("export JWT_TOKEN=eyJh***sw5c", result1)

}

func (s *stringTestSuite) TearDownSuite() {
}

func TestStringTestSuite(t *testing.T) {
	suite.Run(t, new(stringTestSuite))
}
