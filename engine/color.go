package engine

import (
	"encoding/hex"
	"image/color"
	"math"
	"strconv"
	"strings"

	"github.com/qbradq/eye-engine/internal/util"
)

// ColorIndex represents a color from the palette
type ColorIndex uint16

// ColorIndexTransparent is the transparent value for ColorIndex.
const ColorIndexTransparent ColorIndex = 0xFFFF

// HSV represents the Hue, Saturation, and Value components of a color.
type HSV [3]float32

// Rotate returns the color with the hue rotated by d degrees.
func (hsv HSV) Rotate(d float32) HSV {
	h := hsv[0]
	h *= 360
	h += d
	h = float32(math.Mod(float64(h), 360))
	h /= 360
	return HSV{h, hsv[1], hsv[2]}
}

// ColorF32 represents a color in RGBA space using float32's.
type ColorF32 [4]float32

// NewColorFromUint8 returns a new Color values based on 8-bit RGBA values.
func NewColorFromUint8(r, g, b, a uint8) ColorF32 {
	return ColorF32{
		float32(r) / 255,
		float32(g) / 255,
		float32(b) / 255,
		float32(a) / 255,
	}
}

// NewColorFromHex returns a new Color value based on a hex string.
func NewColorFromHex(hex string) ColorF32 {
	// Trim hash if present
	hex = strings.TrimPrefix(hex, "#")
	// Convert hex shortcodes like "f70" to full codes like "ff7700"
	if len(hex) == 3 {
		hex = string([]byte{
			hex[0], hex[0],
			hex[1], hex[1],
			hex[2], hex[2],
			'f', 'f',
		})
	}
	if len(hex) == 4 {
		hex = string([]byte{
			hex[0], hex[0],
			hex[1], hex[1],
			hex[2], hex[2],
			hex[3], hex[3],
		})
	}
	// Extend RGB codes to RGBA
	if len(hex) == 6 {
		hex += "ff"
	}
	// Make sure we have a full hex code
	if len(hex) != 8 {
		return ColorF32{0, 0, 0, 1}
	}
	fn := func(hex string) float32 {
		v, err := strconv.ParseUint(hex[0:2], 16, 8)
		if err != nil {
			return 0
		}
		return float32(v) / 255
	}
	return ColorF32{
		fn(hex[0:2]),
		fn(hex[2:4]),
		fn(hex[4:6]),
		fn(hex[6:8]),
	}
}

// NewColorFromHSV returns a new Color value in RGBA space from the HSV values
// given.
func NewColorFromHSV(hsv HSV, a float32) ColorF32 {
	// Grayscale case
	if hsv[1] == 0 {
		return ColorF32{hsv[2], hsv[2], hsv[2], a}
	}
	// Intermediate values
	i := int(hsv[0]*6) % 6
	f := (hsv[0] * 6) - float32(i)
	p := hsv[2] * (1 - hsv[1])
	q := hsv[2] * (1 - f*hsv[1])
	t := hsv[2] * (1 - (1-f)*hsv[1])
	// Determine RGB based on hue sector
	switch i {
	case 0:
		return ColorF32{hsv[2], t, p, a}
	case 1:
		return ColorF32{q, hsv[2], p, a}
	case 2:
		return ColorF32{p, hsv[2], t, a}
	case 3:
		return ColorF32{p, q, hsv[2], a}
	case 4:
		return ColorF32{t, p, hsv[2], a}
	default:
		return ColorF32{hsv[2], p, q, a}
	}
}

// RGBA implements the color.Color interface.
func (c ColorF32) RGBA() (r, g, b, a uint32) {
	r = uint32(c[0] * 255 * c[3])
	g = uint32(c[1] * 255 * c[3])
	b = uint32(c[2] * 255 * c[3])
	a = uint32(c[3] * 255)
	return
}

// RGBA8 returns the 8-bit RGBA values.
func (c ColorF32) RGBA8() [4]uint8 {
	return [4]byte{
		uint8(c[0] * 255),
		uint8(c[1] * 255),
		uint8(c[2] * 255),
		uint8(c[3] * 255),
	}
}

// ToHSV returns the color as an HSV value.
func (c ColorF32) ToHSV() HSV {
	ret := HSV{}
	min := util.MinFloat32(c[0], c[1], c[2])
	max := util.MaxFloat32(c[0], c[1], c[2])
	d := max - min
	if c[0] >= c[1] && c[0] >= c[2] {
		ret[0] = (c[1] - c[2]) / d
	} else if c[1] >= c[0] && c[1] >= c[2] {
		ret[0] = (c[2]-c[0])/d + 2
	} else {
		ret[0] = (c[0]-c[1])/d + 4
	}
	ret[0] *= 60
	ret[0] = float32(math.Mod(float64(ret[0]), 360))
	ret[0] /= 360
	if max == 0 {
		ret[1] = 0
	} else {
		ret[1] = d / max
	}
	ret[2] = max
	return ret
}

// ToHex returns the color's RGBA8 values as a hex string.
func (c ColorF32) ToHex() string {
	ret := []byte{}
	rgba := c.RGBA8()
	ret = hex.AppendEncode(ret, rgba[:])
	return string(ret)
}

// ToHexRGB returns the color's RGB8 values as a hex string.
func (c ColorF32) ToHexRGB() string {
	return c.ToHex()[:6]
}

// ToColorRGBA returns the color as a color.RGBA value.
func (c ColorF32) ToColorRGBA() color.RGBA {
	rgba := c.RGBA8()
	return color.RGBA{
		R: rgba[0],
		G: rgba[1],
		B: rgba[2],
		A: rgba[3],
	}
}
