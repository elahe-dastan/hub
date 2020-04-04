package test

import (
	"testing"

	"github.com/elahe-dastan/applifier/config"
	"github.com/elahe-dastan/applifier/internal/server"
	"github.com/stretchr/testify/require"
)

const clientCount = 100
const benchmarkServerPort = "8080"

func TestBenchmark(t *testing.T) {
	srv := server.New()
	cfg := config.ServerConfig{Address: ":" + benchmarkServerPort}
	require.NoError(t, srv.Start(cfg))
}
