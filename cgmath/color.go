package cgmath

import (
    "fmt"
    "io"
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

func (c *Color) Scale(t float64) *Color {
    return &Color{c.R * t, c.G * t, c.B * t}
}


func WriteColor(w io.Writer, c *Color) {
    fmt.Fprintf(w, "%d %d %d\n", int(255.999 * c.R), int(255.999 * c.G), int(255.999 * c.B))
}
