package util

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"strconv"
)

var MyFS draw.Drawer = myFS{}

type myFS struct{}

func (myFS) Draw(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point) {
	//TODO `src` からパレットを作成する
	p, err := createPalette(dst, src)
	if err != nil {
		fmt.Println(err.Error())
	}
	pImg := dst.(*image.Paletted)
	pImg.Palette = append(p, color.Black, color.White)
	dst = pImg
}

func createPalette(dst draw.Image, src image.Image) (color.Palette, error) {
	cMap := make(map[color.Color]bool)
	r := src.Bounds()
	for y := 0; y != r.Dy(); y++ {
		for x := 0; x != r.Dx(); x++ {
			c := src.At(x, y)
			r, g, b, _ := c.RGBA()               //uint32
			nr, ng, nb := convertColors(r, g, b) //uint8
			fmt.Printf("R: %v, G: %v, B: %v\n", nr, ng, nb)
			var nc color.Color = color.RGBA{nr, ng, nb, 255}
			if _, ok := cMap[nc]; !ok {
				cMap[nc] = true
			}
			dst.Set(x, y, nc)
		}
	}
	if len(cMap) >= 256 {
		return nil, fmt.Errorf("palette size over")
	}
	p := make([]color.Color, len(cMap))
	for c, _ := range cMap {
		p = append(p, c)
	}
	return p, nil
}

func convertColors(r, g, b uint32) (nr, ng, nb uint8) {
	nr, ng, nb = straightColor(r), straightColor(g), straightColor(b)
	return
}

func straightColor(v uint32) uint8 {
	v = v >> 8 / 42 //あまりを捨てることで6段階に区分けする
	v *= 42         //MAX 252
	return uint8(v)
}

func format(r, g, b uint32) (string, string, string) {
	return strconv.FormatUint(uint64(r), 16), strconv.FormatUint(uint64(g), 16), strconv.FormatUint(uint64(b), 16)
}
