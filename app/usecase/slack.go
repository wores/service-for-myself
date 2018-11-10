package usecase

import (
	"context"

	"github.com/nlopes/slack/slackevents"
	"github.com/wores/service-for-myself/app/config"
	"github.com/wores/service-for-myself/app/infrastructure/repository"
	"github.com/wores/service-for-myself/app/usecase/input"
	"github.com/wores/service-for-myself/app/usecase/output"
	"google.golang.org/appengine/log"
)

type slackUsecase struct {
	env             *config.Env
	slackRepository repository.SlackRepository
	ocrRepository   repository.OcrRepository
}

// SlackUsecase ユースケース層のインターフェイス
type SlackUsecase interface {
	DetectAndPostTextFromImageURL(ctx context.Context, slackAPIEvent *input.InputSlackEvent) error
}

// NewSlackUsecase slackUsecaseのインスタンスを生成して返す
func NewSlackUsecase(
	env *config.Env,
	slackRepository repository.SlackRepository,
	ocrRepository repository.OcrRepository,
) SlackUsecase {
	return &slackUsecase{env, slackRepository, ocrRepository}
}

// DetectTextFromImageURL URLの画像からテキストを検出する
func (slackUsecase *slackUsecase) DetectAndPostTextFromImageURL(
	ctx context.Context,
	slackAPIEvent *input.InputSlackEvent,
) error {
	event := slackAPIEvent.GetEventsAPIEvent()
	if event.Type == slackevents.CallbackEvent {
		log.Debugf(ctx, "CallbackEvent以外は何もしない event type: %s", event.Type)
		return nil
	}

	innerEvent := event.InnerEvent
	switch ev := innerEvent.Data.(type) {
	case *slackevents.MessageEvent:
		if ev.SubType == "bot_message" {
			log.Debugf(ctx, "botからのメッセージは処理しない")
			return nil
		}

		if ev.SubType != "file_share" {
			log.Debugf(ctx, "ファイル以外の投稿は何もしない sub type: %s", ev.SubType)
			slackUsecase.slackRepository.PostMessageToSpecifiedChannel(
				ctx,
				ev.Channel,
				"画像を投稿したらOCRしてあげるよ",
			)
			return nil
		}

		log.Infof(ctx, "file content = %#v", ev.Files)
		// 画像からテキストを抽出するAPIを実行する
		detectedText, err := slackUsecase.ocrRepository.DetectTextFromImage(
			ctx,
			ev.Files[0].URLPrivate,
			slackUsecase.env.GetSlack().GetAccessTokenWithBearer(),
		)
		if err != nil {
			log.Errorf(ctx, "error: %#v", err)
			return output.ErrOcr
		}

		// 画像から抽出したテキストをSlackへ投稿する
		err = slackUsecase.slackRepository.PostMessageToSpecifiedChannel(
			ctx,
			ev.Channel,
			detectedText,
		)
		if err != nil {
			log.Errorf(ctx, "error: %#v", err)
			return output.ErrPostMessage
		}
		break
	default:
		log.Warningf(ctx, "想定してないのが来たかも event type: #v", innerEvent.Data)
		break
	}

	return nil
}
