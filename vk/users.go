package vk

import (
	"encoding/json"
	"errors"
	"github.com/valyala/fasthttp"
	"strconv"
	"strings"
)

type (
	UsersGetRs struct {
		Users []User `json:"response"`
	}
	User struct {
		Id              int    `json:"id"`
		FirstName       string `json:"first_name"`
		LastName        string `json:"last_name"`
		CanAccessClosed bool   `json:"can_access_closed"`
		IsClosed        bool   `json:"is_closed"`
	}
)

func UsersGet(userIds []int, fields []string, nameCase string) (*UsersGetRs, error) {
	var errRs ErrorRs

	rq := fasthttp.AcquireRequest()
	rs := fasthttp.AcquireResponse()

	defer func() {
		fasthttp.ReleaseRequest(rq)
		fasthttp.ReleaseResponse(rs)
	}()

	prepareRequest(apiHost, ContentTypeUrlEncoded, rq)

	rq.URI().SetPath(apiUrl + "/users.get")

	var ids strings.Builder
	for i := 0; i < len(userIds); i++ {
		ids.WriteString(strconv.Itoa(userIds[i]))
		if i != len(userIds)-1 {
			ids.WriteString(",")
		}
	}
	rq.PostArgs().Add("user_ids", ids.String())

	var addFields strings.Builder
	for i := 0; i < len(fields); i++ {
		addFields.WriteString(fields[i])
		if i != len(fields)-1 {
			ids.WriteString(",")
		}
	}
	rq.PostArgs().Add("fields", addFields.String())

	rq.PostArgs().Add("name_case", nameCase)

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

	var result UsersGetRs
	err = json.Unmarshal(body, &result)

	return &result, nil
}
