package vk

import "encoding/json"

type (
	ErrorRs struct {
		Content struct {
			Code int    `json:"error_code"`
			Msg  string `json:"error_msg"`
		} `json:"error"`
	}
)

func (e ErrorRs) Error() string {
	if e.Content.Code == 0 {
		return ""
	}
	content, _ := json.Marshal(&e)
	return string(content)
}
