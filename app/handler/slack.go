package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/wores/service-for-myself/app/config"
	"github.com/wores/service-for-myself/app/usecase"
	"github.com/wores/service-for-myself/app/usecase/input"

	"github.com/nlopes/slack/slackevents"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

type slackAPIHandler struct {
	env     *config.Env
	usecase usecase.SlackUsecase
}

// SlackAPIHandler SlackAPIから叩かれるHandler
type SlackAPIHandler interface {
	Ocr(w http.ResponseWriter, r *http.Request)
}

// NewSlackAPIHandler slackAPIHandlerのインスタンスを生成する
func NewSlackAPIHandler(env *config.Env, usecase usecase.SlackUsecase) SlackAPIHandler {
	return &slackAPIHandler{env, usecase}
}

// SlackOcrAPIHandler Slackから投稿された画像をOCRで検出して、その文字列をSlackへ投稿する
func (handler *slackAPIHandler) Ocr(w http.ResponseWriter, r *http.Request) {
	envSlack := handler.env.GetSlack()
	defer r.Body.Close()
	ctx := appengine.NewContext(r)

	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	bufString := buf.String()
	log.Debugf(ctx, "body = %s", bufString)
	inputSlackEvents, parseErr := input.ParseFromStringToInputSlackEvent(
		ctx,
		bufString,
		envSlack.GetVerificationToken(),
	)
	if parseErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 初回だけ通る認可
	if inputSlackEvents.GetEventsAPIEvent().Type == slackevents.URLVerification {
		log.Debugf(ctx, "verrification")
		handler.authorize(ctx, w, bufString)
		return
	}

	if handler.usecase.DetectAndPostTextFromImageURL(ctx, inputSlackEvents) == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (handler *slackAPIHandler) authorize(ctx context.Context, w http.ResponseWriter, bufString string) {
	var challengeRes *slackevents.ChallengeResponse
	err := json.Unmarshal([]byte(bufString), &challengeRes)
	if err != nil {
		log.Errorf(ctx, "error = %#v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(challengeRes.Challenge))
}
