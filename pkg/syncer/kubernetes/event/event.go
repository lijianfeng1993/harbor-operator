package event

import (
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
)

// Event the client that push event to kubernetes
type Event interface {
	NewEventAdd(object runtime.Object, reason, message string)
}

// EventOption is the Event client interface implementation that using API calls to kubernetes.
type EventOption struct {
	eventsCli record.EventRecorder
}

// NewEvent returns a new Event client
func NewEvent(eventCli record.EventRecorder) Event {
	return &EventOption{
		eventsCli: eventCli,
	}
}

// NewEventAdd implement the Event.Interface
func (e *EventOption) NewEventAdd(object runtime.Object, reason, message string) {
	e.eventsCli.Event(object, v1.EventTypeNormal, reason, message)
}
