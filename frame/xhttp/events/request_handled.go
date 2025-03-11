package events

import (
	"github.com/xframe-go/x/contracts"
	"github.com/xframe-go/x/frame/events"
)

type symbol string

type (
	RequestHandled struct {
		payload *RequestData
		*events.Event[*RequestData]
	}

	RequestData struct {
		Method string
		Path   string
	}
)

var (
	RequestHandledEvent contracts.Event[*RequestData] = &RequestHandled{
		Event: events.NewEvent[*RequestData](symbol("RequestHandled")),
	}
)
