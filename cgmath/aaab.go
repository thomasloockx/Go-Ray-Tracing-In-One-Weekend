package cgmath

import (
    "fmt"
    "math"
)

type Aabb struct {
    Minimum, Maximum Vec3
}

func (a *Aabb) String() string {
    return fmt.Sprintf("Aabb(min=%v, max=%v)", a.Minimum, a.Maximum)
}

func (a *Aabb) Hit(r *Ray, tMin float64, tMax float64) bool {
    var invD, t0, t1 float64

    // X slab
    invD = 1.0 / r.Dir.X 
    t0 = (a.Minimum.X - r.Orig.X) * invD
    t1 = (a.Maximum.X - r.Orig.X) * invD
    if invD < 0 {
        t0, t1 = t1, t0
    }

    if t0 > tMin {
        tMin = t0
    }

    if t1 < tMax {
        tMax = t1
    }

    if (tMax <= tMin) {
        return false
    }

    // Y slab
    invD = 1.0 / r.Dir.Y
    t0 = (a.Minimum.Y - r.Orig.Y) * invD
    t1 = (a.Maximum.Y - r.Orig.Y) * invD
    if invD < 0 {
        t0, t1 = t1, t0
    }

    if t0 > tMin {
        tMin = t0
    }

    if t1 < tMax {
        tMax = t1
    }

    if (tMax <= tMin) {
        return false
    }

    // Z slab
    invD = 1.0 / r.Dir.Z 
    t0 = (a.Minimum.Z - r.Orig.Z) * invD
    t1 = (a.Maximum.Z - r.Orig.Z) * invD
    if invD < 0 {
        t0, t1 = t1, t0
    }

    if t0 > tMin {
        tMin = t0
    }

    if t1 < tMax {
        tMax = t1
    }

    if (tMax <= tMin) {
        return false
    }

    return true
}

func surroundingBox(box0 *Aabb, box1 *Aabb) *Aabb {
    small := Vec3{
        X: math.Min(box0.Minimum.X, box1.Minimum.X),
        Y: math.Min(box0.Minimum.Y, box1.Minimum.Y),
        Z: math.Min(box0.Minimum.Z, box1.Minimum.Z),
    }
    big := Vec3{
        X: math.Max(box0.Maximum.X, box1.Maximum.X),
        Y: math.Max(box0.Maximum.Y, box1.Maximum.Y),
        Z: math.Max(box0.Maximum.Z, box1.Maximum.Z),
    }
    return &Aabb{Minimum: small, Maximum: big}
}
