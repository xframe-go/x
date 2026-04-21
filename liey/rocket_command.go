package liey

import "github.com/spf13/cobra"

func (r *Rocket) createRootCommand() {
	r.rootCommand = &cobra.Command{
		Use: "rocket",
	}
}

func (r *Rocket) AddCommand(commands ...*cobra.Command) *Rocket {
	r.rootCommand.AddCommand(commands...)
	return r
}

func (r *Rocket) Start() {
	if err := r.rootCommand.Execute(); err != nil {
		return
	}
}
