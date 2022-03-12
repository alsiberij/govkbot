package main

import (
	"bot/game"
	"bot/vk"
	_ "embed"
	"fmt"
	"image/color"
	"log"
	"strconv"
	"strings"
	"time"
)

//go:embed token.txt
var token []byte

func vkMessageRouter(msg *vk.MessageLongPoll) {
	if !(msg != nil && msg.Text != "") {
		return
	}

	if msg.Text[0] == '~' {
		switch true {

		//~gen-life-gif/160/190/10/500/100/R0
		//todo not send msg but edit existing
		case strings.Contains(msg.Text, "gen-life-gif"):
			params := strings.Split(msg.Text, "/")
			fmt.Println(params)
			if len(params) != 7 {
				_, err := vk.MessagesSend(msg.PeerId, "Invalid params set. Try ~gen-life-gif/width/height/cell/generations/threads", nil)
				if err != nil {
					log.Println("UNABLE TO SEND MESSAGE")
				}
				return
			}
			width, err := strconv.ParseUint(params[1], 10, 32)
			if err != nil {
				_, err = vk.MessagesSend(msg.PeerId, "Invalid width", nil)
				if err != nil {
					log.Println("UNABLE TO SEND MESSAGE")
				}
				return
			}
			height, err := strconv.ParseUint(params[2], 10, 32)
			if err != nil {
				_, err = vk.MessagesSend(msg.PeerId, "Invalid height", nil)
				if err != nil {
					log.Println("UNABLE TO SEND MESSAGE")
				}
				return
			}
			cellSize, err := strconv.ParseUint(params[3], 10, 32)
			if err != nil {
				_, err = vk.MessagesSend(msg.PeerId, "Invalid cell size", nil)
				if err != nil {
					log.Println("UNABLE TO SEND MESSAGE")
				}
				return
			}
			gens, err := strconv.ParseUint(params[4], 10, 32)
			if err != nil {
				_, err = vk.MessagesSend(msg.PeerId, "Invalid generations", nil)
				if err != nil {
					log.Println("UNABLE TO SEND MESSAGE")
				}
				return
			}
			threads, err := strconv.ParseUint(params[5], 10, 32)
			if err != nil {
				_, err = vk.MessagesSend(msg.PeerId, "Invalid threads", nil)
				if err != nil {
					log.Println("UNABLE TO SEND MESSAGE")
				}
				return
			}

			if params[6][0] != 'R' && params[6][0] != 'G' && params[6][0] != 'B' {
				_, err = vk.MessagesSend(msg.PeerId, "Color should be R/G/B", nil)
				if err != nil {
					log.Println("UNABLE TO SEND MESSAGE")
				}
				return
			}

			if width*cellSize > 3840 || height*cellSize > 3840 {
				_, err = vk.MessagesSend(msg.PeerId, "Resolution is too big", nil)
				if err != nil {
					log.Println("UNABLE TO SEND MESSAGE")
				}
				return
			}
			if gens > 1000 {
				_, err = vk.MessagesSend(msg.PeerId, "Too many generations", nil)
				if err != nil {
					log.Println("UNABLE TO SEND MESSAGE")
				}
				return
			}

			var needDifferentColor bool
			if len(params[6]) > 1 && params[6][1] == '1' {
				needDifferentColor = true
			}

			palette := make([]color.Color, 0, 4)
			switch params[6][0] {
			case 'R':
				palette = []color.Color{
					color.RGBA{R: gameOfLife.PrimaryColorSaturation, G: gameOfLife.SecondaryColorSaturation, B: gameOfLife.SecondaryColorSaturation, A: 255},
				}
				if needDifferentColor {
					palette = append(palette, color.RGBA{R: gameOfLife.PrimaryColorSaturation + gameOfLife.PrimaryColorDiff, G: gameOfLife.SecondaryColorSaturation, B: gameOfLife.SecondaryColorSaturation, A: 255})
					palette = append(palette, color.RGBA{R: gameOfLife.PrimaryColorSaturation - gameOfLife.PrimaryColorDiff, G: gameOfLife.SecondaryColorSaturation, B: gameOfLife.SecondaryColorSaturation, A: 255})
				}
			case 'G':
				palette = []color.Color{
					color.RGBA{R: gameOfLife.SecondaryColorSaturation, G: gameOfLife.PrimaryColorSaturation, B: gameOfLife.SecondaryColorSaturation, A: 255},
				}
				if needDifferentColor {
					palette = append(palette, color.RGBA{R: gameOfLife.SecondaryColorSaturation, G: gameOfLife.PrimaryColorSaturation + gameOfLife.PrimaryColorDiff, B: gameOfLife.SecondaryColorSaturation, A: 255})
					palette = append(palette, color.RGBA{R: gameOfLife.SecondaryColorSaturation, G: gameOfLife.PrimaryColorSaturation - gameOfLife.PrimaryColorDiff, B: gameOfLife.SecondaryColorSaturation, A: 255})
				}
			case 'B':
				palette = []color.Color{
					color.RGBA{R: gameOfLife.SecondaryColorSaturation, G: gameOfLife.SecondaryColorSaturation, B: gameOfLife.PrimaryColorSaturation, A: 255},
				}
				if needDifferentColor {
					palette = append(palette, color.RGBA{R: gameOfLife.SecondaryColorSaturation, G: gameOfLife.SecondaryColorSaturation, B: gameOfLife.PrimaryColorSaturation + gameOfLife.PrimaryColorDiff, A: 255})
					palette = append(palette, color.RGBA{R: gameOfLife.SecondaryColorSaturation, G: gameOfLife.SecondaryColorSaturation, B: gameOfLife.PrimaryColorSaturation - gameOfLife.PrimaryColorDiff, A: 255})
				}
			}
			palette = append(palette, color.RGBA{R: 0, G: 0, B: 0, A: 255})

			err = gameOfLife.Generate(uint(width), uint(height), uint(cellSize), uint(gens), uint(threads), palette, "life.gif")
			if err != nil {
				_, err = vk.MessagesSend(msg.PeerId, "Unable to create GIF, "+err.Error(), nil)
				if err != nil {
					log.Println("UNABLE TO SEND MESSAGE")
				}
				return
			}

			uploadServer, err := vk.DocGetMessageUploadServer("doc", msg.PeerId, false)
			if err != nil {
				_, err = vk.MessagesSend(msg.PeerId, "Unable to get upload url, "+err.Error(), nil)
				if err != nil {
					log.Println("UNABLE TO SEND MESSAGE")
				}
				return
			}

			file, err := vk.DocsUploadToMessageServer(uploadServer, "life.gif")
			if err != nil {
				_, err = vk.MessagesSend(msg.PeerId, "Unable to upload file, "+err.Error(), nil)
				if err != nil {
					log.Println("UNABLE TO SEND MESSAGE")
				}
				return
			}

			doc, err := vk.DocsSave(file, "life")
			if err != nil {
				_, err = vk.MessagesSend(msg.PeerId, "Unable to save doc, "+err.Error(), nil)
				if err != nil {
					log.Println("UNABLE TO SEND MESSAGE")
				}
				return
			}

			_, err = vk.MessagesSend(msg.PeerId, "", &doc.Content)
			if err != nil {
				log.Println("UNABLE TO SEND MESSAGE")
			}
			return
		}
	}
}

func vkUserOnlineOfflineLogger(userId, ts int64, isOnline bool) {
	rs, err := vk.UsersGet([]int{int(userId)}, nil, "")
	if err != nil {
		log.Println(err)
		return
	}
	if len(rs.Users) == 0 {
		log.Println("Users.get returned 0 users //ONLINE STATUS LOGGER")
		return
	}
	var logOnlineStatus strings.Builder
	var status string
	if isOnline {
		status = "ONLINE"
	} else {
		status = "OFFLINE"
	}
	_, _ = fmt.Fprintf(&logOnlineStatus,
		"USER \"%s\" BECOMES %s AT %s\n",
		rs.Users[0].FirstName+" "+rs.Users[0].LastName, status,
		time.Unix(ts, 0).Format("15:04:05, 02.01.2006"))

	fmt.Print(logOnlineStatus.String())
}

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

	vk.NewMsgLongPollHandler = vkMessageRouter
	vk.UserOnlineHandler = vkUserOnlineOfflineLogger

	vk.LongPoll(longPollServer)
}
