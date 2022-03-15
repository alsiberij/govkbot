package main

import (
	_ "embed"
	"gobotvk/vk"
	"log"
)

//go:embed token.txt
var token []byte

func main() {
	err := vk.Auth(string(token))
	if err != nil {
		log.Fatal(err)
		return
	}

	longPollServer, err := vk.GetLongPollServer()
	if err != nil {
		log.Fatal(err)
		return
	}

	vk.NewMsgLongPollHandler = vkMessageHandler

	vk.LongPoll(&longPollServer)
}
