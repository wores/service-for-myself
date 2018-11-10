package config

import "os"

// Slack Slack用の環境変数
type Slack struct {
	accessToken       string
	verificationToken string
}

func (slack *Slack) GetAccessToken() string {
	return slack.accessToken
}

// GetAccessToken AccessTokenをBearer 付きで返す
func (slack *Slack) GetAccessTokenWithBearer() string {
	return "Bearer " + slack.accessToken
}

// GetVerificationToken VerificationTokenを返す
func (slack *Slack) GetVerificationToken() string {
	return slack.verificationToken
}

// NewSlack 環境変数Slack用のインスタンスを生成する
func NewSlack() *Slack {
	slack := new(Slack)
	slack.accessToken = os.Getenv("SLACK_TOKEN")
	slack.verificationToken = os.Getenv("VERIFICATION_TOKEN")
	return slack
}
