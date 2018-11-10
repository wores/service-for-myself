package main

import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	errText := "エラーが発生したら GAE → Cloud PubSub → Cloud Functions → Slack へ通知が行くかのテスト"
	log.Errorf(c, errText)
	http.Error(w, errText, http.StatusInternalServerError)
}
