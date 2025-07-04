package config

type Config struct {
	IPFS_API_URL     string `envconfig:"IPFS_API_URL" default:"1s"`
	IPFS_GATEWAY_URL string `envconfig:"IPFS_GATEWAY_URL" default:"1s"`
}
