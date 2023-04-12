package main

import "07-simple-rest-api/routers"

func main() {
	var PORT = ":8080"

	routers.StartServer().Run(PORT)
}
