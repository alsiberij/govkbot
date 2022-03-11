package vk

import (
	"errors"
	"github.com/valyala/fasthttp"
	"math/rand"
	"strconv"
	"time"
)

const (
	apiHost               = "api.vk.com"
	longPollHost          = "im.vk.com"
	docsMessageUploadHost = "pu.vk.com"
	apiUrl                = "/method"

	ContentTypeUrlEncoded    = "application/x-www-form-urlencoded"
	ContentMultipartFormData = "multipart/form-data"

	PeerMinId = 2000000000
)

var (
	Random = rand.New(rand.NewSource(time.Now().UnixNano()))

	apiClient                = fasthttp.HostClient{Addr: apiHost, IsTLS: true}
	longPollClient           = fasthttp.HostClient{Addr: longPollHost, IsTLS: true}
	docsMessagesUploadClient = fasthttp.HostClient{Addr: docsMessageUploadHost, IsTLS: true}
)

func SetVersion(apiVersion string) error {
	_, err := strconv.ParseFloat(version, 32)
	if err != nil {
		return errors.New("invalid api version")
	}
	version = apiVersion
	return nil
}

func prepareRequest(host, contentType string, rq *fasthttp.Request) {
	rq.Header.SetContentType(contentType)
	rq.Header.SetMethod("POST")
	rq.URI().QueryArgs().Add("access_token", accessToken)
	rq.URI().QueryArgs().Add("v", version)
	rq.URI().SetScheme("https")
	rq.URI().SetHost(host)
}
