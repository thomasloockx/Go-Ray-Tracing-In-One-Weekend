package cgmath

import (
    "fmt"
    "io"
    "math"
)

type Color struct {
    R, G, B float64
}

func (c *Color) Lerp(d *Color, t float64) *Color {
    return &Color{
        R: Lerp(c.R, d.R, t),
        G: Lerp(c.G, d.G, t),
        B: Lerp(c.B, d.B, t),
    }
}

func (c *Color) Add(d *Color) *Color {
    return &Color{c.R +  d.R, c.G + d.G, c.B + d.B}
}

func (c *Color) Scale(t float64) *Color {
    return &Color{c.R * t, c.G * t, c.B * t}
}

func (c *Color) Mul(d *Color) *Color {
    return &Color{c.R * d.R, c.G * d.G, c.B * d.B}
}

func (c *Color) Accumulate(d *Color) {
    c.R += d.R
    c.G += d.G
    c.B += d.B
}

func WriteColor(w io.Writer, c *Color, samplesPerPixel int) {
    // Divide the color by the number of samples and gamma-correct for gamma=2.0
    scale := 1.0 / float64(samplesPerPixel)
    r := math.Sqrt(c.R * scale)
    g := math.Sqrt(c.G * scale)
    b := math.Sqrt(c.B * scale)


    fmt.Fprintf(w, "%d %d %d\n", 
        int(256 * Clamp(r, 0.0, 0.999)),
        int(256 * Clamp(g, 0.0, 0.999)),
        int(256 * Clamp(b, 0.0, 0.999)))
}
