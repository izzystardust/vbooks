package main

import (
	"log"

	"git.sr.ht/~izzy/vbooks/server"
)

func main() {
	log.Fatal(server.Start(":3001"))
}
