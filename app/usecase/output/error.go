package output

import "errors"

var (
	// ErrOcr Ocrの処理でエラー
	ErrOcr = errors.New("failed to detect")
	// ErrPostMessage Slackへの投稿でエラー
	ErrPostMessage = errors.New("failed to post message")
)
