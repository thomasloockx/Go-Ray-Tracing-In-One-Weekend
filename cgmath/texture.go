package cgmath

import (
    "math"
    "image"
    _ "image/jpeg"
    "log"
    "os"
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

type ImageTexture struct {
    image image.Image
    width int
    height int
}

func MakeImageTexture(imagePath string) *ImageTexture {
    reader, err := os.Open(imagePath)
    if err != nil {
        log.Fatal(err)
    }
    defer reader.Close()

    image, _, err := image.Decode(reader)
    bounds := image.Bounds()

    return &ImageTexture{
        image: image,
        width: bounds.Max.X - bounds.Min.X,
        height: bounds.Max.Y - bounds.Min.Y,
    }
}

func (t *ImageTexture) Value(u float64, v float64, p *Vec3) Color {
    u = Clamp(u, 0.0, 1.0)
    v = 1.0 - Clamp(v, 0.0, 1.0)

    i := int(u * float64(t.width))
    j := int(v * float64(t.height))

    r, g, b, _ := t.image.At(i, j).RGBA()

    return Color{
        R: float64(r) / 0xffff,
        G: float64(g) / 0xffff,
        B: float64(b) / 0xffff,
    }
}
