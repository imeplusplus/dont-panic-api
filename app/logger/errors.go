package logger

import (
	"fmt"
)

var (
	errorResourceAlreadyExistsMessage = LogEvent{1, "resource %s already exists"}
	errorResourceNotFoundMessage      = LogEvent{2, "resource %s not found"}
	errorNoVerticesInQueryMessage     = LogEvent{3, "the Gremlin Query did not return any vertices"}
)

type ErrorResourceAlreadyExists struct{ ResourceName string }
type ErrorResourceNotFound struct{ ResourceName string }
type ErrorNoVerticesInQuery struct{}

func (e ErrorResourceAlreadyExists) Error() string {
	return fmt.Sprintf(errorResourceAlreadyExistsMessage.message, e.ResourceName)
}

func (e ErrorResourceNotFound) Error() string {
	return fmt.Sprintf(errorResourceNotFoundMessage.message, e.ResourceName)
}

func (e ErrorNoVerticesInQuery) Error() string {
	return fmt.Sprintf(errorNoVerticesInQueryMessage.message)
}
