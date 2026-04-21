package liey

import (
	"cnb.cool/liey/liey-go/auth"
	"cnb.cool/liey/liey-go/contracts"
	"cnb.cool/liey/liey-go/db"
	"cnb.cool/liey/liey-go/event"
	"cnb.cool/liey/liey-go/server"
	"cnb.cool/liey/liey-go/storage"
	gocap "github.com/ackcoder/go-cap"
	"github.com/spf13/cobra"
)

var (
	rocket *Rocket
)

type Rocket struct {
	rootCommand *cobra.Command

	server *server.EchoServer

	logger contracts.Logger

	id contracts.IdGenerator

	auth *auth.Manager

	captcha *gocap.Cap

	storage *storage.Manager

	bus *event.Bus[any]

	db *db.DB
}

func New() *Rocket {
	if rocket != nil {
		return rocket
	}

	rocket = &Rocket{}

	rocket.createRootCommand()

	return rocket
}
