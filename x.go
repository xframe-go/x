package x

import (
	"github.com/xframe-go/x/frame/app"
	"github.com/xframe-go/x/utils/singleton"
)

var xApp = singleton.New[*app.Application](app.New)

func App() *app.Application {
	return xApp.Get()
}
