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
	RsLongPollGetServer struct {
		Content struct {
			Server string `json:"server"`
			Key    string `json:"key"`
			Ts     int    `json:"ts"`
		} `json:"response"`
	}

	RsLongPoll struct {
		Ts      int             `json:"ts"`
		Updates [][]interface{} `json:"updates"`
	}
)

const (
	EventNewMessage = 4
)

var (
	NewMsgLongPollHandler func(msg *NewMessageLongPollEvent)
)

func GetLongPollServer() (RsLongPollGetServer, error) {
	var result RsLongPollGetServer

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
		return result, err
	}
	if rs.StatusCode() != 200 {
		return result, errors.New("ответ сервера: " + strconv.Itoa(rs.StatusCode()))
	}

	body := rs.Body()

	var rsErr RsError
	err = json.Unmarshal(body, &rsErr)
	if err != nil {
		return result, err
	}
	if rsErr.Error() != "" {
		return result, rsErr
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func LongPoll(server *RsLongPollGetServer) {
	rq := fasthttp.AcquireRequest()
	rs := fasthttp.AcquireResponse()

	defer func() {
		fasthttp.ReleaseRequest(rq)
		fasthttp.ReleaseResponse(rs)
	}()

	prepareRequest(LongPollHost, ContentTypeUrlEncoded, rq)
	rq.URI().SetPath(strings.Split(server.Content.Server, "/")[1])

	rqPostArgs := rq.PostArgs()
	rqPostArgs.Add("act", "a_check")
	rq.PostArgs().Add("mode", "74")
	rqPostArgs.Add("key", server.Content.Key)
	rqPostArgs.Add("wait", "25")
	rqPostArgs.Add("version", "3")

	var RsLp RsLongPoll
	RsLp.Ts = server.Content.Ts
	for {
		rs.Reset()
		rq.URI().QueryArgs().Set("ts", strconv.Itoa(RsLp.Ts))

		err := longPollClient.Do(rq, rs)
		if err != nil {
			return
		}

		body := rs.Body()

		var bodyBuffer bytes.Buffer
		bodyBuffer.Write(body)

		dec := json.NewDecoder(&bodyBuffer)
		dec.UseNumber()
		err = dec.Decode(&RsLp)
		if err != nil {
			return
		}

		for i := range RsLp.Updates {
			updateType, _ := RsLp.Updates[i][0].(json.Number).Int64()
			log.Println(RsLp.Updates[i])
			switch updateType {
			case EventNewMessage:
				if NewMsgLongPollHandler == nil {
					continue
				}
				msg, err := NewMessageLongPoll(RsLp.Updates[i])
				if err != nil {
					log.Println(err)
					continue
				}
				go NewMsgLongPollHandler(&msg)
			}
		}
	}

}
