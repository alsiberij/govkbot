package vk

import (
	"errors"
	"regexp"
)

var (
	accessToken string

	isValidToken = regexp.MustCompile("[0-9a-f]{85}").MatchString
)

func Auth(token string) error {
	if !isValidToken(token) {
		return errors.New("invalid token")
	}
	accessToken = token
	return nil
}
