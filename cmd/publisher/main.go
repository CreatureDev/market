package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	marketv1 "github.com/CreatureDev/market/gen/go/market/v1"
	"github.com/CreatureDev/market/pkg/config"
	"github.com/CreatureDev/market/pkg/service"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var publisherCmd = &cobra.Command{
	Use:   "publisher",
	Short: "product publisher for the XRPL Market",
}

func run(conf config.PublisherConfig) {
	publisherCtx := context.Background()
	mux := runtime.NewServeMux()
	if err := marketv1.RegisterPublisherServiceHandlerServer(publisherCtx, mux, service.NewPublisherService(conf)); err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	errChan := make(chan error)
	stopChan := make(chan os.Signal, 2)

	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)

	select {
	case e := <-errChan:
		fmt.Println(e.Error())
		os.Exit(-1)
	case s := <-stopChan:
		fmt.Println("Interrupt ", s)
		return
	}

}

func main() {
	var conf config.PublisherConfig
	publisherCmd.PersistentFlags().StringVar(&conf.Secret, "secret", "", "secret for NFT issuing account")
	viper.BindPFlag("secret", publisherCmd.Flags().Lookup("secret"))
	publisherCmd.PersistentFlags().StringVar(&conf.SecretPath, "secret_path", "~/.xrpl/secret.ed25519", "Path to secret for NFT issueing account")
	viper.BindPFlag("secret_path", publisherCmd.Flags().Lookup("secret_path"))
	publisherCmd.PersistentFlags().StringVar(&conf.MarketURL, "market_url", "", "URL of the market to publish to")
	viper.BindPFlag("market_url", publisherCmd.Flags().Lookup("market_url"))

	publisherCmd.Run = func(_ *cobra.Command, _ []string) {
		run(conf)
	}
}
