package mser

import "image"

// An ExtremalRegion represents a maximum intensity region.
type ExtremalRegion struct {
	level, pixel, area  int
	x, y                int
	parent, next, child *ExtremalRegion
}

// Level returns pixel level of ExtremalRegion.
func (r *ExtremalRegion) Level() int { return r.level }

// Point returns a point belongs to the ExtremalRegion.
func (r *ExtremalRegion) Point() image.Point { return image.Point{r.x, r.y} }

// Area returns area of ExtremalRegion.
func (r *ExtremalRegion) Area() int { return r.area }

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
	r.area++
}

func (r *ExtremalRegion) merge(child *ExtremalRegion) {
	r.area += child.area
	child.parent = r
	child.next = r.child
	r.child = child
}

// BuildERTree returns ERs tree from given image.
func BuildERTree(im image.Gray) *ExtremalRegion {
	bounds := im.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	if width == 0 && height == 0 {
		return nil
	}

	priority := 256
	heap := make([][]int, 256)
	stack := []*ExtremalRegion{}
	accessible := make([]bool, width*height)
	stack = append(stack, &ExtremalRegion{level: 256})

	curPixel, curEdge, curLevel := 0, 0, int(im.Pix[0])
	accessible[0] = true

step3:
	for {
		stack = append(stack, &ExtremalRegion{
			level: curLevel,
			pixel: curPixel,
			x:     curPixel % width,
			y:     curPixel / width})
		for {
			x, y := curPixel%width, curPixel/width
			for ; curEdge < 4; curEdge++ {
				neighbourPixel := curPixel
				if curEdge == 0 && x < width-1 {
					neighbourPixel = curPixel + 1
				} else if curEdge == 1 && y < height-1 {
					neighbourPixel = curPixel + width
				} else if curEdge == 2 && x > 0 {
					neighbourPixel = curPixel - 1
				} else if curEdge == 3 && y > 0 {
					neighbourPixel = curPixel - width
				}
				if neighbourPixel != curPixel && !accessible[neighbourPixel] {
					neighbourLevel := int(im.Pix[neighbourPixel])
					accessible[neighbourPixel] = true
					if neighbourLevel >= curLevel {
						heap[neighbourLevel] = append(heap[neighbourLevel], neighbourPixel<<4)
						if neighbourLevel < priority {
							priority = neighbourLevel
						}
					} else {
						heap[curLevel] = append(heap[curLevel], (curPixel<<4)|(curEdge+1))
						if curLevel < priority {
							priority = curLevel
						}
						curPixel, curEdge, curLevel = neighbourPixel, 0, neighbourLevel
						continue step3
					}
				}
			}
			stack[len(stack)-1].accumulate(x, y)
			if priority == 256 {
				return stack[len(stack)-1]
			}
			last := heap[priority][len(heap[priority])-1]
			heap[priority] = heap[priority][:len(heap[priority])-1]
			curPixel, curEdge = last>>4, last&15

			for priority < 256 && len(heap[priority]) == 0 {
				priority++
			}

			if int(im.Pix[curPixel]) != curLevel {
				curLevel = int(im.Pix[curPixel])
				stack = processStack(curLevel, curPixel, stack, width)
			}
		}
	}
}

func processStack(newPixelGreyLevel, pixel int, stack []*ExtremalRegion, width int) []*ExtremalRegion {
	for {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if newPixelGreyLevel < stack[len(stack)-1].level {
			stack = append(stack, &ExtremalRegion{
				level: newPixelGreyLevel,
				pixel: pixel,
				x:     pixel % width,
				y:     pixel / width})
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
