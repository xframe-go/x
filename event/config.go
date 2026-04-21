package event

import "github.com/xframe-go/x/contracts"

type Config struct {
	Driver   contracts.EventDriver
	Handlers []Handler
}
