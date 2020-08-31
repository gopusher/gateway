package main

import (
	"github.com/gopusher/gateway/app/gateway/app/cmd/app"
)

func main() {
	cmd := app.NewGatewayCommand()
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
