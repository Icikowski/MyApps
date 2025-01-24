package main

import (
	"log"
	"os"

	"git.sr.ht/~icikowski/myapps/cli"
)

func main() {
	err := cli.MyApps.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
