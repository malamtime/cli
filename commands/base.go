package commands

type errorResponse struct {
	ErrorCode    int    `json:"code"`
	ErrorMessage string `json:"error"`
}

var commitID string

func InjectVar(commitId string) {
	commitID = commitId
}
