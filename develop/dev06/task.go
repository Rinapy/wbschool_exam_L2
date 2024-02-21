package main

import (
	"dev06/cut"
	"log"
)

func main() {
	c, err := cut.NewApp()
	if err != nil {
		log.Fatal(err)
	}
	c.Run()
}
