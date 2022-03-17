package main

import (
	"errors"
	"gobotvk/crypto"
	gof "gobotvk/gameOfLife"
	"gobotvk/vk"
	"image/color"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	SleepTimeBeforeDelete = 10 * time.Second

	CommandPrefix = '$'

	GameOfLifeCommand        = "gen-life"
	GameOfLifeMaxWidth       = 3840
	GameOfLifeMaxHeight      = 2160
	GameOfLifeMaxGenerations = 1000

	CryptoHashCommand = "sha256"
)

func vkMessageHandler(msg *vk.NewMessageLongPollEvent) {
	if !(msg != nil && msg.Text != "") {
		return
	}

	if msg.Text[0] == CommandPrefix {
		switch true {

		//~gen-life/width/height/cell/generations/colorStyle
		case strings.Contains(msg.Text, GameOfLifeCommand):
			GenLifeGif(msg)

		//~sha256 (replied message)
		case strings.Contains(msg.Text, CryptoHashCommand):
			HashMessage(msg)
		}
	}
}

func GenLifeGif(msg *vk.NewMessageLongPollEvent) {
	if msg == nil {
		log.Println("Message is nil //$gen-life")
		return
	}

	params, err := validateGameOfLifeParameters(strings.Split(msg.Text, "/"))
	if err != nil {
		NotifyAboutError(msg.PeerId, err)
		return
	}

	params.FileName = "life.gif"

	err = gof.Generate(params)
	if err != nil {
		NotifyAboutError(msg.PeerId, err)
		return
	}

	uploadServer, err := vk.DocGetMessageUploadServer("doc", msg.PeerId, false)
	if err != nil {
		NotifyAboutError(msg.PeerId, err)
		return
	}

	file, err := vk.DocsUploadToMessageServer(&uploadServer, "life.gif")
	if err != nil {
		NotifyAboutError(msg.PeerId, err)
		return
	}

	doc, err := vk.DocsSave(&file, "life")
	if err != nil {
		NotifyAboutError(msg.PeerId, err)
		return
	}

	_, err = vk.MessagesSend(msg.PeerId, "", []*vk.Document{&doc.Content}, msg.Id)
	if err != nil {
		log.Println("~gen-life-gif failed")
	}
	return
}

func validateGameOfLifeParameters(params []string) (*gof.Parameters, error) {
	if len(params) != 6 {
		return nil, errors.New("неверное количество параметров, попробуйте " + string(CommandPrefix) + GameOfLifeCommand + "/w/h/c/g/cS")
	}

	width, err := strconv.ParseUint(params[1], 10, 32)
	if err != nil {
		return nil, errors.New("неверное значение ширины")
	}

	height, err := strconv.ParseUint(params[2], 10, 32)
	if err != nil {
		return nil, errors.New("неверное значение высоты")
	}

	cellSize, err := strconv.ParseUint(params[3], 10, 32)
	if err != nil {
		return nil, errors.New("неверное значение размера клетки")
	}

	if width*cellSize > GameOfLifeMaxWidth || height*cellSize > GameOfLifeMaxHeight {
		return nil, errors.New("итоговый размер изображения слишком большой")
	}

	gens, err := strconv.ParseUint(params[4], 10, 32)
	if err != nil {
		return nil, errors.New("неверное значение количества поколений")
	} else if gens > GameOfLifeMaxGenerations {
		return nil, errors.New("превышено максимальное количество поколений")
	}

	palette := make([]color.Color, 0, 4)
	if len(params[5]) == 2 {
		needDifColor := params[5][1] == '1'
		switch params[5][0] {
		case 'R':
			palette = []color.Color{
				color.RGBA{R: gof.PC, G: gof.SC, B: gof.SC, A: 255},
			}
			if needDifColor {
				palette = append(palette, color.RGBA{R: gof.PC + gof.CD, G: gof.SC, B: gof.SC, A: 255})
				palette = append(palette, color.RGBA{R: gof.PC - gof.CD, G: gof.SC, B: gof.SC, A: 255})
			}
		case 'G':
			palette = []color.Color{
				color.RGBA{R: gof.SC, G: gof.PC, B: gof.SC, A: 255},
			}
			if needDifColor {
				palette = append(palette, color.RGBA{R: gof.SC, G: gof.PC + gof.CD, B: gof.SC, A: 255})
				palette = append(palette, color.RGBA{R: gof.SC, G: gof.PC - gof.CD, B: gof.SC, A: 255})
			}
		case 'B':
			palette = []color.Color{
				color.RGBA{R: gof.SC, G: gof.SC, B: gof.PC, A: 255},
			}
			if needDifColor {
				palette = append(palette, color.RGBA{R: gof.SC, G: gof.SC, B: gof.PC + gof.CD, A: 255})
				palette = append(palette, color.RGBA{R: gof.SC, G: gof.SC, B: gof.PC - gof.CD, A: 255})
			}
		default:
			return nil, errors.New("доступны только цвета R, G, B")
		}
		palette = append(palette, color.RGBA{R: 0, G: 0, B: 0, A: 255})
	} else {
		return nil, errors.New("выберите цвет и стиль клеток")
	}

	r := new(gof.Parameters)
	r.Width = width
	r.Height = height
	r.CellSize = cellSize
	r.Generations = gens
	r.Palette = palette

	return r, nil
}

func HashMessage(msg *vk.NewMessageLongPollEvent) {
	if msg == nil {
		log.Println("Message is nil //$sha256")
		return
	}

	if msg.RepliedId == 0 {
		NotifyAboutError(msg.PeerId, errors.New("выберите нужное сообщение ответив на него"))
		return
	}

	messages, err := vk.MessagesGetFromConversation([]int64{msg.RepliedId}, msg.PeerId)
	if err != nil {
		log.Println(err)
	}

	if len(messages.Response.Messages) == 0 {
		NotifyAboutError(msg.PeerId, errors.New("ошибка получения целевого сообщения"))
		return
	}

	_, err = vk.MessagesSend(msg.PeerId, "SHA256: "+crypto.SHA256([]byte(messages.Response.Messages[0].Text)).ToHexString(), nil, messages.Response.Messages[0].Id)
	if err != nil {
		NotifyAboutError(msg.PeerId, err)
		return
	}
}

func NotifyAboutError(peerId int64, err error) {
	notificationMsg, notifyErr := vk.MessagesSend(peerId, "Произошла ошибка: "+err.Error(), nil, 0)
	if notifyErr != nil {
		log.Println("Failed to notify about: " + err.Error())
		return
	}

	time.Sleep(SleepTimeBeforeDelete)

	notifyErr = vk.MessagesDelete([]int64{notificationMsg.NewMessageId}, true)
	if notifyErr != nil {
		log.Println("Failed to delete notification: " + err.Error())
		return
	}
}
