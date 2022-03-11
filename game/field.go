package gameOfLife

import (
	"image/color"
	"math/rand"
	"time"
)

const (
	SecondaryColorSaturation = 25
	PrimaryColorSaturation   = 155
	PrimaryColorDiff         = 30
)

var R = rand.New(rand.NewSource(time.Now().UnixNano()))

type (
	Field struct {
		Zone    [][]Cell
		Palette color.Palette
	}
	Cell struct {
		IsAlive bool
		Color   color.Color
	}
)

func CreateField(m, n uint) *Field {
	f := new(Field)
	f.Zone = make([][]Cell, m)
	for i := uint(0); i < m; i++ {
		f.Zone[i] = make([]Cell, n)
	}
	return f
}

func (f *Field) FillRandom() {
	for i := 0; i < len(f.Zone); i++ {
		for j := 0; j < len(f.Zone[i]); j++ {
			f.Zone[i][j].IsAlive = R.Uint32()%2 == 1
			if f.Zone[i][j].IsAlive {
				f.Zone[i][j].Color = f.Palette[R.Uint32()%uint32(len(f.Palette)-1)]
			} else {
				f.Zone[i][j].Color = f.Palette[len(f.Palette)-1]
			}
		}
	}
}

func (f *Field) NextGen() {
	copyField := CreateField(uint(len(f.Zone)), uint(len(f.Zone[0])))

	for i := range f.Zone {
		for j := range f.Zone[i] {
			copyField.Zone[i][j].IsAlive = f.Zone[i][j].IsAlive
			copyField.Zone[i][j].Color = f.Zone[i][j].Color
		}
	}

	for i := range copyField.Zone {
		for j := range copyField.Zone[i] {
			neighbours := copyField.GetNeighbours(i, j)
			if !copyField.Zone[i][j].IsAlive && neighbours == 3 {
				f.Zone[i][j].IsAlive = true
				f.Zone[i][j].Color = f.Palette[R.Uint32()%uint32(len(f.Palette)-1)]
			}
			if copyField.Zone[i][j].IsAlive {
				if neighbours == 2 || neighbours == 3 {
					f.Zone[i][j].IsAlive = true
					f.Zone[i][j].Color = f.Palette[R.Uint32()%uint32(len(f.Palette)-1)]
				} else {
					f.Zone[i][j].IsAlive = false
					f.Zone[i][j].Color = color.RGBA{R: 0, G: 0, B: 0, A: 255}
				}
			}
		}
	}
}

func (f Field) GetNeighbours(i, j int) uint {
	rows := len(f.Zone)
	cols := len(f.Zone[i])
	//todo optimize
	var neighbours uint
	for l := i - 1; l <= i+1; l++ {
		l2 := l
		if l == -1 {
			l = rows - 1
		}
		l %= rows
		for m := j - 1; m <= j+1; m++ {
			m2 := m
			if m == -1 {
				m = cols - 1
			}
			m %= cols

			if (l != i || m != j) && f.Zone[l][m].IsAlive {
				neighbours++
			}
			m = m2
		}
		l = l2
	}
	return neighbours
}
