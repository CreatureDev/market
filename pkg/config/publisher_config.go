package config

type PublisherConfig struct {
	Secret     string `json:"secret" yaml:"secret"`
	SecretPath string `json:"secret_path" yaml:"secret_path"`
	MarketURL  string `json:"market_url" yaml:"market_url"`
}
