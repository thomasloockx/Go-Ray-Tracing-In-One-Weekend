package main

import (
    cgm "raytracer/cgmath"
    "fmt"
    "math"
    "os"
)

// Avoid self intersection (shadow acne) by offsetting the ray position.
const RAY_EPSILON = 0.001


func rayColor(r *cgm.Ray, world cgm.Hittable, depth int) *cgm.Color {
    // If we exceeded the ray bounce limit, no more light is gathered.
    if depth <= 0 {
        return &cgm.Color{R: 0, G: 0, B: 0}
    }

    var rec cgm.HitRecord
    if world.Hit(r, RAY_EPSILON, math.Inf(1), &rec) {
        target := rec.P.Add(&rec.Normal).Add(cgm.RandomInHemisphere(&rec.Normal))
        newRay := cgm.Ray{
            Orig: rec.P,
            Dir: *target.Sub(&rec.P),
        }
        // Gray
        return rayColor(&newRay, world, depth - 1).Scale(0.5)
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
    samplesPerPixel := 100
    maxDepth := 50

    // World
    world := cgm.HittableList{}
    s1 := &cgm.Sphere{Center: cgm.Vec3{0, 0, -1}, Radius: 0.5}
    s2 := &cgm.Sphere{Center: cgm.Vec3{0, -100.5, -1}, Radius: 100}
    world.Add(s1)
    world.Add(s2)

    // Camera
    cam := cgm.MakeCamera()

    // Render
    fmt.Printf("P3\n") 
    fmt.Printf("%d %d\n", imageWidth, imageHeight)
    fmt.Printf("255\n") 

    for j := imageHeight - 1; j >= 0; j-- {
        fmt.Fprintf(os.Stderr, "\rScanlines remaining: %d", j)
        for i := 0; i < imageWidth; i++ {
            pixelColor := cgm.Color{}
            for s := 0; s < samplesPerPixel; s++ {
                u := (float64(i) + cgm.Rand()) / float64(imageWidth - 1)
                v := (float64(j) + cgm.Rand()) / float64(imageHeight - 1)
                r := cam.MakeRay(u, v)
                pixelColor.Accumulate(rayColor(&r, &world, maxDepth))
            }
            cgm.WriteColor(os.Stdout, &pixelColor, samplesPerPixel)
        }
    }
    fmt.Fprintf(os.Stderr, "\nDone.\n")
}
