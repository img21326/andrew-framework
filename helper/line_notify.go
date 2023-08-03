package helper

import (
	"github.com/imroc/req/v3"
	"github.com/spf13/viper"
)

var LineNotifyInstance *LineNotify

type LineNotify struct {
	token string
}

func GetLineNotify() *LineNotify {
	if LineNotifyInstance == nil {
		LineNotifyInstance = &LineNotify{
			token: viper.GetViper().GetString("LINE_NOTIFY_TOKEN"),
		}
	}
	return LineNotifyInstance
}

func (s *LineNotify) Send(msg string) error {
	_, err := req.R().
		SetHeader("Authorization", "Bearer "+s.token).
		SetFormData(map[string]string{
			"message": msg,
		}).
		Post("https://notify-api.line.me/api/notify")
	return err
}
