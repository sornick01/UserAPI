package main

import (
	"github.com/sornick01/UserAPI/server"
	"log"
)

func main() {
	a := server.NewApp()

	err := a.Run("3333")
	if err != nil {
		log.Fatal(err)
	}
}
