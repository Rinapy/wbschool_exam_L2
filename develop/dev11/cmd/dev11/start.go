package cmd

import (
	"dev11/internal/dev11/server"
	"os"
)

func StartServer() chan os.Signal {
	cfg := server.DefaultCfg()
	s := server.NewServer(cfg)
	sigint := s.Run()
	return sigint
}
