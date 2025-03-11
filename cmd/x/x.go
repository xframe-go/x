package x

import (
	"github.com/spf13/cobra"
	"github.com/xframe-go/x/cmd/x/makes"
	"github.com/xframe-go/x/contracts"
)

type Command struct {
	root cobra.Command
}

func New(app contracts.Application) *Command {
	cmd := &Command{
		root: cobra.Command{
			Use: "x",
			RunE: func(cmd *cobra.Command, args []string) error {
				return cmd.Help()
			},
		},
	}

	var (
		xSrv = createXServerCmd(app)
	)

	cmd.root.AddCommand(xSrv)
	cmd.root.AddGroup(&cobra.Group{
		ID:    "make",
		Title: "Make",
	})
	cmd.root.AddCommand(makes.CreateMakeGroup(app)...)

	var envFile string
	cmd.root.PersistentFlags().StringVar(&envFile, "env", ".env", "env file")

	return cmd
}

func (cmd *Command) Execute() error {
	return cmd.root.Execute()
}
