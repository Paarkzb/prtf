package suite

import (
	"net/http"
	"os"
	"sso/internal/config"
	"testing"
	"time"
)

const (
	grpcHost = "localhost"
)

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient http.Client
}

func NewSuite(t *testing.T) *Suite {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoad(configPath())

	// ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.HTTP.Timeout)

	t.Cleanup(func() {
		t.Helper()
		// cancelCtx()
	})

	cc := http.Client{
		Transport: &http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    30 * time.Second,
			DisableCompression: true,
		},
	}

	return &Suite{
		T:          t,
		Cfg:        cfg,
		AuthClient: cc,
	}
}

func configPath() string {
	const key = "CONFIG_PATH"

	if v := os.Getenv(key); v != "" {
		return v
	}

	return "../config/config_local.yaml"
}
