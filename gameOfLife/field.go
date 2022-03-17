package gameOfLife

import (
	"image/color"
	"math/rand"
	"time"
)

const (
	SC = 45  //SECONDARY COLOR
	PC = 175 //PRIMARY COLOR
	CD = 30  //COLOR DIFFERENCE
)

var R = rand.New(rand.NewSource(time.Now().UnixNano()))

type (
	Field struct {
		CGen    [][]Cell
		nGen    [][]Cell
		Palette color.Palette
	}
	Cell struct {
		IsAlive bool
		Color   color.Color
	}
)

func NewField(m, n uint64) *Field {
	f := new(Field)
	f.CGen = make([][]Cell, m)
	f.nGen = make([][]Cell, m)
	for i := uint64(0); i < m; i++ {
		f.CGen[i] = make([]Cell, n)
		f.nGen[i] = make([]Cell, n)
	}
	return f
}

func (f *Field) FillRandom() {
	for i := 0; i < len(f.CGen); i++ {
		for j := 0; j < len(f.CGen[i]); j++ {
			f.CGen[i][j].IsAlive = R.Uint32()%2 == 1
			if f.CGen[i][j].IsAlive {
				f.CGen[i][j].Color = f.Palette[R.Uint32()%uint32(len(f.Palette)-1)]
			} else {
				f.CGen[i][j].Color = f.Palette[len(f.Palette)-1]
			}
		}
	}
}

func (f *Field) NextGen() {
	for i := range f.CGen {
		for j := range f.CGen[i] {
			neighbours := f.GetNeighbours(i, j)
			if !f.CGen[i][j].IsAlive && neighbours == 3 || f.CGen[i][j].IsAlive && (neighbours == 2 || neighbours == 3) {
				f.nGen[i][j].IsAlive = true
				f.nGen[i][j].Color = f.Palette[R.Uint32()%uint32(len(f.Palette)-1)]
			} else {
				f.nGen[i][j].IsAlive = false
				f.nGen[i][j].Color = f.Palette[len(f.Palette)-1]
			}
		}
	}
	f.CGen, f.nGen = f.nGen, f.CGen
}

func (f Field) GetNeighbours(i, j int) uint {
	rows := len(f.CGen)
	cols := len(f.CGen[i])
	var neighbours uint
	for l := i - 1; l <= i+1; l++ {
		lf := (l + rows) % rows
		for m := j - 1; m <= j+1; m++ {
			mf := (m + cols) % cols
			if !(lf == i && mf == j) && f.CGen[lf][mf].IsAlive {
				neighbours++
			}
		}
	}
	return neighbours
}
