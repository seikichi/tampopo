// Package mser provides Maximum Stable Extremal Region algorithms.
package mser

import "image"

// A MSER represents Maximum Stable Extremal Region.
type MSER struct {
	level, area int
	point       image.Point
	variation   float64

	parent   *MSER
	children []*MSER
}

// Level returns pixel level of ExtremalRegion.
func (r *MSER) Level() int { return r.level }

// Area returns area of MSER.
func (r *MSER) Area() int { return r.area }

// Point returns a point belongs to the MSER.
func (r *MSER) Point() image.Point { return r.point }

// Parent returns parent region.
func (r *MSER) Parent() *MSER { return r.parent }

// Variation returns variaion of MSER.
func (r *MSER) Variation() float64 { return r.variation }

// Children returns MSER children.
func (r *MSER) Children() []*MSER { return r.children }

func (r *ExtremalRegion) process(delta, minArea, maxArea int, maxVariation float64) {
	parent := r
	for parent.parent != nil && parent.parent.level <= r.level+delta {
		parent = parent.parent
	}

	r.variation = float64(parent.area-r.area) / float64(r.area)
	stable := (parent == nil) || (r.variation <= parent.variation)
	stable = stable && r.area >= minArea && r.area <= maxArea && r.variation <= maxVariation

	for child := r.child; child != nil; child = child.next {
		child.process(delta, minArea, maxArea, maxVariation)
		r.stable = r.stable || (stable && r.variation < child.variation)
	}

	r.stable = r.stable || (r.child == nil && stable)
}

func (r *ExtremalRegion) check(variation float64, area int) bool {
	if r.area <= area {
		return true
	}
	if r.stable && r.variation < variation {
		return false
	}
	for child := r.child; child != nil; child = child.next {
		if !child.check(variation, area) {
			return false
		}
	}
	return true
}

func (r *ExtremalRegion) buildMSERForest(minDiversity float64) []*MSER {
	if r.stable {
		minParentArea := int(float64(r.area)/(1.0-minDiversity) + 0.5)
		parent := r

		for parent.parent != nil && parent.parent.area < minParentArea {
			parent = parent.parent
			if parent.stable && parent.variation <= r.variation {
				r.stable = false
				break
			}
		}
		if r.stable {
			maxChildArea := int(float64(r.area)*(1.0-minDiversity) + 0.5)
			if !r.check(r.variation, maxChildArea) {
				r.stable = false
			}
		}
	}
	children := []*MSER{}
	for child := r.child; child != nil; child = child.next {
		children = append(children, child.buildMSERForest(minDiversity)...)
	}
	if r.stable {
		region := &MSER{
			level:     r.level,
			area:      r.area,
			point:     r.point,
			variation: r.variation,
			children:  children,
		}
		for _, child := range children {
			child.parent = region
		}
		return []*MSER{region}
	}
	return children
}

// Params represents MSER algorithm paraemters.
type Params struct {
	Delta                                        int
	MinArea, MaxArea, MaxVariation, MinDiversity float64
}

// BuildMSERForest returns MSERs forest from given image.
func BuildMSERForest(im *image.Gray, params *Params) []*MSER {
	tree := BuildERTree(im)

	bounds := im.Bounds()
	size := bounds.Dx() * bounds.Dy()
	minArea := int(params.MinArea * float64(size))
	maxArea := int(params.MaxArea * float64(size))
	tree.process(params.Delta, minArea, maxArea, params.MaxVariation)
	return tree.buildMSERForest(params.MinDiversity)
}
