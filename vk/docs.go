package vk

import (
	"bot/vk/logs"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/valyala/fasthttp"
	"io"
	"mime/multipart"
	"os"
	"strconv"
)

type (
	DocsGetMessageUploadServerRs struct {
		Content struct {
			Url string `json:"upload_url"`
		} `json:"response"`
	}

	DocsUploadRs struct {
		File string `json:"file"`
	}

	DocsSaveRs struct {
		Content Document `json:"response"`
	}

	Document struct {
		Type string `json:"type"`
		Doc  struct {
			Id      int64  `json:"id"`
			OwnerId int64  `json:"owner_id"`
			Title   string `json:"title"`
			Size    int64  `json:"size"`
			Ext     string `json:"ext"`
			Ts      int64  `json:"date"`
			Type    int64  `json:"type"`
			Url     string `json:"url"`
			//todo Preview(gif)
			IsLicenced int64 `json:"isLicenced"`
		} `json:"doc"`

		/*
			AudioMessage struct {

			} `json:"audio_msg"`
		*/
	}
)

func DocGetMessageUploadServer(docType string, peerId int64, isGroupChat bool) (*DocsGetMessageUploadServerRs, error) {
	rq := fasthttp.AcquireRequest()
	rs := fasthttp.AcquireResponse()

	defer func() {
		fasthttp.ReleaseRequest(rq)
		fasthttp.ReleaseResponse(rs)
	}()

	prepareRequest(apiHost, ContentTypeUrlEncoded, rq)

	rq.URI().SetPath(apiUrl + "/docs.getMessagesUploadServer")

	rq.PostArgs().Add("docType", docType)

	peer := peerId
	if isGroupChat {
		peer += PeerMinId
	}
	rq.PostArgs().Add("peer_id", strconv.FormatInt(peer, 10))

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

	var result DocsGetMessageUploadServerRs
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func DocsUploadToMessageServer(server *DocsGetMessageUploadServerRs, file string) (*DocsUploadRs, error) {
	rq := fasthttp.AcquireRequest()
	rs := fasthttp.AcquireResponse()

	defer func() {
		fasthttp.ReleaseRequest(rq)
		fasthttp.ReleaseResponse(rs)
	}()

	rq.SetRequestURI(server.Content.Url)
	rq.URI().SetScheme("https")
	rq.Header.SetMethod("POST")

	bodyBuff := &bytes.Buffer{}
	w := multipart.NewWriter(bodyBuff)
	fw, err := w.CreateFormFile("file", file)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(fw, f)
	if err != nil {
		return nil, err
	}

	_ = f.Close()
	_ = w.Close()

	rq.SetBody(bodyBuff.Bytes())
	rq.Header.SetContentType(w.FormDataContentType())

	err = docsMessagesUploadClient.Do(rq, rs)
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

	var result DocsUploadRs
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func DocsSave(file *DocsUploadRs, title string) (*DocsSaveRs, error) {
	rq := fasthttp.AcquireRequest()
	rs := fasthttp.AcquireResponse()

	defer func() {
		fasthttp.ReleaseRequest(rq)
		fasthttp.ReleaseResponse(rs)
	}()

	prepareRequest(apiHost, ContentTypeUrlEncoded, rq)

	rq.URI().SetPath(apiUrl + "/docs.save")

	rq.PostArgs().Add("file", file.File)
	rq.PostArgs().Add("title", title)

	err := apiClient.Do(rq, rs)
	if err != nil {
		return nil, err
	}
	if rs.StatusCode() != 200 {
		return nil, errors.New("status code " + strconv.Itoa(rs.StatusCode()) + "returned")
	}

	body := rs.Body()
	logs.Log(body)

	var errRs ErrorRs
	err = json.Unmarshal(body, &errRs)
	if err != nil {
		return nil, err
	}
	if errRs.Error() != "" {
		return nil, errRs
	}

	var result DocsSaveRs
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
