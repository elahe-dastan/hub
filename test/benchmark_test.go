package test

import (
	"testing"

	"github.com/elahe-dastan/applifier/config"
	"github.com/elahe-dastan/applifier/internal/client"
	"github.com/elahe-dastan/applifier/internal/server"
	"github.com/stretchr/testify/require"
)

const clientCount = 100
const benchmarkServerPort = "8080"

func TestBenchmark(t *testing.T) {
	srv := server.New()
	cfg := config.Config{Address: ":" + benchmarkServerPort}

	go func() {
		require.NoError(t, srv.Start(cfg))
	}()

	for i := 0; i < clientCount; i++ {
		cli := client.New()

		go func() {
			require.NoError(t, cli.Connect("127.0.0.1:"+benchmarkServerPort))
		}()
	}
}
