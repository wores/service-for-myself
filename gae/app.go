package main

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/wores/gae-wrs-sample/app/config"
	"github.com/wores/gae-wrs-sample/app/infrastructure/router"
	"github.com/wores/gae-wrs-sample/app/registry"
	"google.golang.org/appengine"
)

func main() {
	env := config.GetEnv()
	slackHandler := registry.NewSlackHander(env)
	router.NewRouter(slackHandler)

	http.HandleFunc("/error", ErrorHandler)

	appengine.Main()

	// ロジックとは無関係
	// github.com/gorilla/websocketは使用していないが
	// github.com/nlopes/slack内で使用している依存関係のため
	if websocket.TextMessage == 2 {
		return
	}
}
