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
	"github.com/zeromicro/go-zero/core/proc"
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

	bus *event.Bus

	db *db.DB

	shutdown []func()
}

func New() *Rocket {
	if rocket != nil {
		return rocket
	}

	rocket = &Rocket{}

	rocket.createRootCommand()

	return rocket
}

func AddShutdownListener(listener func()) {
	rocket.shutdown = append(rocket.shutdown, proc.AddShutdownListener(listener))
}

func Wait() {
	for _, waiter := range rocket.shutdown {
		waiter()
	}
}
