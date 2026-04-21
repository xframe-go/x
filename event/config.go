package event

import "cnb.cool/liey/liey-go/contracts"

type Config struct {
	Driver   contracts.EventDriver
	Handlers []Handler
}
