package model

type GinGraphQLContextType struct {
	IP     string
	UserID int
}

var commitID string

func InjectVar(commitId string) {
	commitID = commitId
}
