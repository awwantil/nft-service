package config

import coreconfig "main/tools/pkg/core_config"

type Config struct {
	App              coreconfig.App
	Database         coreconfig.Database
	Logging          coreconfig.Logging
	Redis            coreconfig.Redis
	JWT              coreconfig.JWT
	Secret           string `envconfig:"APP_SECRET"` // Secret of the application
	IPFS_API_URL     string `envconfig:"IPFS_API_URL" default:"1s"`
	IPFS_GATEWAY_URL string `envconfig:"IPFS_GATEWAY_URL" default:"1s"`
}
