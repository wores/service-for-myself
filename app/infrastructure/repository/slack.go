package repository

import "context"

// SlackRepository SlackAPI用インターフェイス
type SlackRepository interface {
	PostMessageToSpecifiedChannel(ctx context.Context, channel string, text string) error
}
