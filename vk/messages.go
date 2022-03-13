package vk

import (
	"encoding/json"
	"errors"
	"github.com/valyala/fasthttp"
	"log"
	"strconv"
	"strings"
)

type (
	MessagesSendRs struct {
		NewMessageId int64 `json:"response"`
	}

	MessagesEditRs struct {
	}
)

func MessagesSend(peerId int64, message string, attachment *Document, replyTo int64) (*MessagesSendRs, error) {
	rq := fasthttp.AcquireRequest()
	rs := fasthttp.AcquireResponse()

	defer func() {
		fasthttp.ReleaseRequest(rq)
		fasthttp.ReleaseResponse(rs)
	}()

	prepareRequest(apiHost, ContentTypeUrlEncoded, rq)

	rq.URI().SetPath(apiUrl + "/messages.send")

	rq.PostArgs().Add("peer_id", strconv.FormatInt(peerId, 10))
	rq.PostArgs().Add("random_id", strconv.FormatInt(int64(Random.Int31()), 10))

	//todo not doc only
	if attachment != nil {
		rq.PostArgs().Add("attachment", "doc"+strconv.FormatInt(attachment.Doc.OwnerId, 10)+"_"+strconv.FormatInt(attachment.Doc.Id, 10))
	}

	if message != "" {
		rq.PostArgs().Add("message", message)
	}

	if replyTo != 0 {
		rq.PostArgs().Add("reply_to", strconv.FormatInt(replyTo, 10))
	}

	err := apiClient.Do(rq, rs)
	if err != nil {
		return nil, err
	}
	if rs.StatusCode() != 200 {
		return nil, errors.New("status code " + strconv.Itoa(rs.StatusCode()) + "returned")
	}

	body := rs.Body()

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

func MessagesEdit(messageId, peerId int64, message string) error {
	rq := fasthttp.AcquireRequest()
	rs := fasthttp.AcquireResponse()

	defer func() {
		fasthttp.ReleaseRequest(rq)
		fasthttp.ReleaseResponse(rs)
	}()

	prepareRequest(apiHost, ContentTypeUrlEncoded, rq)

	rq.URI().SetPath(apiUrl + "/messages.edit")

	rq.PostArgs().Add("message_id", strconv.FormatInt(messageId, 10))
	rq.PostArgs().Add("peer_id", strconv.FormatInt(peerId, 10))
	rq.PostArgs().Add("random_id", strconv.FormatInt(int64(Random.Int31()), 10))

	if message != "" {
		rq.PostArgs().Add("message", message)
	}

	err := apiClient.Do(rq, rs)
	if err != nil {
		return err
	}
	if rs.StatusCode() != 200 {
		return errors.New("status code " + strconv.Itoa(rs.StatusCode()) + "returned")
	}

	body := rs.Body()

	log.Println(string(body))

	var errRs ErrorRs
	err = json.Unmarshal(body, &errRs)
	if err != nil {
		return err
	}
	if errRs.Error() != "" {
		return errRs
	}

	var result MessagesEditRs
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}

	return nil
}

func MessagesDelete(messageIds []int64, deleteForAll bool) error {
	rq := fasthttp.AcquireRequest()
	rs := fasthttp.AcquireResponse()

	defer func() {
		fasthttp.ReleaseRequest(rq)
		fasthttp.ReleaseResponse(rs)
	}()

	prepareRequest(apiHost, ContentTypeUrlEncoded, rq)

	rq.URI().SetPath(apiUrl + "/messages.delete")

	var ids strings.Builder
	for i := 0; i < len(messageIds); i++ {
		ids.WriteString(strconv.FormatInt(messageIds[i], 10))
		if i != len(messageIds)-1 {
			ids.WriteString(",")
		}
	}
	rq.PostArgs().Add("message_ids", ids.String())

	if deleteForAll {
		rq.PostArgs().Add("delete_for_all", "1")
	}

	err := apiClient.Do(rq, rs)
	if err != nil {
		return err
	}
	if rs.StatusCode() != 200 {
		return errors.New("status code " + strconv.Itoa(rs.StatusCode()) + "returned")
	}

	body := rs.Body()

	var errRs ErrorRs
	err = json.Unmarshal(body, &errRs)
	if err != nil {
		return err
	}
	if errRs.Error() != "" {
		return errRs
	}

	var result MessagesEditRs
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}

	return nil
}
