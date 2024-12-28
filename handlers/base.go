package handlers

import "github.com/malamtime/cli/model"

var stConfig model.ConfigService

func Init(cs model.ConfigService) {
	stConfig = cs
}
