package cgmath

import (
    "math"
)

type Camera struct {
    origin Vec3
    lowerLeftCorner Vec3
    horizontal Vec3
    vertical Vec3
    u, v, w Vec3
    lensRadius float64
}

func MakeCamera(lookFrom *Vec3, lookAt *Vec3, vUp *Vec3, vFov float64, aspectRatio float64, aperture float64, focusDist float64) Camera {
    cam := Camera{}

    theta := DegToRad(vFov)
    h := math.Tan(theta / 2.0)
    viewportHeight := 2.0 * h
    viewportWidth := aspectRatio * viewportHeight

    w := lookFrom.Sub(lookAt).UnitVector()
    u := vUp.Cross(w).UnitVector()
    v := w.Cross(u)

    cam.origin = *lookFrom
    cam.horizontal = *(u.Scale(focusDist * viewportWidth))
    cam.vertical = *(v.Scale(focusDist * viewportHeight))
    cam.lowerLeftCorner = *(cam.origin.Sub(cam.horizontal.Div(2.0)).Sub(cam.vertical.Div(2.0)).Sub(w.Scale(focusDist)))
    cam.lensRadius = 0.5 * aperture
    cam.u = *u
    cam.v = *v
    cam.w = *w

    return cam
}

func (c *Camera) MakeRay(u, v float64) Ray {
    rd := RandomInUnitDisk().Scale(c.lensRadius)
    offset := c.u.Scale(rd.X).Add(c.v.Scale(rd.Y))

    return Ray{
        Orig: *c.origin.Add(offset),
        Dir: *(c.lowerLeftCorner.Add(c.horizontal.Scale(u)).Add(c.vertical.Scale(v)).Sub(&c.origin).Sub(offset)),
    }
}
