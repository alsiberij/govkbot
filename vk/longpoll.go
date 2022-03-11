package vk

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/valyala/fasthttp"
	"log"
	"strconv"
	"strings"
)

type (
	LongPollServerRs struct {
		Content struct {
			Server string `json:"server"`
			Key    string `json:"key"`
			Ts     int    `json:"ts"`
		} `json:"response"`
	}
	LongPollRs struct {
		Ts      int             `json:"ts"`
		Updates [][]interface{} `json:"updates"`
	}
)

const (
	EventNewMessage  = 4
	EventUserOnline  = 8
	EventUserOffline = 9
)

var (
	NewMsgLongPollHandler func(messageId, flags, peerId, ts int64, text string)
	UserOnlineHandler     func(userId, ts int64, isOnline bool)
)

func GetLongPollServer() (*LongPollServerRs, error) {
	var errRs ErrorRs

	rq := fasthttp.AcquireRequest()
	rs := fasthttp.AcquireResponse()

	defer func() {
		fasthttp.ReleaseRequest(rq)
		fasthttp.ReleaseResponse(rs)
	}()

	prepareRequest(apiHost, ContentTypeUrlEncoded, rq)

	rq.URI().SetPath(apiUrl + "/messages.getLongPollServer")

	rq.PostArgs().Add("lp_version", "3")

	err := apiClient.Do(rq, rs)
	if err != nil {
		return nil, err
	}
	if rs.StatusCode() != 200 {
		return nil, errors.New("status code " + strconv.Itoa(rs.StatusCode()) + "returned")
	}

	body := rs.Body()

	err = json.Unmarshal(body, &errRs)
	if err != nil {
		return nil, err
	}
	if errRs.Error() != "" {
		return nil, errRs
	}

	var result LongPollServerRs
	err = json.Unmarshal(body, &result)

	return &result, nil
}

func LongPoll(server *LongPollServerRs) {
	rq := fasthttp.AcquireRequest()
	rs := fasthttp.AcquireResponse()

	defer func() {
		fasthttp.ReleaseRequest(rq)
		fasthttp.ReleaseResponse(rs)
	}()

	prepareRequest(longPollHost, ContentTypeUrlEncoded, rq)
	rq.URI().SetPath(strings.Split(server.Content.Server, "/")[1])

	rqPostArgs := rq.PostArgs()
	rqPostArgs.Add("act", "a_check")
	rqPostArgs.Add("key", server.Content.Key)
	rqPostArgs.Add("wait", "25")
	rqPostArgs.Add("version", "3")

	var lpRs LongPollRs
	lpRs.Ts = server.Content.Ts
	for {
		rs.Reset()
		rq.URI().QueryArgs().Set("ts", strconv.Itoa(lpRs.Ts))

		err := longPollClient.Do(rq, rs)
		if err != nil {
			return
		}

		body := rs.Body()
		var bodyBuffer bytes.Buffer
		bodyBuffer.Write(body)

		dec := json.NewDecoder(&bodyBuffer)
		dec.UseNumber()

		err = dec.Decode(&lpRs)
		if err != nil {
			return
		}

		for i := range lpRs.Updates {
			updateType, _ := lpRs.Updates[i][0].(json.Number).Int64()
			log.Println(lpRs.Updates[i])
			switch updateType {
			case EventNewMessage:
				if NewMsgLongPollHandler == nil {
					continue
				}
				msgId, _ := lpRs.Updates[i][1].(json.Number).Int64()
				msgFlags, _ := lpRs.Updates[i][2].(json.Number).Int64()
				peerId, _ := lpRs.Updates[i][3].(json.Number).Int64()
				ts, _ := lpRs.Updates[i][4].(json.Number).Int64()
				text, _ := lpRs.Updates[i][5].(string)
				go NewMsgLongPollHandler(msgId, msgFlags, peerId, ts, text)
			case EventUserOnline, EventUserOffline:
				isOnline := updateType == EventUserOnline
				userId, _ := lpRs.Updates[i][1].(json.Number).Int64()
				ts, _ := lpRs.Updates[i][3].(json.Number).Int64()
				go UserOnlineHandler(-userId, ts, isOnline)
			}
		}
	}

}
