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

type (
	Parameters struct {
		Width       uint64
		Height      uint64
		CellSize    uint64
		Generations uint64
		Routines    uint64
		Palette     color.Palette
		FileName    string
	}
)

func Generate(params *Parameters) error {

	var wg sync.WaitGroup

	widthImg := params.Width * params.CellSize
	heightImg := params.Height * params.CellSize

	if heightImg%params.Routines != 0 {
		return errors.New("высота изображения должна быть кратна количеству потоков")
	}

	field := NewField(params.Height, params.Width)

	field.Palette = params.Palette

	field.FillRandom()

	result := &gif.GIF{}
	for n := uint64(0); n < params.Generations; n++ {
		log.Println("GEN: ", n)

		upLeft := image.Point{X: 0, Y: 0}
		lowRight := image.Point{X: int(widthImg), Y: int(heightImg)}

		img := image.NewPaletted(image.Rectangle{Min: upLeft, Max: lowRight}, params.Palette)

		wg.Add(int(params.Routines))
		for i := uint64(0); i < params.Routines; i++ {
			ii := i
			go func() {
				draw(img, field, heightImg/params.Routines, ii, widthImg, params.CellSize)
				wg.Done()
			}()
		}
		wg.Wait()

		result.Image = append(result.Image, img)
		result.Delay = append(result.Delay, 10)

		field.NextGen()
	}

	f, err := os.Create(params.FileName)
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

func draw(canvas *image.Paletted, field *Field, totalRows, indexRow, width, cellSize uint64) {
	for x := uint64(0); x < width; x++ {
		for y := totalRows * indexRow; y < totalRows*(indexRow+1); y++ {
			canvas.Set(int(x), int(y), field.CGen[y/cellSize][x/cellSize].Color)
		}
	}
}
