package cgmath

import (
    "math"
)

type HitRecord struct {
    P Vec3
    Normal Vec3
    T float64
    FrontFace bool
    Material Material
}

// Set the normal so that it always points opposite the incident ray.
func (h *HitRecord) SetFaceNormal(r *Ray, outwardNormal *Vec3) {
    h.FrontFace = r.Dir.Dot(outwardNormal) < 0
    if h.FrontFace {
        h.Normal = *outwardNormal
    } else {
        h.Normal = *outwardNormal.Negate()
    }
}

type Hittable interface {
    Hit(r *Ray, tMin float64, tMax float64, h *HitRecord) bool
}

type Sphere struct {
    Center Vec3
    Radius float64
    Material Material
}

func (s *Sphere) Hit(r *Ray, tMin float64, tMax float64, rec *HitRecord) bool {
    oc := r.Orig.Sub(&s.Center)
    a := r.Dir.LengthSquared()
    halfB := oc.Dot(&r.Dir)
    c := oc.LengthSquared() - s.Radius * s.Radius

    discriminant := halfB * halfB - a * c
    if discriminant < 0.0 {
        return false
    }

    sqrtd := math.Sqrt(discriminant)

    root := (-halfB - sqrtd) / a
    if (root < tMin || tMax < root) {
        root = (-halfB + sqrtd) / a
        if (root < tMin || tMax < root) {
            return false
        }
    }

    rec.T = root
    rec.P = *r.At(rec.T)
    outwardNormal := rec.P.Sub(&s.Center).Div(s.Radius)
    rec.SetFaceNormal(r, outwardNormal)
    rec.Material = s.Material
    return true
}

type HittableList struct {
    objects []Hittable
}

func (hl *HittableList) Add(h Hittable) {
    hl.objects = append(hl.objects, h)
}

func (hl *HittableList) Clear() {
    hl.objects = make([]Hittable, 0, 16)
}

func (hl *HittableList) Hit(r *Ray, tMin float64, tMax float64, h *HitRecord) bool {
    var tmpRecord HitRecord
    hitAnything := false
    closestSoFar := tMax

    for _, object := range hl.objects {
        if (object.Hit(r, tMin, closestSoFar, &tmpRecord)) {
            hitAnything = true
            closestSoFar = tmpRecord.T
            *h = tmpRecord
        }
    }

    return hitAnything
}
