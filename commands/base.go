package commands

import "github.com/malamtime/cli/model"

var commitID string

var configService model.ConfigService

func InjectVar(commitId string, cs model.ConfigService) {
	commitID = commitId
	configService = cs
}
