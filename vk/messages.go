package vk

import (
	"encoding/json"
	"errors"
	"github.com/valyala/fasthttp"
	"log"
	"strconv"
)

type (
	MessageLongPoll struct {
		Id        int64  `json:"messageId"`
		Flags     int64  `json:"flags"`
		PeerId    int64  `json:"peerId"`
		Ts        int64  `json:"ts"`
		Text      string `json:"text"`
		Title     string `json:"title"`
		RepliedId int64  `json:"repliedId"`
	}

	MessagesSendRs struct {
	}
)

func NewMessageLongPoll(event []interface{}) (*MessageLongPoll, error) {
	msgId, _ := event[1].(json.Number).Int64()
	msgFlags, _ := event[2].(json.Number).Int64()
	msgPeerId, _ := event[3].(json.Number).Int64()
	msgTs, _ := event[4].(json.Number).Int64()
	msgText := event[5].(string)
	msgTitle := event[6].(map[string]interface{})
	msgReply := event[7].(map[string]interface{})

	res := new(MessageLongPoll)

	res.Id = msgId
	res.Flags = msgFlags
	res.PeerId = msgPeerId
	res.Ts = msgTs
	res.Text = msgText
	res.Title = msgTitle["title"].(string)

	if msgReply["reply"] != nil {
		var convMsg struct {
			Id int64 `json:"conversation_message_id"`
		}
		_ = json.Unmarshal([]byte(msgReply["reply"].(string)), &convMsg)
		res.RepliedId = convMsg.Id
	}

	return res, nil
}

func MessagesSend(peerId int64, message string, attachment *Document) (*MessagesSendRs, error) {
	rq := fasthttp.AcquireRequest()
	rs := fasthttp.AcquireResponse()

	defer func() {
		fasthttp.ReleaseRequest(rq)
		fasthttp.ReleaseResponse(rs)
	}()

	prepareRequest(apiHost, ContentTypeUrlEncoded, rq)

	rq.URI().SetPath(apiUrl + "/messages.send")

	rq.PostArgs().Add("peer_id", strconv.FormatInt(peerId, 10))

	rq.PostArgs().Add("random_id", strconv.FormatInt(Random.Int63(), 10))

	if attachment != nil {
		rq.PostArgs().Add("attachment", "doc"+strconv.FormatInt(attachment.Doc.OwnerId, 10)+"_"+strconv.FormatInt(attachment.Doc.Id, 10))
	}

	if message != "" {
		rq.PostArgs().Add("message", message)
	}

	err := apiClient.Do(rq, rs)
	if err != nil {
		return nil, err
	}
	if rs.StatusCode() != 200 {
		return nil, errors.New("status code " + strconv.Itoa(rs.StatusCode()) + "returned")
	}

	body := rs.Body()

	log.Println(string(body))

	var errRs ErrorRs
	err = json.Unmarshal(body, &errRs)
	if err != nil {
		return nil, err
	}
	if errRs.Error() != "" {
		return nil, errRs
	}

	var result MessagesSendRs
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
