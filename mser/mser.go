// Package mser provides Maximum Stable Extremal Region algorithms.
package mser

import (
	"fmt"
	"image"
)

// Params represents MSER algorithm paraemters.
type Params struct {
	Delta                                        int
	MinArea, MaxArea, MaxVariation, MinDiversity float64
}

// ExtractMSERForest extracts the MSER component forest of the image.
func ExtractMSERForest(im *image.Gray, params Params) []*ExtremalRegion {
	tree := ExtractERTree(im)

	bounds := im.Bounds()
	size := bounds.Dx() * bounds.Dy()
	minArea := int(params.MinArea * float64(size))
	maxArea := int(params.MaxArea * float64(size))
	tree.process(params.Delta, minArea, maxArea, params.MaxVariation)
	return tree.extractMSER(params.MinDiversity)
}

func (r *ExtremalRegion) process(delta, minArea, maxArea int, maxVariation float64) {
	parent := r
	for parent.parent != nil && parent.parent.level <= r.level+delta {
		parent = parent.parent
	}

	r.variation = float64(parent.area-r.area) / float64(r.area)
	stable := (r.parent == nil) || (r.variation <= r.parent.variation)
	stable = stable && r.area >= minArea && r.area <= maxArea && r.variation <= maxVariation

	fmt.Println(r.area, r.variation)

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

func (r *ExtremalRegion) extractMSER(minDiversity float64) []*ExtremalRegion {
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
	children := []*ExtremalRegion{}
	for child := r.child; child != nil; child = child.next {
		children = append(children, child.extractMSER(minDiversity)...)
	}
	for i, child := range children {
		child.parent = nil
		if i+1 < len(children) {
			child.next = children[i+1]
		}
	}

	if r.stable {
		root := *r
		for _, child := range children {
			child.parent = &root
		}
		root.parent, root.child, root.next = nil, nil, nil
		if len(children) > 0 {
			root.child = children[0]
		}
		return []*ExtremalRegion{&root}
	}
	return children
}
