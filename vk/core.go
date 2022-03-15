package vk

import (
	"github.com/valyala/fasthttp"
	"math/rand"
	"time"
)

const (
	apiHost               = "api.vk.com"
	LongPollHost          = "im.vk.com"
	docsMessageUploadHost = "pu.vk.com"
	apiUrl                = "/method"
	version               = "5.131"

	ContentTypeUrlEncoded    = "application/x-www-form-urlencoded"
	ContentMultipartFormData = "multipart/form-data"

	PeerMinId = 2000000000
)

var (
	Random = rand.New(rand.NewSource(time.Now().UnixNano()))

	apiClient                = fasthttp.HostClient{Addr: apiHost, IsTLS: true}
	longPollClient           = fasthttp.HostClient{Addr: LongPollHost, IsTLS: true}
	docsMessagesUploadClient = fasthttp.HostClient{Addr: docsMessageUploadHost, IsTLS: true}
)

func prepareRequest(host, contentType string, rq *fasthttp.Request) {
	rq.Header.SetContentType(contentType)
	rq.Header.SetMethod("POST")
	rq.URI().QueryArgs().Add("access_token", accessToken)
	rq.URI().QueryArgs().Add("v", version)
	rq.URI().SetScheme("https")
	rq.URI().SetHost(host)
}
