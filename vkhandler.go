package main

import (
	gof "bot/game"
	"bot/vk"
	"errors"
	"image/color"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	SleepTimeBeforeDelete    = 5 * time.Second
	GameOfLifeMaxWidth       = 3840
	GameOfLifeMaxHeight      = 2160
	GameOfLifeMaxGenerations = 1000
)

func vkMessageHandler(msg *vk.NewMessageLongPollEvent) {
	if !(msg != nil && msg.Text != "") {
		return
	}

	if msg.Text[0] == '~' {
		switch true {

		//~gen-life-gif/width/height/cell/generations/routines/colorStyle
		case strings.Contains(msg.Text, "gen-life-gif"):
			GenLifeGif(msg)
		}

	}
}

func GenLifeGif(msg *vk.NewMessageLongPollEvent) {
	params, err := validateGameOfLifeParameters(strings.Split(msg.Text, "/"))
	if err != nil {
		NotifyAboutError(msg.Id, msg.PeerId, err)
		return
	}

	params.FileName = "life.gif"

	err = gof.Generate(params)
	if err != nil {
		NotifyAboutError(msg.Id, msg.PeerId, err)
		return
	}

	uploadServer, err := vk.DocGetMessageUploadServer("doc", msg.PeerId, false)
	if err != nil {
		NotifyAboutError(msg.Id, msg.PeerId, err)
		return
	}

	file, err := vk.DocsUploadToMessageServer(uploadServer, "life.gif")
	if err != nil {
		NotifyAboutError(msg.Id, msg.PeerId, err)
		return
	}

	doc, err := vk.DocsSave(file, "life")
	if err != nil {
		NotifyAboutError(msg.Id, msg.PeerId, err)
		return
	}

	_, err = vk.MessagesSend(msg.PeerId, "", &doc.Content, msg.Id)
	if err != nil {
		log.Println("~gen-life-gif failed")
	}
	return
}

func validateGameOfLifeParameters(params []string) (*gof.Parameters, error) {
	if len(params) != 7 {
		return nil, errors.New("неверное количество параметров, попробуйте ~gen-life-gif/w/h/c/g/t/cS")
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

	routines, err := strconv.ParseUint(params[5], 10, 32)
	if err != nil {
		return nil, errors.New("неверное значение количества потоков")
	} else if height*cellSize%routines != 0 {
		return nil, errors.New("высота изображения должна быть кратна количеству потоков")
	}

	palette := make([]color.Color, 0, 4)
	if len(params[6]) == 2 {
		needDifColor := params[6][1] == '1'
		switch params[6][0] {
		case 'R':
			palette = []color.Color{
				color.RGBA{R: gof.PC, G: gof.SC, B: gof.SC, A: 255},
			}
			if needDifColor {
				palette = append(palette, color.RGBA{R: gof.PC + gof.PCD, G: gof.SC, B: gof.SC, A: 255})
				palette = append(palette, color.RGBA{R: gof.PC - gof.PCD, G: gof.SC, B: gof.SC, A: 255})
			}
		case 'G':
			palette = []color.Color{
				color.RGBA{R: gof.SC, G: gof.PC, B: gof.SC, A: 255},
			}
			if needDifColor {
				palette = append(palette, color.RGBA{R: gof.SC, G: gof.PC + gof.PCD, B: gof.SC, A: 255})
				palette = append(palette, color.RGBA{R: gof.SC, G: gof.PC - gof.PCD, B: gof.SC, A: 255})
			}
		case 'B':
			palette = []color.Color{
				color.RGBA{R: gof.SC, G: gof.SC, B: gof.PC, A: 255},
			}
			if needDifColor {
				palette = append(palette, color.RGBA{R: gof.SC, G: gof.SC, B: gof.PC + gof.PCD, A: 255})
				palette = append(palette, color.RGBA{R: gof.SC, G: gof.SC, B: gof.PC - gof.PCD, A: 255})
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
	r.Routines = routines
	r.Palette = palette

	return r, nil
}

func NotifyAboutError(messageId, peerId int64, err error) {
	notifyErr := vk.MessagesEdit(messageId, peerId, "Произошла ошибка: "+err.Error())
	if notifyErr != nil {
		log.Println("Failed to notify about: " + err.Error())
		return
	}

	time.Sleep(SleepTimeBeforeDelete)

	notifyErr = vk.MessagesDelete([]int64{messageId}, true)
	if notifyErr != nil {
		log.Println("Failed to delete notification: " + err.Error())
		return
	}
}
