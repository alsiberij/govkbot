package vk

import (
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
	RsDocsGetMessageUploadServer struct {
		Content struct {
			Url string `json:"upload_url"`
		} `json:"response"`
	}

	RsDocsUpload struct {
		File string `json:"file"`
	}

	RsDocsSave struct {
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
	}
)

func DocGetMessageUploadServer(docType string, peerId int64, isGroupChat bool) (RsDocsGetMessageUploadServer, error) {
	var result RsDocsGetMessageUploadServer

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

func DocsUploadToMessageServer(server *RsDocsGetMessageUploadServer, file string) (RsDocsUpload, error) {
	var result RsDocsUpload

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
		return result, err
	}

	f, err := os.Open(file)
	if err != nil {
		return result, err
	}

	_, err = io.Copy(fw, f)
	if err != nil {
		return result, err
	}

	_ = f.Close()
	_ = w.Close()

	rq.SetBody(bodyBuff.Bytes())
	rq.Header.SetContentType(w.FormDataContentType())

	err = docsMessagesUploadClient.Do(rq, rs)
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

func DocsSave(file *RsDocsUpload, title string) (RsDocsSave, error) {
	var result RsDocsSave

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
