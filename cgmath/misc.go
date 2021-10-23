package cgmath

import (
    "math"
    "math/rand"
)

func Lerp(x, y, t float64) float64 {
    return (1.0 - t) * x + t * y
}

func DegToRad(degrees float64) float64 {
    return degrees * math.Pi / 180.0
}

func Rand() float64 {
    return rand.Float64()
}

func RandInRange(min float64, max float64) float64 {
    return min + (max - min) * rand.Float64()
}

func Clamp(x, min, max float64) float64 {
    if x < min {
        return min
    }
    if x > max {
        return max
    }
    return x
}
