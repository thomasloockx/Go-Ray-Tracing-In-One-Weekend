package cgmath

type Ray struct {
    Orig, Dir Vec3
}

func (r *Ray) At(t float64) *Vec3 {
    return r.Orig.Add(r.Dir.Scale(t))
}
