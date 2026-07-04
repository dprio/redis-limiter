package main

import "github.com/dprio/redis-limiter/cmd/app"

func main() {
	app := app.New()

	if err := app.Start(); err != nil {
		println(err.Error())
	}
	println("Acabou !")
}
