package main

import (
	"context"
	"sync"

	"backend-challenge/config"
	"backend-challenge/db"
	"backend-challenge/api"
	"backend-challenge/user"
)

func main() {
	wg := new(sync.WaitGroup)
	defer wg.Wait()

	mainCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := config.LoadConfig(); err != nil {
		panic(err)
	}
	db.GetConn(mainCtx)
	WebApp := api.GetWebApi(mainCtx)
	user.SetUpRoutes(WebApp.GetApp())

	// WebApp.GetApp().Get("/api/*", swagger.HandlerDefault)
	WebApp.Start(mainCtx,wg)
}
