package captcha

import (
	crand "crypto/rand"
	"image"
	"image/color"
	"image/png"
	"io"
	"math/rand"
	"strconv"
	"time"
)

const (
	stdWidth  = 100
	stdHeight = 40
	maxSkew   = 2
)

const (
	fontWidth  = 5
	fontHeight = 8
	blackChar  = 1
)

var font = [][]byte{
	{ // 0
		0, 1, 1, 1, 0,
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1,
		0, 1, 1, 1, 0,
	},
	{ // 1
		0, 0, 1, 0, 0,
		0, 1, 1, 0, 0,
		1, 0, 1, 0, 0,
		0, 0, 1, 0, 0,
		0, 0, 1, 0, 0,
		0, 0, 1, 0, 0,
		0, 0, 1, 0, 0,
		1, 1, 1, 1, 1,
	},
	{ // 2
		0, 1, 1, 1, 0,
		1, 0, 0, 0, 1,
		0, 0, 0, 0, 1,
		0, 0, 0, 1, 1,
		0, 1, 1, 0, 0,
		1, 0, 0, 0, 0,
		1, 0, 0, 0, 0,
		1, 1, 1, 1, 1,
	},
	{ // 3
		1, 1, 1, 1, 0,
		0, 0, 0, 0, 1,
		0, 0, 0, 1, 0,
		0, 1, 1, 1, 0,
		0, 0, 0, 1, 0,
		0, 0, 0, 0, 1,
		0, 0, 0, 0, 1,
		1, 1, 1, 1, 0,
	},
	{ // 4
		1, 0, 0, 1, 0,
		1, 0, 0, 1, 0,
		1, 0, 0, 1, 0,
		1, 0, 0, 1, 0,
		1, 1, 1, 1, 1,
		0, 0, 0, 1, 0,
		0, 0, 0, 1, 0,
		0, 0, 0, 1, 0,
	},
	{ // 5
		1, 1, 1, 1, 1,
		1, 0, 0, 0, 0,
		1, 0, 0, 0, 0,
		1, 1, 1, 1, 0,
		0, 0, 0, 0, 1,
		0, 0, 0, 0, 1,
		0, 0, 0, 0, 1,
		1, 1, 1, 1, 0,
	},
	{ // 6
		0, 0, 1, 1, 1,
		0, 1, 0, 0, 0,
		1, 0, 0, 0, 0,
		1, 1, 1, 1, 0,
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1,
		0, 1, 1, 1, 0,
	},
	{ // 7
		1, 1, 1, 1, 1,
		0, 0, 0, 0, 1,
		0, 0, 0, 0, 1,
		0, 0, 0, 1, 0,
		0, 0, 1, 0, 0,
		0, 1, 0, 0, 0,
		0, 1, 0, 0, 0,
		0, 1, 0, 0, 0,
	},
	{ // 8
		0, 1, 1, 1, 0,
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1,
		0, 1, 1, 1, 0,
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1,
		0, 1, 1, 1, 0,
	},
	{ // 9
		0, 1, 1, 1, 0,
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1,
		1, 1, 0, 0, 1,
		0, 1, 1, 1, 1,
		0, 0, 0, 0, 1,
		0, 0, 0, 0, 1,
		1, 1, 1, 1, 0,
	},
}

type Image struct {
	*image.NRGBA
	color   *color.NRGBA
	width   int //a digit width
	height  int //a digit height
	dotsize int
}

func init() {
	rand.Seed(int64(time.Second))
}

func NewImage(digits []byte, width, height int) *Image {
	img := new(Image)
	r := image.Rect(img.width, img.height, stdWidth, stdHeight)
	img.NRGBA = image.NewNRGBA(r)

	img.color = &color.NRGBA{
		uint8(rand.Intn(129)),
		uint8(rand.Intn(129)),
		uint8(rand.Intn(129)),
		0xFF,
	}
	// Draw background (10 random circles of random brightness)
	img.calculateSizes(width, height, len(digits))
	img.fillWithCircles(10, img.dotsize)

	maxx := width - (img.width+img.dotsize)*len(digits) - img.dotsize
	maxy := height - img.height - img.dotsize*2

	x := rnd(img.dotsize*2, maxx)
	y := rnd(img.dotsize*2, maxy)

	// Draw digits.
	for _, n := range digits {
		img.drawDigit(font[n], x, y)
		x += img.width + img.dotsize
	}

	// Draw strike-through line.
	// 中间线不要
	//img.strikeThrough()

	return img
}

func (img *Image) WriteTo(w io.Writer) (int64, error) {
	return 0, png.Encode(w, img)
}

func (img *Image) calculateSizes(width, height, ncount int) {

	// Goal: fit all digits inside the image.
	var border int
	if width > height {
		border = height / 5
	} else {
		border = width / 5
	}
	// Convert everything to floats for calculations.
	w := float64(width - border*2)  //268
	h := float64(height - border*2) //48
	// fw takes into account 1-dot spacing between digits.

	fw := float64(fontWidth) + 1 //6

	fh := float64(fontHeight) //8
	nc := float64(ncount)     //7

	// Calculate the width of a single digit taking into account only the
	// width of the image.
	nw := w / nc //38
	// Calculate the height of a digit from this width.
	nh := nw * fh / fw //51

	// Digit too high?

	if nh > h {
		// Fit digits based on height.
		nh = h //nh = 44
		nw = fw / fh * nh
	}
	// Calculate dot size.
	img.dotsize = int(nh / fh)
	// Save everything, making the actual width smaller by 1 dot to account
	// for spacing between digits.
	img.width = int(nw)
	img.height = int(nh) - img.dotsize
}

