package main

import (
	"net/http"

	"github.com/wores/service-for-myself/app/config"
	"github.com/wores/service-for-myself/app/infrastructure/router"
	"github.com/wores/service-for-myself/app/registry"
	"google.golang.org/appengine"
)

func main() {
	env := config.GetEnv()
	slackHandler := registry.NewSlackHander(env)
	router.NewRouter(slackHandler)

	http.HandleFunc("/error", ErrorHandler)

	appengine.Main()
}
