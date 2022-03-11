package vk

import (
	"encoding/json"
	"errors"
	"github.com/valyala/fasthttp"
	"log"
	"strconv"
)

type (
	MessagesSendRs struct {
	}
)

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
