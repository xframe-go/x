package makes

import (
	"github.com/spf13/cobra"
	"github.com/xframe-go/x/contracts"
)

func createMakeControllerCmd(app contracts.Application) *cobra.Command {
	cmd := &cobra.Command{
		Use: "make:controller",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	return cmd
}
