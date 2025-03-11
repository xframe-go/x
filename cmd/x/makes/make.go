package makes

import (
	"github.com/spf13/cobra"
	"github.com/xframe-go/x/contracts"
)

func CreateMakeGroup(app contracts.Application) []*cobra.Command {
	return []*cobra.Command{
		createMakeModelCmd(app),
		createMakeControllerCmd(app),
	}
}
