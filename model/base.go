package model

import "go.opentelemetry.io/otel"

type GinGraphQLContextType struct {
	IP     string
	UserID int
}

var commitID string
var modelTracer = otel.Tracer("model")

func InjectVar(commitId string) {
	commitID = commitId
}
