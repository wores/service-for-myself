package infrastructure

import (
	"context"
	"errors"

	"github.com/nlopes/slack"
	"github.com/wores/service-for-myself/app/config"

	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

// Slack Slack関連の制御を行う構造体
type SlackAPI struct {
	env *config.Slack
}

// New Slackのインスタンスを生成して返す
func NewSlackAPI(env *config.Slack) *SlackAPI {
	return &SlackAPI{env}
}

// PostMessageToSpecifiedChannel 指定したチャンネルへメッセージを投稿する
func (s *SlackAPI) PostMessageToSpecifiedChannel(
	ctx context.Context,
	channel string,
	text string,
) error {
	slack.SetHTTPClient(urlfetch.Client(ctx))
	api := slack.New(s.env.GetAccessToken())
	postParams := slack.PostMessageParameters{}
	postParams.Attachments = []slack.Attachment{slack.Attachment{Text: text}}
	ch, _, postErr := api.PostMessage(channel, "", postParams)
	if postErr != nil {
		log.Errorf(ctx, "post err = %#v", postErr.Error())
		return errors.New("failed to post message")
	}
	log.Debugf(ctx, "post to ch = %#v", ch)
	return nil
}
