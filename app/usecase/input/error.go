package input

import "errors"

var (
	// ErrParseSlackEvent SlackEventへの変換でエラー
	ErrParseSlackEvent = errors.New("failed to parse slack event")
)
