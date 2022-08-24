package main

import (
	"github.com/proxx/cmd"
	"log"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatalf("%+v\n", err)
	}
}
