package x

import (
	"github.com/spf13/cobra"
	"github.com/xframe-go/x/contracts"
)

func createXServerCmd(app contracts.Application) *cobra.Command {
	return &cobra.Command{
		Use:   "serve [name]",
		Short: "Serve the application",
		Run: func(cmd *cobra.Command, args []string) {
			name := "default"
			if len(args) > 0 {
				name = args[0]
			}

			if err := app.Server(name).Start(); err != nil {
				app.Log().Fatal(err.Error())
			}
		},
	}
}
