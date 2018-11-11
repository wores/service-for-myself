package infrastructure

import (
	"context"
	"fmt"
	"io"
	"net/http"

	vision "cloud.google.com/go/vision/apiv1"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
	pb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

// OcrAPI OCR用APIの構造体
type OcrAPI struct{}

// NewOcrAPI OcrAPIのインスタンスを生成して返す
func NewOcrAPI() *OcrAPI {
	return new(OcrAPI)
}

// DetectTextFromImage URLの画像からOCRを実行する
func (ocrAPI *OcrAPI) DetectTextFromImage(ctx context.Context, url string, authHeaderValue string) (text string, err error) {
	text = "none"
	// urlから画像取得
	imageBin, fetchErr := ocrAPI.fetchImageFromURL(ctx, url, authHeaderValue)
	if fetchErr != nil {
		return text, fetchErr
	}
	defer imageBin.Close()

	// vision apiに投げてOCRしてもらう
	annotations, detectErr := ocrAPI.postImageToCloudVisionAPI(ctx, imageBin)
	if detectErr != nil {
		err := fmt.Errorf("failed to detect. %#v", detectErr)
		log.Errorf(ctx, err.Error())
		return text, err
	}

	if len(annotations) == 0 {
		text = "解析してみたけど、テキストが見つからなかったよ"
		log.Debugf(ctx, text)
		return text, nil
	}

	log.Debugf(ctx, "Text:")
	text = ""
	for _, annotation := range annotations {
		log.Debugf(ctx, "an = %#v", annotation.BoundingPoly)
	}
	text = "OCR成功: \n```\n" + annotations[0].Description + "\n```"
	log.Debugf(ctx, text)
	return text, nil
}

// postImageToCloudVisionAPI vision apiに検出してもらう画像を投げる
func (ocrAPI *OcrAPI) postImageToCloudVisionAPI(ctx context.Context, imageBin io.ReadCloser) ([]*pb.EntityAnnotation, error) {
	// vision api用のclientを作成
	client, createErr := vision.NewImageAnnotatorClient(ctx)
	if createErr != nil {
		return nil, createErr
	}

	// vision apiに渡す画像を生成
	image, readerErr := vision.NewImageFromReader(imageBin)
	if readerErr != nil {
		return nil, readerErr
	}

	// vision apiにOCRしてもらう
	annotations, detectErr := client.DetectTexts(ctx, image, nil, 10)
	if detectErr != nil {
		return nil, detectErr
	}

	return annotations, nil
}

// URLから画像を取得する
func (ocrAPI *OcrAPI) fetchImageFromURL(ctx context.Context, url string, authHeaderValue string) (image io.ReadCloser, err error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", authHeaderValue)
	client := urlfetch.Client(ctx)
	res, resErr := client.Do(req)
	if resErr != nil {
		err := fmt.Errorf("failed to request image. %#v", resErr)
		log.Errorf(ctx, err.Error())
		return nil, err
	}

	return res.Body, nil
}
