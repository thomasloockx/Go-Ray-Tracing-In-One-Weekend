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
    *scattered = Ray{Orig: rec.P, Dir: *scatterDir, Time: rayIn.Time}
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
    *scattered = Ray{Orig: rec.P, Dir: *reflected, Time: rayIn.Time}
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
    cosTheta := math.Min(-unitDirection.Dot(&rec.Normal), 1.0)
    sinTheta := math.Sqrt(1.0 - cosTheta * cosTheta)

    cannotRefract := refractionRatio * sinTheta > 1.0

    var direction *Vec3
    if cannotRefract || reflectance(cosTheta, refractionRatio) > Rand() {
        direction = Reflect(unitDirection, &rec.Normal)
    } else {
        direction = Refract(unitDirection, &rec.Normal, refractionRatio)
    }

    *scattered = Ray{Orig: rec.P, Dir: *direction, Time: rayIn.Time}
    return true
}

// Schlick approximation of Fresnel equations.
func reflectance(cosine, refractiveIdx float64) float64 {
    r0 := (1 - refractiveIdx) / (1 + refractiveIdx)
    r0 = r0 * r0
    return r0 + (1 - r0) * math.Pow(1 - cosine, 5)
}
