package color

import (
	"fmt"
	"math"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func NewFromRGB(name string, red, green, blue int) Color {
	return Color{
		name:  name,
		red:   float64(red) / 255,
		green: float64(green) / 255,
		blue:  float64(blue) / 255,
	}
}

func NewFromHex(name, hex string) Color {
	if len(hex) != 6 {
		log.Fatal("hex for whole color must have 6 digits")
	}
	r, err := hexToRGB(hex[0:2])
	if err != nil {
		log.Fatal(err)
	}
	g, err := hexToRGB(hex[2:4])
	if err != nil {
		log.Fatal(err)
	}
	b, err := hexToRGB(hex[4:6])
	if err != nil {
		log.Fatal(err)
	}
	return Color{
		name:  name,
		red:   r,
		green: g,
		blue:  b,
	}
}

type Color struct {
	name  string
	red   float64 // 0 - 1
	green float64 // 0 - 1
	blue  float64 // 0 - 1
}

func (c Color) String() string {
	return fmt.Sprintf("%s (%s)", c.name, c.Hex())
}

func (c Color) Name() string {
	return c.name
}

func (c Color) Hex() string {
	r := strconv.FormatInt(int64(c.red*255), 16)
	g := strconv.FormatInt(int64(c.green*255), 16)
	b := strconv.FormatInt(int64(c.blue*255), 16)
	return fmt.Sprintf("#%02s%02s%02s", r, g, b)
}

// https://www.w3.org/TR/WCAG20/#relativeluminancedef
func (c Color) Luminance() float64 {
	r := sRGB(c.red)
	g := sRGB(c.green)
	b := sRGB(c.blue)

	l := 0.2126*r + 0.7152*g + 0.0722*b

	log.Debugf("luminance of %s = %f", c, l)
	return l
}

func hexToRGB(in string) (float64, error) {
	if len(in) != 2 {
		return 0, fmt.Errorf("hex for single value must have 2 digits")
	}
	i, err := strconv.ParseInt(in, 16, 0)
	if err != nil {
		return 0, err
	}
	f := float64(i) / 255
	return f, nil
}

func sRGB(in float64) float64 {
	if in > 0.03928 {
		return math.Pow(((in + 0.055) / 1.055), 2.4)
	}
	return in / 12.92
}
