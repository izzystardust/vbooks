package main

import (
	"log"

	"git.sr.ht/~izzy/vbooks"
)

func main() {
	log.Fatal(vbooks.Start(":3001"))
}
