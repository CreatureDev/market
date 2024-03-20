package main

import (
	"github.com/CreatureDev/market/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var publisherCmd = &cobra.Command{
	Use:   "publisher",
	Short: "product publisher for the XRPL Market",
}

func run(conf config.PublisherConfig) {

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
