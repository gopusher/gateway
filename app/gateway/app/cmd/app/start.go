package app

import (
	"github.com/gopusher/gateway/app/gateway/app/bootstrap"
	"github.com/spf13/cobra"
)

func NewStartCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start gateway server",
		Long:  "Start the gateway application",
		Run: func(cmd *cobra.Command, args []string) {
			cfgFile, err := cmd.Flags().GetString("config")
			if err != nil {
				panic(err)
			}

			bootstrap.Start(cfgFile, cmd.Flags())
		},
	}
	cmd.Flags().StringP("config", "c", "config.yaml", "app config file")
	return cmd
}
