package cgmath

import (
    "math"
)

type Vec3 struct {
    X, Y, Z float64
}

func (v *Vec3) Negate() *Vec3 {
    return &Vec3{-v.X, -v.Y, -v.Z}
}

func (v *Vec3) Length() float64 {
    return math.Sqrt(v.LengthSquared())
}

func (v *Vec3) LengthSquared() float64 {
    return v.X * v.X + v.Y * v.Y + v.Z * v.Z
}

func (v *Vec3) Add(w *Vec3) *Vec3 {
    return &Vec3{v.X + w.X, v.Y + w.Y, v.Z + w.Z}
}

func (v *Vec3) Sub(w *Vec3) *Vec3 {
    return &Vec3{v.X - w.X, v.Y - w.Y, v.Z - w.Z}
}

func (v *Vec3) Mul(w *Vec3) *Vec3 {
    return &Vec3{v.X * w.X, v.Y * w.Y, v.Z * w.Z}
}

func (v *Vec3) Scale(t float64) *Vec3 {
    return &Vec3{v.X * t, v.Y * t, v.Z * t}
}

func (v *Vec3) Div(t float64) *Vec3 {
    return v.Scale(1.0 / t)
}

func (v *Vec3) Dot(w *Vec3) float64 {
    return v.X * w.X + v.Y * w.Y + v.Z * w.Z
}

func (v *Vec3) Cross(w *Vec3) *Vec3 {
    return &Vec3{
        v.Y * w.Z - v.Z * w.Y,
        v.Z * w.X - v.X * w.Z,
        v.X * w.Y - v.Y * w.X,
    }
}

func (v *Vec3) UnitVector() *Vec3 {
    return v.Div(v.Length())
}

func Random() *Vec3 {
    return &Vec3{Rand(), Rand(), Rand()}
}

func RandomInRange(min, max float64) *Vec3 {
    return &Vec3{RandInRange(min, max), RandInRange(min, max), RandInRange(min, max)}
}

func RandomInUnitSphere() *Vec3 {
    for {
        p := RandomInRange(-1.0, 1.0)
        if (p.LengthSquared() < 1) {
            return p
        }
    }
}

func RandomUnitVector() *Vec3 {
    return RandomInUnitSphere().UnitVector()
}

func RandomInHemisphere(normal *Vec3) *Vec3 {
    inUnitSphere := RandomInUnitSphere()
    if inUnitSphere.Dot(normal) > 0.0 {
        return inUnitSphere
    }
    return inUnitSphere.Negate()
}

func (v *Vec3) NearZero() bool {
    s := 1e-8
    return math.Abs(v.X) < s && math.Abs(v.Y) < s && math.Abs(v.Z) < s
}

func Reflect(v, n *Vec3) *Vec3 {
    return v.Sub(n.Scale(2 * v.Dot(n)))
}
