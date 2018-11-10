package repository

import "context"

// OcrRepository OcrAPI用のインターフェイ
type OcrRepository interface {
	DetectTextFromImage(ctx context.Context, url string, authHeaderValue string) (text string, err error)
}
