package captcha

import (
	"bytes"
	"crypto/rand"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math/big"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

const (
	chars      = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	strLen     = 6
	baseWidth  = 120
	baseHeight = 40
	scale      = 3
)

func Generate() (string, []byte, error) {
	code, err := randomString(strLen)
	if err != nil {
		return "", nil, err
	}

	img := image.NewRGBA(image.Rect(0, 0, baseWidth, baseHeight))
	
	draw.Draw(img, img.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)

	addNoise(img)

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.Black),
		Face: basicfont.Face7x13,
	}

	startX := 10
	for _, char := range code {
		yOffset, _ := rand.Int(rand.Reader, big.NewInt(10))
		d.Dot = fixed.Point26_6{
			X: fixed.I(startX),
			Y: fixed.I(20 + int(yOffset.Int64())),
		}
		d.DrawString(string(char))
		startX += 15 + int(yOffset.Int64()%5)
	}

	finalImg := scaleImage(img, scale)

	var buf bytes.Buffer
	if err := png.Encode(&buf, finalImg); err != nil {
		return "", nil, err
	}

	return code, buf.Bytes(), nil
}

func randomString(length int) (string, error) {
	ret := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}
		ret[i] = chars[num.Int64()]
	}
	return string(ret), nil
}

func addNoise(img *image.RGBA) {
	bounds := img.Bounds()
	for i := 0; i < 100; i++ {
		x, _ := rand.Int(rand.Reader, big.NewInt(int64(bounds.Max.X)))
		y, _ := rand.Int(rand.Reader, big.NewInt(int64(bounds.Max.Y)))
		img.Set(int(x.Int64()), int(y.Int64()), color.RGBA{
			R: uint8(x.Int64() % 255),
			G: uint8(y.Int64() % 255),
			B: 100,
			A: 255,
		})
	}
	for i := 0; i < 5; i++ {
		x1, _ := rand.Int(rand.Reader, big.NewInt(int64(bounds.Max.X)))
		y1, _ := rand.Int(rand.Reader, big.NewInt(int64(bounds.Max.Y)))
		x2, _ := rand.Int(rand.Reader, big.NewInt(int64(bounds.Max.X)))
		y2, _ := rand.Int(rand.Reader, big.NewInt(int64(bounds.Max.Y)))
		line(img, int(x1.Int64()), int(y1.Int64()), int(x2.Int64()), int(y2.Int64()), color.RGBA{150, 150, 150, 255})
	}
}

func line(img *image.RGBA, x1, y1, x2, y2 int, col color.Color) {
	dx := abs(x2 - x1)
	dy := abs(y2 - y1)
	sx := -1
	if x1 < x2 {
		sx = 1
	}
	sy := -1
	if y1 < y2 {
		sy = 1
	}
	err := dx - dy
	for {
		img.Set(x1, y1, col)
		if x1 == x2 && y1 == y2 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x1 += sx
		}
		if e2 < dx {
			err += dx
			y1 += sy
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func scaleImage(src image.Image, factor int) image.Image {
	bounds := src.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, bounds.Max.X*factor, bounds.Max.Y*factor))
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			col := src.At(x, y)
			for dy := 0; dy < factor; dy++ {
				for dx := 0; dx < factor; dx++ {
					dst.Set(x*factor+dx, y*factor+dy, col)
				}
			}
		}
	}
	return dst
}