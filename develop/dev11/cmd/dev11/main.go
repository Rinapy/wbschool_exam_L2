package main

import "dev11/internal/dev11/server"

func main() {
	cfg := server.DefaultCfg()

	s := server.NewServer(cfg)
	s.Run()

}
