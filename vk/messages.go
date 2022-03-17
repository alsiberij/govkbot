package vk

import (
	"encoding/json"
	"errors"
	"github.com/valyala/fasthttp"
	"strconv"
	"strings"
)

type (
	RsMessagesSend struct {
		NewMessageId int64 `json:"response"`
	}

	RsMessagesGetFromConversation struct {
		Response struct {
			Count    int64     `json:"count"`
			Messages []Message `json:"items"`
		} `json:"response"`
	}

	Message struct {
		Id               int64  `json:"id"`
		IdInConversation int64  `json:"conversation_message_id"`
		Ts               int64  `json:"date"`
		From             int64  `json:"from"`
		Out              int64  `json:"out"`
		Important        bool   `json:"important"`
		IsHidden         bool   `json:"is_hidden"`
		PeerId           int64  `json:"peer_id"`
		RandomId         int64  `json:"random_id"`
		Text             string `json:"text"`
		//attachments[{...}]
		RepliedMessage MessageShort `json:"reply_message"`
	}

	MessageShort struct {
		Id               int64  `json:"id"`
		IdInConversation int64  `json:"conversation_message_id"`
		Ts               int64  `json:"date"`
		From             int64  `json:"from"`
		PeerId           int64  `json:"peer_id"`
		Text             string `json:"text"`
		//attachments[{...}]
	}
)

func MessagesSend(peerId int64, message string, attachments []*Document, replyTo int64) (RsMessagesSend, error) {
	var result RsMessagesSend

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

	if len(attachments) > 0 {
		var att strings.Builder
		for i := 0; i < len(attachments); i++ {
			if attachments[i] != nil {
				att.WriteString(attachments[i].Type + strconv.FormatInt(attachments[i].Doc.OwnerId, 10))
				att.WriteString("_" + strconv.FormatInt(attachments[i].Doc.Id, 10))
				if i != len(attachments)-1 {
					att.WriteString(",")
				}
			}
		}
		rq.PostArgs().Add("attachment", att.String())
	}

	if message != "" {
		rq.PostArgs().Add("message", message)
	}

	if replyTo != 0 {
		rq.PostArgs().Add("reply_to", strconv.FormatInt(replyTo, 10))
	}

	err := apiClient.Do(rq, rs)
	if err != nil {
		return result, err
	}
	if rs.StatusCode() != 200 {
		return result, errors.New("status code " + strconv.Itoa(rs.StatusCode()) + "returned")
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

func MessagesEdit(messageId, peerId int64, message string, attachment []*Document) error {
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

	var attachments strings.Builder
	for i := 0; i < len(attachment); i++ {
		attachments.WriteString(attachment[i].Type + strconv.FormatInt(attachment[i].Doc.OwnerId, 10) + "_")
		attachments.WriteString(strconv.FormatInt(attachment[i].Doc.Id, 10))
		if i != len(attachment)-1 {
			attachments.WriteString(",")
		}
	}

	rq.PostArgs().Add("attachment", attachments.String())

	err := apiClient.Do(rq, rs)
	if err != nil {
		return err
	}
	if rs.StatusCode() != 200 {
		return errors.New("status code " + strconv.Itoa(rs.StatusCode()) + "returned")
	}

	body := rs.Body()

	var errRs RsError
	err = json.Unmarshal(body, &errRs)
	if err != nil {
		return err
	}
	if errRs.Error() != "" {
		return errRs
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

	var errRs RsError
	err = json.Unmarshal(body, &errRs)
	if err != nil {
		return err
	}
	if errRs.Error() != "" {
		return errRs
	}

	return nil
}

func MessagesGetFromConversation(messageIds []int64, peerId int64) (RsMessagesGetFromConversation, error) {
	var result RsMessagesGetFromConversation

	rq := fasthttp.AcquireRequest()
	rs := fasthttp.AcquireResponse()

	defer func() {
		fasthttp.ReleaseRequest(rq)
		fasthttp.ReleaseResponse(rs)
	}()

	prepareRequest(apiHost, ContentTypeUrlEncoded, rq)

	rq.URI().SetPath(apiUrl + "/messages.getByConversationMessageId")

	rq.PostArgs().Add("peer_id", strconv.FormatInt(peerId, 10))

	var ids strings.Builder
	for i := 0; i < len(messageIds); i++ {
		ids.WriteString(strconv.FormatInt(messageIds[i], 10))
		if i != len(messageIds)-1 {
			ids.WriteString(",")
		}
	}
	rq.PostArgs().Add("conversation_message_ids", ids.String())

	err := apiClient.Do(rq, rs)
	if err != nil {
		return result, err
	}
	if rs.StatusCode() != 200 {
		return result, errors.New("status code " + strconv.Itoa(rs.StatusCode()) + "returned")
	}

	body := rs.Body()

	var errRs RsError
	err = json.Unmarshal(body, &errRs)
	if err != nil {
		return result, err
	}
	if errRs.Error() != "" {
		return result, errRs
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}
