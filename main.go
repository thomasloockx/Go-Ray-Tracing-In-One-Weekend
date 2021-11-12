package main

import (
    cgm "raytracer/cgmath"
    "fmt"
    "math"
    "os"
    "time"
)

// Avoid self intersection (shadow acne) by offsetting the ray position.
const RAY_EPSILON = 0.001

func randomScene() cgm.Hittable {
    world := &cgm.HittableList{}

    groundMaterial := cgm.Lambertian{Albedo: 
        cgm.MakeCheckerTexture(
            cgm.MakeSolidColor(0.2, 0.3, 0.1),
            cgm.MakeSolidColor(0.9, 0.9, 0.9),
        ),
    }

    world.Add(&cgm.Sphere{cgm.Vec3{0, -1000, 0}, 1000, &groundMaterial})

    for a := -11; a < 11; a++ {
        for b := -11; b < 11; b++ {
			chooseMat := cgm.Rand()
            center := cgm.Vec3{float64(a) + 0.9 * cgm.Rand(), 0.2, float64(b) + 0.9 * cgm.Rand()}

            if center.Sub(&cgm.Vec3{4, 0.2, 0}).Length() > 0.9 {
                var sphereMaterial cgm.Material

                if chooseMat < 0.8 {
                    // diffuse
                    albedo := cgm.MakeSolidColor(cgm.Rand() * cgm.Rand(), cgm.Rand() * cgm.Rand(), cgm.Rand() * cgm.Rand())
                    sphereMaterial = &cgm.Lambertian{Albedo: albedo}
                    center2 := center.Add(&cgm.Vec3{0.0, cgm.RandInRange(0, 0.5), 0})
                    sphere := &cgm.MovingSphere{
                        Center0: center,
                        Center1: *center2,
                        Time0: 0.0,
                        Time1: 1.0,
                        Radius: 0.2,
                        Material: sphereMaterial,
                    }
                    world.Add(sphere)
                } else if chooseMat < 0.95 {
                    // metal
					x := cgm.RandInRange(0.5, 1)
                    albedo := cgm.Color{x, x, x}
                    fuzz := cgm.RandInRange(0, 0.5)
                    sphereMaterial = &cgm.Metal{Albedo: albedo, Fuzz: fuzz}
				    world.Add(&cgm.Sphere{center, 0.2, sphereMaterial})
                } else {
                    // glass
                    sphereMaterial = &cgm.Dielectric{RefractiveIndex: 1.5}
				    world.Add(&cgm.Sphere{center, 0.2, sphereMaterial})
                }
            }
        }
    }

    material1 := cgm.Dielectric{RefractiveIndex: 1.5}
    world.Add(&cgm.Sphere{cgm.Vec3{0, 1, 0}, 1.0, &material1})

    material2 := cgm.Lambertian{cgm.MakeSolidColor(0.4, 0.2, 0.1)}
    world.Add(&cgm.Sphere{cgm.Vec3{-4, 1, 0}, 1.0, &material2})

    material3 := cgm.Metal{cgm.Color{0.7, 0.6, 0.5}, 0.0}
    world.Add(&cgm.Sphere{cgm.Vec3{4, 1, 0}, 1.0, &material3})

    return world
}

func twoSpheres() *cgm.HittableList {
    checker := cgm.MakeCheckerTexture(
        cgm.MakeSolidColor(0.2, 0.3, 0.1),
        cgm.MakeSolidColor(0.9, 0.9, 0.9),
    )

    objects := &cgm.HittableList{}
    objects.Add(&cgm.Sphere{cgm.Vec3{0, -10, 0}, 10, &cgm.Lambertian{checker}})
    objects.Add(&cgm.Sphere{cgm.Vec3{0, 10, 0}, 10, &cgm.Lambertian{checker}})

    return objects
}

