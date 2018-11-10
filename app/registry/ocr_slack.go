package registry

import (
	"github.com/wores/service-for-myself/app/config"
	"github.com/wores/service-for-myself/app/handler"
	"github.com/wores/service-for-myself/app/infrastructure"
	"github.com/wores/service-for-myself/app/infrastructure/repository"
	"github.com/wores/service-for-myself/app/usecase"
)

// NewSlackHander SlackHandlerのインスタンスを生成して返す
func NewSlackHander(env *config.Env) handler.SlackAPIHandler {
	usecase := newSlackUsecase(env)
	return handler.NewSlackAPIHandler(env, usecase)
}

func newSlackUsecase(env *config.Env) usecase.SlackUsecase {
	slackRepository, ocrRepository := newSlackRepository(env)
	return usecase.NewSlackUsecase(env, slackRepository, ocrRepository)
}

func newSlackRepository(env *config.Env) (repository.SlackRepository, repository.OcrRepository) {
	slackRepository := infrastructure.NewSlackAPI(env.GetSlack())
	ocrRepository := infrastructure.NewOcrAPI()
	return slackRepository, ocrRepository
}
