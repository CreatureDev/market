package main

import (
	"fmt"

	"github.com/CreatureDev/market/pkg/app"
	"github.com/CreatureDev/market/pkg/config"

	"github.com/leaanthony/clir"
)

func main() {

	var conf config.Config
	cli := clir.NewCli("Market", "XRP Storefront for local files", "a1")

	cli.IntFlag("port", "Port for market to bind to", &conf.Port)

	cli.Action(func() error {
		return app.Run(conf)
	})

	if err := cli.Run(); err != nil {
		fmt.Println(err.Error())
	}

}
