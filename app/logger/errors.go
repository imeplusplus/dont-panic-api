package logger

import (
	"fmt"
)

var (
	errorResourceAlreadyExistsMessage = LogEvent{1, "Resource %s already exists"}
	errorResourceNotFoundMessage      = LogEvent{2, "Resource %s not found"}
	errorNoVertexInResourceMessage    = LogEvent{3, "There is no vertex in resource %s"}
)

type ErrorResourceAlreadyExists struct{ ResourceName string }
type ErrorResourceNotFound struct{ ResourceName string }
type ErrorNoVertexInResource struct{ ResourceName string }

func (e ErrorResourceAlreadyExists) Error() error {
	return fmt.Errorf(errorResourceAlreadyExistsMessage.message, e.ResourceName)
}

func (e ErrorResourceNotFound) Error() error {
	return fmt.Errorf(errorResourceNotFoundMessage.message, e.ResourceName)
}

func (e ErrorNoVertexInResource) Error() error {
	return fmt.Errorf(errorNoVertexInResourceMessage.message, e.ResourceName)
}
