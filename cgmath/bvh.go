package cgmath

import (
    "fmt"
    "sort"
)

type BvhNode struct {
    left, right Hittable
    box Aabb
}

func (n *BvhNode) BoundingBox(time0 float64, time1 float64, outputBox *Aabb) bool {
    *outputBox = n.box
    return true
}

func (n *BvhNode) Hit(r *Ray, tMin float64, tMax float64, rec *HitRecord) bool {
    if !n.box.Hit(r, tMin, tMax) {
        return false
    }

    hitLeft := n.left.Hit(r, tMin, tMax, rec)
    if hitLeft {
        tMax = rec.T
    }
    hitRight := n.right.Hit(r, tMin, tMax, rec)

    return hitLeft || hitRight
}

func (n *BvhNode) String() string {
    return fmt.Sprintf("BvhNode(left=%v, right=%v, box=%v)", n.left, n.right, n.box)
}

func boxXCompare(a, b Hittable) bool {
    var boxA, boxB Aabb
    a.BoundingBox(0, 0, &boxA)
    b.BoundingBox(0, 0, &boxB)
    return boxA.Minimum.X < boxB.Minimum.X
}

func boxYCompare(a, b Hittable) bool {
    var boxA, boxB Aabb
    a.BoundingBox(0, 0, &boxA)
    b.BoundingBox(0, 0, &boxB)
    return boxA.Minimum.Y < boxB.Minimum.Y
}

func boxZCompare(a, b Hittable) bool {
    var boxA, boxB Aabb
    a.BoundingBox(0, 0, &boxA)
    b.BoundingBox(0, 0, &boxB)
    return boxA.Minimum.Z < boxB.Minimum.Z
}

func getComparisonFunctionForAxis(axis int) func(a, b Hittable)(bool) {
    if axis == 0 {
        return boxXCompare
    } else if axis == 1 {
        return boxYCompare
    }
    return boxZCompare
}

type sortByAxis struct {
    objects []Hittable
    axis int
}

func (s *sortByAxis) Len() int {
    return len(s.objects)
}

func (s *sortByAxis) Swap(i, j int) {
    s.objects[i], s.objects[j] = s.objects[j], s.objects[i]
}

func (s *sortByAxis) Less(i, j int) bool {
    // x-axis
    if s.axis == 0 {
        return boxXCompare(s.objects[i], s.objects[j])
    } else if s.axis == 1 {
        return boxYCompare(s.objects[i], s.objects[j])
    } 
    return boxZCompare(s.objects[i], s.objects[j])
}

func MakeBvh(objects []Hittable, time0 float64, time1 float64) *BvhNode {
    axis := RandInt(0, 2)

    var left, right Hittable
    if len(objects) == 1 {
        left = objects[0]
        right = left
    } else if len(objects) == 2 {
        if getComparisonFunctionForAxis(axis)(objects[0], objects[1]) {
            left = objects[0]
            right = objects[1]
        } else {
            left = objects[1]
            right = objects[0]
        }
    } else {
        sort.Sort(&sortByAxis{
            objects: objects,
            axis: axis,
        })

        mid := len(objects) / 2
        left = MakeBvh(objects[0:mid], time0, time1)
        right = MakeBvh(objects[mid:], time0, time1)
    }

    var leftBox, rightBox Aabb
    left.BoundingBox(time0, time1, &leftBox)
    right.BoundingBox(time0, time1, &rightBox)

    return &BvhNode{
        left: left,
        right: right,
        box: *surroundingBox(&leftBox, &rightBox),
    }
}
