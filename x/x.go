package x

import (
	gocap "github.com/ackcoder/go-cap"
	"github.com/spf13/cobra"
	"github.com/xframe-go/x/auth"
	"github.com/xframe-go/x/contracts"
	"github.com/xframe-go/x/db"
	"github.com/xframe-go/x/event"
	"github.com/xframe-go/x/server"
	"github.com/xframe-go/x/storage"
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
