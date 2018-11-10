package router

import (
	"net/http"

	"github.com/wores/service-for-myself/app/handler"
)

func NewRouter(slackHandler handler.SlackAPIHandler) {
	http.HandleFunc("/", slackHandler.Ocr)
}
