package vk

import "encoding/json"

type (
	NewMessageLongPollEvent struct {
		Id          int64    `json:"messageId"`
		Flags       int64    `json:"flags"`
		PeerId      int64    `json:"peerId"`
		Ts          int64    `json:"ts"`
		Text        string   `json:"text"`
		Title       string   `json:"title"`
		RepliedId   int64    `json:"repliedId"`
		Attachments []string `json:"attachments"`
	}
)

func NewMessageLongPoll(event []interface{}) (*NewMessageLongPollEvent, error) {
	msgId, _ := event[1].(json.Number).Int64()
	msgFlags, _ := event[2].(json.Number).Int64()
	msgPeerId, _ := event[3].(json.Number).Int64()
	msgTs, _ := event[4].(json.Number).Int64()
	msgText := event[5].(string)
	msgTitle := event[6].(map[string]interface{})
	msgAttach := event[7].(map[string]interface{})

	res := new(NewMessageLongPollEvent)

	res.Id = msgId
	res.Flags = msgFlags
	res.PeerId = msgPeerId
	res.Ts = msgTs
	res.Text = msgText

	title, ok := msgTitle["title"]
	if ok {
		res.Title = title.(string)
	}

	replyJson, ok := msgAttach["reply"]
	if ok {
		var convMsg struct {
			Id int64 `json:"conversation_message_id"`
		}
		_ = json.Unmarshal([]byte(replyJson.(string)), &convMsg)
		res.RepliedId = convMsg.Id
	}

	attach1, ok := msgAttach["attach1"]
	if ok {
		res.Attachments = append(res.Attachments, msgAttach["attach1_type"].(string)+attach1.(string))
	}

	return res, nil
}
