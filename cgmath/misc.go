package cgmath

import (
    "math"
)

func Lerp(x, y, t float64) float64 {
    return (1.0 - t) * x + t * y
}

func DegToRad(degrees float64) float64 {
    return degrees * math.Pi / 180.0
}
