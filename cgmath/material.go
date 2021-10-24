package cgmath

import (
    "math"
)

type Material interface {
    Scatter(rayIn *Ray, rec *HitRecord, attenuation *Color, scattered *Ray) bool
}

type Lambertian struct {
    Albedo Color
}

func (mat *Lambertian) Scatter(rayIn *Ray, rec *HitRecord, attenuation *Color, scattered *Ray) bool {
    scatterDir := rec.Normal.Add(RandomUnitVector())
    if scatterDir.NearZero() {
        scatterDir = &rec.Normal
    }
    *scattered = Ray{Orig: rec.P, Dir: *scatterDir}
    *attenuation = mat.Albedo
    return true
}

type Metal struct {
    Albedo Color
    Fuzz float64
}

func (mat *Metal) Scatter(rayIn *Ray, rec *HitRecord, attenuation *Color, scattered *Ray) bool {
    fuzz := math.Min(mat.Fuzz, 1.0)
    reflected := Reflect(rayIn.Dir.UnitVector(), &rec.Normal)
    reflected = reflected.Add(RandomInUnitSphere().Scale(fuzz))
    *scattered = Ray{Orig: rec.P, Dir: *reflected}
    *attenuation = mat.Albedo
    return scattered.Dir.Dot(&rec.Normal) > 0
}

type Dielectric struct {
    RefractiveIndex float64       
}

func (mat *Dielectric) Scatter(rayIn *Ray, rec *HitRecord, attenuation *Color, scattered *Ray) bool {
    *attenuation = Color{1.0, 1.0, 1.0}
    refractionRatio := mat.RefractiveIndex
    if rec.FrontFace {
        refractionRatio = 1.0 / mat.RefractiveIndex
    }
    unitDirection := rayIn.Dir.UnitVector()
    refracted := Refract(unitDirection, &rec.Normal, refractionRatio)
    *scattered = Ray{Orig: rec.P, Dir: *refracted}
    return true
}
