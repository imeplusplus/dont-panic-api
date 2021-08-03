package logger

import (
	"fmt"
)

var (
	resourceCreatedMessage             = LogEvent{1001, "Resource %s created.\nContent:\n%s"}
	resourceUpdatedMessage             = LogEvent{1002, "Resource %s updated.\nNew Content:\n%s"}
	resourceReadMessage                = LogEvent{1003, "Resource %s read.\nContent:\n%s"}
	resourceDeletedMessage             = LogEvent{1004, "Resource %s deleted.\n"}
	failedToDecodeJSONMessage          = LogEvent{1005, "Failed to parse JSON into %s.\n"}
	failedToEncodeJSONMessage          = LogEvent{1006, "Failed to encode %s into JSON format.\n"}
	failedToExecuteGremlinQueryMessage = LogEvent{1007, "Failed to execute Gremlin query.\n"}
)

type ResourceCreated struct {
	ResourceName    string
	ResourceContent string
}
type ResourceUpdated struct {
	ResourceName    string
	ResourceContent string
}
type ResourceRead struct {
	ResourceName    string
	ResourceContent string
}
type ResourceDeleted struct{ ResourceName string }
type FailedToDecodeJSON struct{ Resource string }
type FailedToEncodeJSON struct{ Resource string }
type FailedToExecuteGremlinQuery struct{}

func (e ResourceCreated) Info() string {
	return fmt.Sprintf(resourceCreatedMessage.message, e.ResourceName, e.ResourceContent)
}

func (e ResourceUpdated) Info() string {
	return fmt.Sprintf(resourceUpdatedMessage.message, e.ResourceName, e.ResourceContent)
}

func (e ResourceRead) Info() string {
	return fmt.Sprintf(resourceReadMessage.message, e.ResourceName, e.ResourceContent)
}

func (e ResourceDeleted) Info() string {
	return fmt.Sprintf(resourceDeletedMessage.message, e.ResourceName)
}

func (e FailedToDecodeJSON) Info() string {
	return fmt.Sprintf(failedToDecodeJSONMessage.message, e.Resource)
}

func (e FailedToEncodeJSON) Info() string {
	return fmt.Sprintf(failedToEncodeJSONMessage.message, e.Resource)
}

func (e FailedToExecuteGremlinQuery) Info() string {
	return fmt.Sprintf(failedToExecuteGremlinQueryMessage.message)
}
