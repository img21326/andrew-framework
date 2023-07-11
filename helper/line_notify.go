package helper

import (
	"github.com/imroc/req/v3"
)

func LineNotify(msg string, token string) error {
	_, err := req.R().
		SetHeader("Authorization", "Bearer "+token).
		SetFormData(map[string]string{
			"message": msg,
		}).
		Post("https://notify-api.line.me/api/notify")
	return err
}
