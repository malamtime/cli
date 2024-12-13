package model

import "go.opentelemetry.io/otel"

type GinGraphQLContextType struct {
	IP     string
	UserID int
}

var commitID string
var modelTracer = otel.Tracer("model")

const MAX_BUFFER_SIZE = 512 * 1024 // 512Kb

func InjectVar(commitId string) {
	commitID = commitId
}
