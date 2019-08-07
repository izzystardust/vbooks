package main

import (
	"log"
	"os"

	"git.sr.ht/~izzy/vbooks"
)

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3001"
	}

	log.Fatal(vbooks.Start(":" + port))
}