func (img *Image) fillWithCircles(n, maxradius int) {
	color := img.color
	maxx := img.Bounds().Max.X
	maxy := img.Bounds().Max.Y
	for i := 0; i < n; i++ {
		setRandomBrightness(color, 255)
		r := rnd(1, maxradius)
		img.drawCircle(color, rnd(r, maxx-r), rnd(r, maxy-r), r)
	}
}

func (img *Image) drawHorizLine(color color.Color, fromX, toX, y int) {
	for x := fromX; x <= toX; x++ {
		img.Set(x, y, color)
	}
}

func (img *Image) drawCircle(color color.Color, x, y, radius int) {
	f := 1 - radius
	dfx := 1
	dfy := -2 * radius
	xx := 0
	yy := radius

	img.Set(x, y+radius, color)
	img.Set(x, y-radius, color)
	img.drawHorizLine(color, x-radius, x+radius, y)

	for xx < yy {
		if f >= 0 {
			yy--
			dfy += 2
			f += dfy
		}
		xx++
		dfx += 2
		f += dfx
		img.drawHorizLine(color, x-xx, x+xx, y+yy)
		img.drawHorizLine(color, x-xx, x+xx, y-yy)
		img.drawHorizLine(color, x-yy, x+yy, y+xx)
		img.drawHorizLine(color, x-yy, x+yy, y-xx)
	}
}

func (img *Image) strikeThrough() {
	r := 0
	maxx := img.Bounds().Max.X
	maxy := img.Bounds().Max.Y
	y := rnd(maxy/3, maxy-maxy/3)
	for x := 0; x < maxx; x += r {
		r = rnd(1, img.dotsize/3)
		y += rnd(-img.dotsize/2, img.dotsize/2)
		if y <= 0 || y >= maxy {
			y = rnd(maxy/3, maxy-maxy/3)
		}
		img.drawCircle(img.color, x, y, r)
	}
}

func (img *Image) drawDigit(digit []byte, x, y int) {
	skf := rand.Float64() * float64(rnd(-maxSkew, maxSkew))
	xs := float64(x)
	minr := img.dotsize / 2               // minumum radius
	maxr := img.dotsize/2 + img.dotsize/4 // maximum radius
	y += rnd(-minr, minr)
	for yy := 0; yy < fontHeight; yy++ {
		for xx := 0; xx < fontWidth; xx++ {
			if digit[yy*fontWidth+xx] != blackChar {
				continue
			}
			// Introduce random variations.
			or := rnd(minr, maxr)
			ox := x + (xx * img.dotsize) + rnd(0, or/2)
			oy := y + (yy * img.dotsize) + rnd(0, or/2)

			img.drawCircle(img.color, ox, oy, or)
		}
		xs += skf
		x = int(xs)
	}
}

func setRandomBrightness(c *color.NRGBA, max uint8) {
	minc := min3(c.R, c.G, c.B)
	maxc := max3(c.R, c.G, c.B)
	if maxc > max {
		return
	}
	n := rand.Intn(int(max-maxc)) - int(minc)
	c.R = uint8(int(c.R) + n)
	c.G = uint8(int(c.G) + n)
	c.B = uint8(int(c.B) + n)
}

func min3(x, y, z uint8) (o uint8) {
	o = x
	if y < o {
		o = y
	}
	if z < o {
		o = z
	}
	return
}

func max3(x, y, z uint8) (o uint8) {
	o = x
	if y > o {
		o = y
	}
	if z > o {
		o = z
	}
	return
}

// rnd returns a random number in range [from, to].
func rnd(from, to int) int {
	//println(to+1-from)
	return rand.Intn(to+1-from) + from
}

const (
	// Standard length of uniuri string to achive ~95 bits of entropy.
	StdLen = 16
	// Length of uniurl string to achive ~119 bits of entropy, closest
	// to what can be losslessly converted to UUIDv4 (122 bits).
	UUIDLen = 20
)

// Standard characters allowed in uniuri string.
var StdChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

// New returns a new random string of the standard length, consisting of
// standard characters.
func New() string {
	return NewLenChars(StdLen, StdChars)
}

// NewLen returns a new random string of the provided length, consisting of
// standard characters.
func NewLen(length int) string {
	return NewLenChars(length, StdChars)
}

// NewLenChars returns a new random string of the provided length, consisting
// of the provided byte slice of allowed characters (maximum 256).
func NewLenChars(length int, chars []byte) string {
	b := make([]byte, length)
	r := make([]byte, length+(length/4)) // storage for random bytes.
	clen := byte(len(chars))
	maxrb := byte(256 - (256 % len(chars)))
	i := 0
	for {
		if _, err := io.ReadFull(crand.Reader, r); err != nil {
			panic("error reading from random source: " + err.Error())
		}
		for _, c := range r {
			if c >= maxrb {
				// Skip this number to avoid modulo bias.
				continue
			}
			b[i] = chars[c%clen]
			i++
			if i == length {
				return string(b)
			}
		}
	}
	panic("unreachable")
}

func Fetch() (*Image, string) {
	d := make([]byte, 4)
	s := NewLen(4)
	ss := ""
	d = []byte(s)
	for v := range d {
		d[v] %= 10
		ss += strconv.FormatInt(int64(d[v]), 32)
	}
	return NewImage(d, 100, 40), ss
}
