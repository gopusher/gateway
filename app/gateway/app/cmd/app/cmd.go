package app

import (
	"github.com/spf13/cobra"
)

func NewGatewayCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gateway",
		Short: "The gateway app",
		Long:  "The gateway application is a daemon program that serves all gateway requests",
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Help(); err != nil {
				panic(err)
			}
		},
	}
	cmd.AddCommand(NewStartCommand())
	return cmd
}