func twoPerlinSpheres() *cgm.HittableList {
    perlinTexture := cgm.MakeNoiseTexture(4)     
    objects := &cgm.HittableList{}
    objects.Add(&cgm.Sphere{cgm.Vec3{0, -1000, 0}, 1000, &cgm.Lambertian{perlinTexture}})
    objects.Add(&cgm.Sphere{cgm.Vec3{0, 2, 0}, 2, &cgm.Lambertian{perlinTexture}})
    return objects
}

func earth() *cgm.HittableList {
    earthTexture := cgm.MakeImageTexture("earthmap.jpg")
    earthSuface := cgm.Lambertian{earthTexture}
    objects := &cgm.HittableList{}
    objects.Add(&cgm.Sphere{cgm.Vec3{0, 0, 0}, 2, &earthSuface})
    return objects
}

func rayColor(r *cgm.Ray, world cgm.Hittable, depth int) *cgm.Color {
    // If we exceeded the ray bounce limit, no more light is gathered.
    if depth <= 0 {
        return &cgm.Color{R: 0, G: 0, B: 0}
    }

    var rec cgm.HitRecord
    if world.Hit(r, RAY_EPSILON, math.Inf(1), &rec) {
        var scattered cgm.Ray
        var attenuation cgm.Color
        if rec.Material.Scatter(r, &rec, &attenuation, &scattered) {
            return attenuation.Mul(rayColor(&scattered, world, depth - 1))
        }

        return &cgm.Color{R: 0, G: 0, B: 0}
    }

    // Return the sky color.
    unitDir := r.Dir.UnitVector()
    t := 0.5 * (unitDir.Y + 1.0)
    white := &cgm.Color{R: 1.0, G: 1.0, B: 1.0}
    blue := &cgm.Color{R: 0.5, G: 0.7, B: 1.0}
    return white.Lerp(blue, t)
}


func main() {
    startTime := time.Now()

    // Image
    aspectRatio := 16.0 / 9.0
    imageWidth := 400
    imageHeight := int(float64(imageWidth) / aspectRatio)
    samplesPerPixel := 100
    maxDepth := 50

    // World
    var world cgm.Hittable
    var lookFrom cgm.Vec3
    var lookAt cgm.Vec3
    vfov := 40.0
    aperture := 0.0

    scene := 3
    switch scene {
        case 0:
            world = randomScene()
            lookFrom = cgm.Vec3{13, 2, 3}
            lookAt = cgm.Vec3{0, 0, 0}
            vfov = 20.0
            aperture = 0.1
        case 1:
            world = twoSpheres()
            lookFrom = cgm.Vec3{X: 13, Y: 2, Z: 3}
            lookAt = cgm.Vec3{X: 0, Y: 0, Z: 0}
            vfov = 20.0
        case 2:
            world = twoPerlinSpheres()
            lookFrom = cgm.Vec3{13, 2, 3}
            lookAt = cgm.Vec3{X: 0, Y: 0, Z: 0}
            vfov = 20.0
        case 3:
            world = earth()
            lookFrom = cgm.Vec3{13, 2, 3}
            lookAt = cgm.Vec3{X: 0, Y: 0, Z: 0}
            vfov = 20.0
    }

    vUp := cgm.Vec3{X: 0, Y: 1, Z: 0}
    distToFocus := 10.0
    cam := cgm.MakeCamera(&lookFrom, &lookAt, &vUp, vfov, aspectRatio, aperture, distToFocus, 0.0, 1.0)
    bvh := cgm.MakeBvh([]cgm.Hittable{world}, 0.0, 1.0)

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
                pixelColor.Accumulate(rayColor(&r, bvh, maxDepth))
            }
            cgm.WriteColor(os.Stdout, &pixelColor, samplesPerPixel)
        }
    }
    fmt.Fprintf(os.Stderr, "\nDone.\n")

    renderDuration := time.Since(startTime)
    fmt.Fprintf(os.Stderr, "Render time %v\n", renderDuration)
}
