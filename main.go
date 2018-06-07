package main

import (
	"math/rand"
	"os"
	"time"

	"go.dutchsec.com/lootbox/cmd"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	app := cmd.New()
	app.Run(os.Args)
}
