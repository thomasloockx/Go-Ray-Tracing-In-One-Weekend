package cgmath

type Camera struct {
    origin Vec3
    lowerLeftCorner Vec3
    horizontal Vec3
    vertical Vec3
}

func MakeCamera() Camera {
    cam := Camera{}

    aspectRation := 16.0 / 9.0
    viewportHeight := 2.0
    viewportWidth := aspectRation * viewportHeight
    focalLength := 1.0

    cam.origin = Vec3{0.0, 0.0, 0.0}
    cam.horizontal = Vec3{viewportWidth, 0.0, 0.0}
    cam.vertical = Vec3{0.0, viewportHeight, 0.0}
    cam.lowerLeftCorner = *(cam.origin.Sub(cam.horizontal.Div(2.0)).Sub(cam.vertical.Div(2.0)).Sub(&Vec3{X: 0, Y: 0, Z: focalLength}))

    return cam
}

func (c *Camera) MakeRay(u, v float64) Ray {
    return Ray{
        Orig: c.origin,
        Dir: *(c.lowerLeftCorner.Add(c.horizontal.Scale(u)).Add(c.vertical.Scale(v)).Sub(&c.origin)),
    }
}
