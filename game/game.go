package gameOfLife

import (
	"errors"
	"image"
	"image/color"
	"image/gif"
	"log"
	"os"
	"sync"
)

func Gen(w, h, c, g, t uint, palette color.Palette) error {

	var wg sync.WaitGroup

	widthImg := w * c
	heightImg := h * c

	if widthImg%t != 0 || heightImg%t != 0 {
		return errors.New("width or height % 10 != 0")
	}

	field := CreateField(h, w)

	field.Palette = palette

	field.FillRandom()

	result := &gif.GIF{}
	for n := uint(0); n < g; n++ {
		log.Println("GEN: ", n)

		upLeft := image.Point{X: 0, Y: 0}
		lowRight := image.Point{X: int(widthImg), Y: int(heightImg)}

		img := image.NewPaletted(image.Rectangle{Min: upLeft, Max: lowRight}, palette)

		wg.Add(int(t))
		for i := uint(0); i < t; i++ {
			routineNumber := i
			go func() {
				draw(img, field, heightImg/t, routineNumber, widthImg, c)
				wg.Done()
			}()
		}
		wg.Wait()

		result.Image = append(result.Image, img)
		result.Delay = append(result.Delay, 1)

		field.NextGen()
	}

	f, err := os.Create("life.gif")
	if err != nil {
		return err
	}
	err = gif.EncodeAll(f, result)
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	return nil
}

func draw(canvas *image.Paletted, field *Field, rows, i, w, c uint) {
	for x := uint(0); x < w; x++ {
		for y := rows * i; y < rows*(i+1); y++ {
			if field.Zone[y/c][x/c].IsAlive {
				canvas.Set(int(x), int(y), field.Zone[y/c][x/c].Color)
			} else {
				canvas.Set(int(x), int(y), field.Zone[y/c][x/c].Color)
			}
		}
	}
}
