package gameOfLife

import (
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
		Palette     color.Palette
		FileName    string
	}
)

func Generate(params *Parameters) error {

	var wg sync.WaitGroup

	widthImg := params.Width * params.CellSize
	heightImg := params.Height * params.CellSize

	field := NewField(params.Height, params.Width)

	field.Palette = params.Palette

	field.FillRandom()

	result := &gif.GIF{}
	for n := uint64(0); n < params.Generations; n++ {
		log.Println("GEN: ", n)

		upLeft := image.Point{X: 0, Y: 0}
		lowRight := image.Point{X: int(widthImg), Y: int(heightImg)}

		img := image.NewPaletted(image.Rectangle{Min: upLeft, Max: lowRight}, params.Palette)

		routines := heightImg
		if routines > 100 {
			routines = 100
		}
		wg.Add(int(routines))
		for i := uint64(0); i < routines; i++ {
			ii := i
			go func() {
				draw(img, field, heightImg/routines, ii, widthImg, params.CellSize)
				wg.Done()
			}()
		}
		wg.Wait()

		result.Image = append(result.Image, img)
		result.Delay = append(result.Delay, 0)

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
