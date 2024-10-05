package model

var commitID string

func InjectVar(commitId string) {
	commitID = commitId
}
