package main

import "AuthInGo/app"

func main() {
	cfg := app.Config{
		Addr: ":3003",
	}

	app := app.Application{
		Config: cfg,
	}

	app.Run()
}
