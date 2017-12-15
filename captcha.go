package captcha

import (
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"math/rand"
	"time"
	"golang.org/x/image/font/gofont/goregular"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func Make(width int, height int, fontsize int, length int) (image.Image, string) {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	drawer := font.Drawer{}
	f, _ := truetype.Parse(goregular.TTF)
	face := truetype.NewFace(f, &truetype.Options{
		Size: float64(fontsize),
	})
	drawer.Face = face

	drawer.Dst = img
	drawer.Dot = fixed.P(0, height-(height-fontsize)/2)
	chars := generateChars(length)
	lastX := fixed.I(width / (length + 1) / 2)
	drawer.Dot.X = lastX
	for i := 0; i < length; i++ {
		drawer.Src = image.NewUniform(color.RGBA{R: uint8(rand.Int31n(170)), G: uint8(rand.Int31n(170)), B: uint8(rand.Int31n(170)), A: 255})
		drawer.DrawBytes(chars[i: i+1])
		lastX += fixed.I(width / (length + 1))
		drawer.Dot.X = lastX
	}
	opts := DistortionOptions{}
	y := (height-fontsize)/2 + r.Int()%fontsize
	x0 := 0 + r.Int()%20 - 10
	x1 := width + r.Int()%20 - 10
	opts = append(opts, [3]int{
		y, x0, x1,
	})
	img = distortion(img, opts)
	return img, string(chars)
}

type DistortionOptions [][3]int

func cp(img *image.RGBA) (nImage *image.RGBA) {
	nImage = &image.RGBA{}
	nImage.Pix = make([]uint8, len(img.Pix))
	copy(nImage.Pix, img.Pix)
	nImage.Rect = img.Rect
	nImage.Stride = img.Stride
	return
}
func distortion(img *image.RGBA, options DistortionOptions) (nImage *image.RGBA) {
	nImage = cp(img)
	bounds := img.Bounds()
	startX, startY := bounds.Min.X, bounds.Min.Y
	endX, endY := bounds.Max.X, bounds.Max.Y
	lastY := startY
	originWidth := endX - startX
	currentWidth := originWidth
	lastX0 := startX
	options = append(options, [3]int{endY, startX, endX})
	for _, v := range options {
		y := v[0]
		x0, x1 := v[1], v[2]
		height := y - lastY
		toWidth := x1 - x0
		for h := 0; h < height; h++ {
			offset := float64(lastX0) - float64(lastX0-x0)/float64(height)*float64(h)
			lineWidth := float64(currentWidth) - float64(currentWidth-toWidth)/float64(height)*float64(h)
			num := float64(originWidth) / float64(lineWidth)
			for i := startX; i < endX; i++ {
				nImage.Set(i, lastY+h, img.At(int((float64(i)-offset)/num), lastY+h))
			}
		}
		lastY = y
		currentWidth = toWidth
		lastX0 = x0
	}
	return nImage
}

var byteList = make([]byte, 0, 57)

func init() {
	initByteList()
}
func initByteList() {
	var c byte
	for c = '0'; c <= '9'; c++ {
		if c == '1' || c == '6' {
			continue
		}
		byteList = append(byteList, c)
	}
	for c = 'A'; c <= 'Z'; c++ {
		if c == 'L' || c == 'I' {
			continue
		}
		byteList = append(byteList, c)
	}
	for c = 'a'; c <= 'z'; c++ {
		if c == 'l' || c == 'i' || c == 'b' {
			continue
		}
		byteList = append(byteList, c)
	}
}

func generateChars(size int) []byte {
	str := make([]byte, size)
	for ; size > 0; size-- {
		str[size-1] = byteList[r.Uint32()%Uint32(len(byteList))]
	}
	return str
}
