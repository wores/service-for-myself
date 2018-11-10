package config

var env = NewEnv()

// GetEnv 環境変数を扱うEnvを返す
func GetEnv() *Env {
	return env
}

// Env 環境変数(gae/secret.yamlに定義)
type Env struct {
	slack *Slack
}

// GetSlack プロパティのslackを返す
func (env *Env) GetSlack() *Slack {
	return env.slack
}

// NewEnv 環境変数用のインスタンスを生成する
func NewEnv() *Env {
	env := new(Env)
	env.slack = NewSlack()
	return env
}
