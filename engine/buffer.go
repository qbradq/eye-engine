package engine

import (
	"errors"
	"image"
	"io"

	_ "image/png"
)

// Buffer is a rectangular buffer of colors.
type Buffer struct {
	Width   int          // Width of the buffer in pixels
	Height  int          // Height of the buffer in pixels
	Pixels  []ColorIndex // Pixel buffer
	Palette *Palette     // Palette used by the buffer
}

// NewBuffer creates a new buffer with the given size.
func NewBuffer(w, h int, p *Palette) *Buffer {
	ret := &Buffer{
		Width:   w,
		Height:  h,
		Pixels:  make([]ColorIndex, w*h),
		Palette: p,
	}
	ret.Clear()
	return ret
}

// DecodeBuffer decodes an image from the reader and returns it as a buffer.
func DecodeBuffer(r io.Reader, p *Palette) (*Buffer, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}
	max := img.Bounds().Max
	buf := NewBuffer(max.X, max.Y, p)
	for y := range max.Y {
		for x := range max.X {
			c := img.At(x, y)
			i := buf.Palette.IndexOf(c)
			buf.Pixels[y*max.X+x] = i
		}
	}
	return buf, nil
}

// WriteTo writes the bytes of the buffer as RGBA bytes.
func (b *Buffer) Bytes(dest []byte) error {
	if b.Width*b.Height*4 != len(dest) {
		return errors.New("buffer size mismatch")
	}
	for i, c := range b.Pixels {
		if c == ColorIndexTransparent {
			dest[i*4+0] = 0
			dest[i*4+1] = 0
			dest[i*4+2] = 0
			dest[i*4+3] = 0
		} else {
			rgba := b.Palette[c]
			dest[i*4+0] = rgba.R
			dest[i*4+1] = rgba.G
			dest[i*4+2] = rgba.B
			dest[i*4+3] = 255
		}
	}
	return nil
}

// Fill fills the buffer with the color.
func (b *Buffer) Fill(c ColorIndex) {
	for i := range b.Pixels {
		b.Pixels[i] = c
	}
}

// Clear clears the buffer to opaque black.
func (b *Buffer) Clear() {
	b.Fill(ColorIndexTransparent)
}

// SetPixel sets a single pixel.
func (b *Buffer) SetPixel(x, y int, c ColorIndex) {
	if x < 0 || x >= b.Width || y < 0 || y >= b.Height {
		return
	}
	b.Pixels[y*b.Width+x] = c
}

// Draw copies all non-transparent pixels from src into this buffer.
func (b *Buffer) Draw(src *Buffer, dx, dy, sx, sy, w, h int) {
	for oy := range h {
		idy := dy + oy
		if idy < 0 || idy >= b.Height {
			continue
		}
		isy := sy + oy
		if isy < 0 || isy >= src.Height {
			continue
		}
		for ox := range w {
			idx := dx + ox
			if idx < 0 || idx >= b.Width {
				continue
			}
			isx := sx + ox
			if isx < 0 || isx >= src.Width {
				continue
			}
			c := src.Pixels[isy*src.Width+isx]
			if c == ColorIndexTransparent {
				continue
			}
			b.Pixels[idy*b.Width+idx] = c
		}
	}
}

// DrawMask places pixel p at every location a non-transparent pixel appears in
// src. This is useful for rendering fonts and masks for example.
func (b *Buffer) DrawMask(src *Buffer, dx, dy, sx, sy, w, h int, p ColorIndex) {
	for oy := range h {
		idy := dy + oy
		if idy < 0 || idy >= b.Height {
			continue
		}
		isy := sy + oy
		if isy < 0 || isy >= src.Height {
			continue
		}
		for ox := range w {
			idx := dx + ox
			if idx < 0 || idx >= b.Width {
				continue
			}
			isx := sx + ox
			if isx < 0 || isx >= src.Width {
				continue
			}
			c := src.Pixels[isy*src.Width+isx]
			if c == ColorIndexTransparent {
				continue
			}
			b.Pixels[idy*b.Width+idx] = p
		}
	}
}
