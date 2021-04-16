package color

import (
	"fmt"
	"math"
	"strconv"
)

func NewFromRGB(name string, red, green, blue int) *Color {
	return &Color{
		name:  name,
		red:   float64(red) / 255,
		green: float64(green) / 255,
		blue:  float64(blue) / 255,
	}
}

func NewFromHex(name, hex string) (*Color, error) {
	if len(hex) != 6 {
		return nil, fmt.Errorf("hex for whole color must have 6 digits")
	}
	r, err := hexToRGB(hex[0:2])
	if err != nil {
		return nil, err
	}
	g, err := hexToRGB(hex[2:4])
	if err != nil {
		return nil, err
	}
	b, err := hexToRGB(hex[4:6])
	if err != nil {
		return nil, err
	}
	return &Color{
		name:  name,
		red:   r,
		green: g,
		blue:  b,
	}, nil
}

// Color represents a color with red, green and blue values, as well as a name.
type Color struct {
	name  string
	red   float64 // 0 - 1
	green float64 // 0 - 1
	blue  float64 // 0 - 1
}

func (c *Color) String() string {
	return fmt.Sprintf("%s (%s)", c.name, c.Hex())
}

// Name returns the receiver's name
func (c *Color) Name() string {
	return c.name
}

// Hex returns a hexadecimal representation of the receiver. For instance:
// #00000 for black.
func (c *Color) Hex() string {
	r := strconv.FormatInt(int64(c.red*255), 16)
	g := strconv.FormatInt(int64(c.green*255), 16)
	b := strconv.FormatInt(int64(c.blue*255), 16)
	return fmt.Sprintf("#%02s%02s%02s", r, g, b)
}

// Luminance returns the receiver's luminance as computed by the formula here:
// https://www.w3.org/TR/WCAG20/#relativeluminancedef
func (c *Color) Luminance() float64 {
	r := sRGB(c.red)
	g := sRGB(c.green)
	b := sRGB(c.blue)

	l := 0.2126*r + 0.7152*g + 0.0722*b

	return l
}

// ContrastRatio returns the contrast ratio between the receiver and another
// given color, as outlined by the process here:
// https://medium.muz.li/the-science-of-color-contrast-an-expert-designers-guide-33e84c41d156
func (c *Color) ContrastRatio(other *Color) float64 {
	thisL := c.Luminance()
	otherL := other.Luminance()
	lighter := thisL
	darker := otherL
	if otherL > thisL {
		lighter = otherL
		darker = thisL
	}
	return (lighter + 0.05) / (darker + 0.05)
}

// DistanceTo returns the distance from the receiver to another given color
// using a naive 3-d space.
func (c *Color) DistanceTo(other *Color) float64 {
	r := math.Abs(c.red - other.red)
	g := math.Abs(c.green - other.green)
	b := math.Abs(c.blue - other.blue)
	return math.Sqrt(math.Pow(r, 2) + math.Pow(g, 2) + math.Pow(b, 2))
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

const unreportedContrastRatioThreshold = 3

func ContrastRatioDescription(contrastRatio float64) string {
	if contrastRatio < unreportedContrastRatioThreshold {
		return "--"
	}
	name := "AAA"
	if contrastRatio < 4.5 {
		name = "AA+"
	} else if contrastRatio < 7 {
		name = "AA"
	}
	return name
}
