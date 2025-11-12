package engine

import (
	"bufio"
	"fmt"
	"image/color"
	"path/filepath"
	"strings"

	"github.com/qbradq/eye-engine/data"
)

// Palette represents the palette of the engine.
type Palette [256]color.RGBA

// GeneratePalette returns a pointer to a new Palette.
func GeneratePalette() *Palette {
	ret := Palette{}
	for r := range 8 {
		for g := range 4 {
			for b := range 8 {
				idx := (r << 5) | (g << 3) | b
				c := ColorF32{
					float32(r) / 7.0,
					float32(g) / 3.0,
					float32(b) / 7.0,
					1.0,
				}
				ret[idx] = c.ToColorRGBA()
			}
		}
	}
	return &ret
}

// LoadPalette returns a new Palette based on the named palette. Pass "" or
// "default" for the default palette.
func LoadPalette(name string) (*Palette, error) {
	// Detect default case
	if name == "" {
		name = "default"
	}
	// Open the palette file
	f, err := data.FS.Open(filepath.Join("palettes", name+".gpl"))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	// Load the base palette
	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		return nil, fmt.Errorf("empty palette file '%s'", name)
	}
	line := strings.TrimSpace(scanner.Text())
	if line != "GIMP Palette" {
		return nil, fmt.Errorf("palette file '%s' missing header", name)
	}
	baseColors := []ColorF32{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 4 {
			return nil, fmt.Errorf(
				"palette file '%s' malformed color line",
				name,
			)
		}
		baseColors = append(baseColors, NewColorFromHex(parts[3]))
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning palette file '%s': %w", name, err)
	}
	if len(baseColors) != 16 {
		return nil, fmt.Errorf("palette file '%s' should have 16 entries", name)
	}
	// Generate the full palette
	ret := Palette{}
	{
		r0 := float32(2.6)
		for iColor, base := range baseColors {
			for iDistanceFraction := range 16 {
				k := float32(iDistanceFraction)
				c := r0 / (r0 + k)
				c *= c
				hsv := base.ToHSV()
				hsv[2] = hsv[2] * c
				ret[iColor+iDistanceFraction*16] = NewColorFromHSV(hsv, 1.0).ToColorRGBA()
			}
		}
	}
	return &ret, nil
}

// IndexOf returns the ColorIndex that best matches the input color.
func (p *Palette) IndexOf(c color.Color) ColorIndex {
	_, _, _, a := c.RGBA()
	if a < 255 {
		return ColorIndexTransparent
	}
	tp := make(color.Palette, 16)
	for i := range 16 {
		tp[i] = p[256+i]
	}
	return ColorIndex(tp.Index(c))
}
