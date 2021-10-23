package cgmath

func Lerp(x, y, t float64) float64 {
    return (1.0 - t) * x + t * y
}
