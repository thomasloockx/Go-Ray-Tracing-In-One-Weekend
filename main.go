package main

import (
    cgm "raytracer/cgmath"
    "fmt"
    "math"
    "os"
)

const IMAGE_WIDTH = 256
const IMAGE_HEIGHT = 256

func hitSphere(center *cgm.Vec3, radius float64, ray *cgm.Ray) float64 {
    oc := ray.Orig.Sub(center)    
    a := ray.Dir.LengthSquared()
    halfB := oc.Dot(&ray.Dir)
    c := oc.LengthSquared() - radius * radius
    discriminant := halfB * halfB - a * c
    if discriminant < 0 {
        return -1
    }
    return (-halfB - math.Sqrt(discriminant)) / a
}

func rayColor(r *cgm.Ray) *cgm.Color {
    t := hitSphere(&cgm.Vec3{0, 0, -1}, 0.5, r)
    if t > 0.0 {
        n := r.At(t).Sub(&cgm.Vec3{0, 0, -1}).UnitVector()
        return (&cgm.Color{n.X + 1, n.Y + 1, n.Z + 1}).Scale(0.5)
    }
    unitDir := r.Dir.UnitVector()
    t = 0.5 * (unitDir.Y + 1.0)
    white := &cgm.Color{R: 1.0, G: 1.0, B: 1.0}
    blue := &cgm.Color{R: 0.5, G: 0.7, B: 1.0}
    return white.Lerp(blue, t)
}


func main() {
    // Image
    aspectRatio := 16.0 / 9.0
    imageWidth := 400 
    imageHeight := int(float64(imageWidth) / aspectRatio)

    // Camera
    viewPortHeight := 2.0
    viewPortWidth := aspectRatio * viewPortHeight
    focalLength := 1.0

    origin := &cgm.Vec3{X: 0, Y: 0, Z: 0}
    horizontal := &cgm.Vec3{X: viewPortWidth, Y: 0, Z: 0}
    vertical := &cgm.Vec3{X: 0, Y: viewPortHeight, Z: 0}
    lowerLeftCorner := origin.Sub(horizontal.Div(2.0)).Sub(vertical.Div(2.0)).Sub(&cgm.Vec3{X: 0, Y: 0, Z: focalLength})

    // Render
    fmt.Printf("P3\n") 
    fmt.Printf("%d %d\n", imageWidth, imageHeight)
    fmt.Printf("255\n") 

    for j := imageHeight - 1; j >= 0; j-- {
        fmt.Fprintf(os.Stderr, "\rScanlines remaining: %d", j)
        for i := 0; i < imageWidth; i++ {
            u := float64(i) / float64(imageWidth - 1)
            v := float64(j) / float64(imageHeight - 1)
            r := &cgm.Ray{
                Orig: *origin, 
                Dir: *(lowerLeftCorner.Add(horizontal.Scale(u)).Add(vertical.Scale(v)).Sub(origin)),
            }
            c := rayColor(r)
            cgm.WriteColor(os.Stdout, c)
        }
    }
    fmt.Fprintf(os.Stderr, "\nDone.\n")
}
