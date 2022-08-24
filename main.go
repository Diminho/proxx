package main

import (
	"log"

	"github.com/proxx/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatalf("%+v\n", err)
	}
}
