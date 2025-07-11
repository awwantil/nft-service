package database

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"

	coreconfig "main/tools/pkg/core_config"
)

// NewRedisClient создает клиент подключения к redis DB
func NewRedisClient(ctx context.Context, cfg coreconfig.Redis) (*redis.Client, error) {

	// Construct the address for the Redis server
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	// Initialize Redis options with connection details
	options := &redis.Options{
		Addr:     addr,
		Username: cfg.Username,
		Password: cfg.Password,
		DB:       cfg.Database,
	}

	// Configure TLS settings if a certificate path is provided
	if cfg.CertPath != "" {
		// Read the contents of the certificate file
		caCert, err := os.ReadFile(cfg.CertPath)
		if err != nil {
			return nil, err
		}

		// Create a certificate pool and append the certificate
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		// Configure TLS settings with the certificate pool
		options.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
			RootCAs:    caCertPool,
		}
	}
	rdb := redis.NewClient(options)

	// Ping the Redis server to check the connection
	if pong := rdb.Ping(ctx); pong.String() != "ping: PONG" {
		return rdb, fmt.Errorf("ping error: %w", pong.Err())
	}

	// Return the initialized Redis client
	return rdb, nil
}
