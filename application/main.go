package main

import (
	"log"
	"os"

	"icikowski.pl/myapps/cli"
)

func main() {
	err := cli.MyApps.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
