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

func Generate(width, height, cellSize, generations, routines uint, palette color.Palette, fileName string) error {

	var wg sync.WaitGroup

	widthImg := width * cellSize
	heightImg := height * cellSize

	if widthImg%routines != 0 || heightImg%routines != 0 {
		return errors.New("invalid amount of threads")
	}

	field := NewField(height, width)

	field.Palette = palette

	field.FillRandom()

	result := &gif.GIF{}
	for n := uint(0); n < generations; n++ {
		log.Println("GEN: ", n)

		upLeft := image.Point{X: 0, Y: 0}
		lowRight := image.Point{X: int(widthImg), Y: int(heightImg)}

		img := image.NewPaletted(image.Rectangle{Min: upLeft, Max: lowRight}, palette)

		wg.Add(int(routines))
		for i := uint(0); i < routines; i++ {
			ii := i
			go func() {
				draw(img, field, heightImg/routines, ii, widthImg, cellSize)
				wg.Done()
			}()
		}
		wg.Wait()

		result.Image = append(result.Image, img)
		result.Delay = append(result.Delay, 10)

		field.NextGen()
	}

	f, err := os.Create(fileName)
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

func draw(canvas *image.Paletted, field *Field, totalRows, indexRow, width, cellSize uint) {
	for x := uint(0); x < width; x++ {
		for y := totalRows * indexRow; y < totalRows*(indexRow+1); y++ {
			canvas.Set(int(x), int(y), field.CGen[y/cellSize][x/cellSize].Color)
		}
	}
}
