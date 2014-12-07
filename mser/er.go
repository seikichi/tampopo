package mser

import (
	"image"
	"image/color"
)

// An ExtremalRegion represents a maximum intensity region.
type ExtremalRegion struct {
	// seed point and threshold
	point image.Point
	level int

	// incrementally computable features
	area           int
	rect           image.Rectangle
	rawMoments     [2]int
	centralMoments [3]int

	// fields for stability computation
	variation float64
	stable    bool

	// pointers preserving the tree structure of the component tree
	parent, next, child *ExtremalRegion
}

// Level returns pixel level of ExtremalRegion.
func (r *ExtremalRegion) Level() int { return r.level }

// Area returns area of ExtremalRegion.
func (r *ExtremalRegion) Area() int { return r.area }

// Point returns a point belongs to the ExtremalRegion.
func (r *ExtremalRegion) Point() image.Point { return r.point }

// Variation returns variaion of ExtremalRegion.
func (r *ExtremalRegion) Variation() float64 { return r.variation }

// Bounds returns bounding box of the region.
func (r *ExtremalRegion) Bounds() image.Rectangle { return r.rect }

// Parent returns parent region.
func (r *ExtremalRegion) Parent() *ExtremalRegion { return r.parent }

// Children returns children regions.
func (r *ExtremalRegion) Children() []*ExtremalRegion {
	children := []*ExtremalRegion{}
	for child := r.child; child != nil; child = child.next {
		children = append(children, child)
	}
	return children
}

func (r *ExtremalRegion) accumulate(x, y int) {
	if r.area == 0 {
		r.rect = image.Rect(x, y, x+1, y+1)
	} else {
		r.rect = r.rect.Union(image.Rect(x, y, x+1, y+1))
	}
	r.area++

	r.rawMoments[0] += x
	r.rawMoments[1] += x
	r.centralMoments[0] += x * x
	r.centralMoments[1] += x * y
	r.centralMoments[2] += y * y
}

func (r *ExtremalRegion) merge(child *ExtremalRegion) {
	if r.area == 0 {
		r.rect = child.rect
	} else {
		r.rect = r.rect.Union(child.rect)
	}
	r.area += child.area

	r.rawMoments[0] += child.rawMoments[0]
	r.rawMoments[1] += child.rawMoments[1]
	r.centralMoments[0] += child.centralMoments[0]
	r.centralMoments[1] += child.centralMoments[1]
	r.centralMoments[2] += child.centralMoments[2]

	child.parent, child.next, r.child = r, r.child, child
}

type searchState struct {
	point image.Point
	edge  int
}

// ExtractERTree extracts the ER component tree of the image.
func ExtractERTree(im *image.Gray) *ExtremalRegion {
	bounds := im.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	if width == 0 && height == 0 {
		return nil
	}

	priority := 256
	heap := make([][]searchState, 256)
	stack := []*ExtremalRegion{}
	accessible := make([]bool, width*height)
	stack = append(stack, &ExtremalRegion{level: 256})

	curPixel, curEdge := bounds.Min, 0
	curLevel := int(im.At(bounds.Min.X, bounds.Min.Y).(color.Gray).Y)
	accessible[0] = true

step3:
	for {
		stack = append(stack, &ExtremalRegion{level: curLevel, point: curPixel})
		for {
			for ; curEdge < 4; curEdge++ {
				var neighbourPixel image.Point
				switch curEdge {
				case 0:
					neighbourPixel = image.Point{curPixel.X + 1, curPixel.Y}
				case 1:
					neighbourPixel = image.Point{curPixel.X, curPixel.Y + 1}
				case 2:
					neighbourPixel = image.Point{curPixel.X - 1, curPixel.Y}
				default:
					neighbourPixel = image.Point{curPixel.X, curPixel.Y - 1}
				}

				index := (neighbourPixel.Y-bounds.Min.Y)*width + (neighbourPixel.X - bounds.Min.X)
				if neighbourPixel.In(bounds) && !accessible[index] {
					neighbourLevel := int(im.At(neighbourPixel.X, neighbourPixel.Y).(color.Gray).Y)
					accessible[index] = true
					if neighbourLevel >= curLevel {
						heap[neighbourLevel] = append(heap[neighbourLevel], searchState{point: neighbourPixel})
						if neighbourLevel < priority {
							priority = neighbourLevel
						}
					} else {
						heap[curLevel] = append(heap[curLevel], searchState{point: curPixel, edge: curEdge + 1})
						if curLevel < priority {
							priority = curLevel
						}
						curPixel, curEdge, curLevel = neighbourPixel, 0, neighbourLevel
						continue step3
					}
				}
			}
			stack[len(stack)-1].accumulate(curPixel.X, curPixel.Y)
			if priority == 256 {
				return stack[len(stack)-1]
			}
			last := heap[priority][len(heap[priority])-1]
			heap[priority] = heap[priority][:len(heap[priority])-1]
			curPixel, curEdge = last.point, last.edge

			for priority < 256 && len(heap[priority]) == 0 {
				priority++
			}

			newPixelGreyLevel := int(im.At(curPixel.X, curPixel.Y).(color.Gray).Y)
			if newPixelGreyLevel != curLevel {
				curLevel = newPixelGreyLevel
				stack = processStack(newPixelGreyLevel, curPixel, stack)
			}
		}
	}
}

func processStack(newPixelGreyLevel int, pixel image.Point, stack []*ExtremalRegion) []*ExtremalRegion {
	for {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if newPixelGreyLevel < stack[len(stack)-1].level {
			stack = append(stack, &ExtremalRegion{level: newPixelGreyLevel, point: pixel})
			stack[len(stack)-1].merge(top)
			return stack
		}
		stack[len(stack)-1].merge(top)
		if newPixelGreyLevel <= stack[len(stack)-1].level {
			break
		}
	}
	return stack
}
