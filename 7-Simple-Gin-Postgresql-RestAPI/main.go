package main

import (
	"book/database"
	"book/router"
)

func main() {
	database.StartDB()

	var PORT string = ":4000"
	router.StartServer().Run(PORT)
}
