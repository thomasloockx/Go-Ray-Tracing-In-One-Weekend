package main

import (
    cgm "raytracer/cgmath"
    "fmt"
    "math"
    "os"
)

func rayColor(r *cgm.Ray, world cgm.Hittable) *cgm.Color {
    var rec cgm.HitRecord
    if world.Hit(r, 0, math.Inf(1), &rec) {
        n := rec.Normal
        return (&cgm.Color{R: n.X + 1, G: n.Y + 1, B: n.Z + 1}).Scale(0.5)
    }

    // Return the sky color.
    unitDir := r.Dir.UnitVector()
    t := 0.5 * (unitDir.Y + 1.0)
    white := &cgm.Color{R: 1.0, G: 1.0, B: 1.0}
    blue := &cgm.Color{R: 0.5, G: 0.7, B: 1.0}
    return white.Lerp(blue, t)
}


func main() {
    // Image
    aspectRatio := 16.0 / 9.0
    imageWidth := 400 
    imageHeight := int(float64(imageWidth) / aspectRatio)

    // World
    world := cgm.HittableList{}
    s1 := &cgm.Sphere{Center: cgm.Vec3{0, 0, -1}, Radius: 0.5}
    s2 := &cgm.Sphere{Center: cgm.Vec3{0, -100.5, -1}, Radius: 100}
    world.Add(s1)
    world.Add(s2)

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
            c := rayColor(r, &world)
            cgm.WriteColor(os.Stdout, c)
        }
    }
    fmt.Fprintf(os.Stderr, "\nDone.\n")
}
