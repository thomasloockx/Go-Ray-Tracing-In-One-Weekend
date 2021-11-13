package cgmath

import (
    "fmt"
    "math"
)

type HitRecord struct {
    P Vec3
    Normal Vec3
    T, U, V float64
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
    BoundingBox(time0 float64, time1 float64, outputBox *Aabb) bool
    fmt.Stringer
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
    rec.U, rec.V = s.getUv(outwardNormal)
    rec.Material = s.Material
    return true
}

func (s *Sphere) BoundingBox(time0 float64, time1 float64, outputBox *Aabb) bool {
    *outputBox = Aabb{
        Minimum: *s.Center.Sub(&Vec3{s.Radius, s.Radius, s.Radius}),
        Maximum: *s.Center.Add(&Vec3{s.Radius, s.Radius, s.Radius}),
    }
    return true
}

func (s *Sphere) String() string {
    return fmt.Sprintf("Sphere(Radius=%02f, Center=%v)", s.Radius, s.Center)
}

func (s *Sphere) getUv(p *Vec3) (float64, float64) {
    theta := math.Acos(-p.Y)
    phi := math.Atan2(-p.Z, p.X) + math.Pi
    u := phi / (2 * math.Pi)
    v := theta / math.Pi
    return u, v
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


func (hl *HittableList) BoundingBox(time0 float64, time1 float64, outputBox *Aabb) bool {
    if len(hl.objects) == 0 {
        return false
    }

    var tempBox Aabb
    firstBox := true

    for _, object := range hl.objects {
        if !object.BoundingBox(time0, time1, &tempBox) {
            return false
        }
        if firstBox {
            *outputBox = tempBox
            firstBox = false
        } else {
            *outputBox = *surroundingBox(outputBox, &tempBox)
        }
    }

    return true
}

func (hl *HittableList) String() string {
    return fmt.Sprintf("HittableList(objects=%v)", hl.objects)
}

type MovingSphere struct {
    Center0, Center1 Vec3
    Time0, Time1 float64
    Radius float64
    Material Material
}

func (s *MovingSphere) Center(time float64) *Vec3 {
    frac := (time - s.Time0) / (s.Time1 - s.Time0)
    return s.Center0.Add(s.Center1.Sub(&s.Center0).Scale(frac))
}

func (s *MovingSphere) Hit(r *Ray, tMin float64, tMax float64, rec *HitRecord) bool {
    sCenter := s.Center(r.Time)
    oc := r.Orig.Sub(sCenter)
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
    outwardNormal := rec.P.Sub(sCenter).Div(s.Radius)
    rec.SetFaceNormal(r, outwardNormal)
    rec.Material = s.Material
    return true
}

func (s *MovingSphere) BoundingBox(time0 float64, time1 float64, outputBox *Aabb) bool {
    box0 := Aabb{
        Minimum: *s.Center0.Sub(&Vec3{s.Radius, s.Radius, s.Radius}),
        Maximum: *s.Center0.Add(&Vec3{s.Radius, s.Radius, s.Radius}),
    }
    box1 := Aabb{
        Minimum: *s.Center1.Sub(&Vec3{s.Radius, s.Radius, s.Radius}),
        Maximum: *s.Center1.Add(&Vec3{s.Radius, s.Radius, s.Radius}),
    }
    *outputBox = *surroundingBox(&box0, &box1)
    return true
}

func (s *MovingSphere) String() string {
    return fmt.Sprintf("MovingSphere(Radius=%02f, Center0=%v, Center1=%v)", s.Radius, s.Center0, s.Center1)
}

type XyRect struct {
    X0, X1, Y0, Y1 float64
    K float64
    Material Material
}

func (rect *XyRect) Hit(r *Ray, tMin float64, tMax float64, rec *HitRecord) bool {
    t := (rect.K - r.Orig.Z) / r.Dir.Z
    if t < tMin || t > tMax {
        return false
    }

    x := r.Orig.X + t * r.Dir.X
    if x < rect.X0 || x > rect.X1 {
        return false
    }

    y := r.Orig.Y + t * r.Dir.Y
    if y < rect.Y0 || y > rect.Y1 {
        return false
    }

    rec.U = (x - rect.X0) / (rect.X1 - rect.X0)
    rec.V = (y - rect.Y0) / (rect.Y1 - rect.Y0)
    rec.T = t
    rec.SetFaceNormal(r, &Vec3{0, 0, 1})
    rec.Material = rect.Material
    rec.P = *r.At(t)

    return true
}



func (r *XyRect) BoundingBox(time0 float64, time1 float64, outputBox *Aabb) bool {
    *outputBox = Aabb{
        Minimum: Vec3{r.X0, r.Y0, r.K - 0.0001},
        Maximum: Vec3{r.X1, r.Y1, r.K + 0.0001},
    }
    return true
}

func (rect *XyRect) String() string {
    return fmt.Sprintf("XyRect(x=[%02f, %02f], y=[%02f, %02f], k=%02f)", rect.X0, rect.X1, rect.Y0, rect.Y1, rect.K)
}

type XzRect struct {
    X0, X1, Z0, Z1 float64
    K float64
    Material Material
}

func (rect *XzRect) Hit(r *Ray, tMin float64, tMax float64, rec *HitRecord) bool {
    t := (rect.K - r.Orig.Y) / r.Dir.Y
    if t < tMin || t > tMax {
        return false
    }

    x := r.Orig.X + t * r.Dir.X
    if x < rect.X0 || x > rect.X1 {
        return false
    }

    z := r.Orig.Z + t * r.Dir.Z
    if z < rect.Z0 || z > rect.Z1 {
        return false
    }

    rec.U = (x - rect.X0) / (rect.X1 - rect.X0)
    rec.V = (z - rect.Z0) / (rect.Z1 - rect.Z0)
    rec.T = t
    rec.SetFaceNormal(r, &Vec3{0, 1, 0})
    rec.Material = rect.Material
    rec.P = *r.At(t)

    return true
}

func (r *XzRect) BoundingBox(time0 float64, time1 float64, outputBox *Aabb) bool {
    *outputBox = Aabb{
        Minimum: Vec3{r.X0, r.K - 0.0001, r.Z0},
        Maximum: Vec3{r.X1, r.K + 0.0001, r.Z1},
    }
    return true
}

func (rect *XzRect) String() string {
    return fmt.Sprintf("XzRect(x=[%02f, %02f], z=[%02f, %02f], k=%02f)", rect.X0, rect.X1, rect.Z0, rect.Z1, rect.K)
}

type YzRect struct {
    Y0, Y1, Z0, Z1 float64
    K float64
    Material Material
}

func (rect *YzRect) Hit(r *Ray, tMin float64, tMax float64, rec *HitRecord) bool {
    t := (rect.K - r.Orig.X) / r.Dir.X
    if t < tMin || t > tMax {
        return false
    }

    z := r.Orig.Z + t * r.Dir.Z
    if z < rect.Z0 || z > rect.Z1 {
        return false
    }

    y := r.Orig.Y + t * r.Dir.Y
    if y < rect.Y0 || y > rect.Y1 {
        return false
    }

    rec.U = (y - rect.Y0) / (rect.Y1 - rect.Y0)
    rec.V = (z - rect.Z0) / (rect.Z1 - rect.Z0)
    rec.T = t
    rec.SetFaceNormal(r, &Vec3{1, 0, 0})
    rec.Material = rect.Material
    rec.P = *r.At(t)

    return true
}



func (r *YzRect) BoundingBox(time0 float64, time1 float64, outputBox *Aabb) bool {
    *outputBox = Aabb{
        Minimum: Vec3{r.K - 0.0001, r.Y0, r.Z0},
        Maximum: Vec3{r.K + 0.0001, r.Y1, r.Z1},
    }
    return true
}

func (rect *YzRect) String() string {
    return fmt.Sprintf("YzRect(y=[%02f, %02f], z=[%02f, %02f], k=%02f)", rect.Y0, rect.Y1, rect.Z0, rect.Z1, rect.K)
}
