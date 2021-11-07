package cgmath

import (
    "math"
)

type Texture interface {
    Value(u float64, v float64, p *Vec3) Color
}

type SolidColor struct {
    color Color
}

func MakeSolidColor(r, g, b float64) *SolidColor {
    return &SolidColor{color: Color{r, g, b}}
}

func (s* SolidColor) Value(u float64, v float64, p *Vec3) Color {
    return s.color
}

type CheckerTexture struct {
    odd, even Texture
}

func MakeCheckerTexture(odd, even Texture) *CheckerTexture {
    return &CheckerTexture{odd: odd, even: even}
}

func (t *CheckerTexture) Value(u float64, v float64, p *Vec3) Color {
    sines := math.Sin(10 * p.X) * math.Sin(10 * p.Y) * math.Sin(10 * p.Z)
    if sines < 0 {
        return t.odd.Value(u, v, p)
    } else {
        return t.even.Value(u, v, p)        
    }
}
