package mser

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"io"
	"sort"
	"testing"
)

func newRect(x, y, w, h int) image.Rectangle { return image.Rect(x, y, w, h) }

func newGray(width, height int, pix []uint8) *image.Gray {
	return &image.Gray{Pix: pix, Stride: width, Rect: image.Rect(0, 0, width, height)}
}

type fromFields struct {
	level int
	area  int
	rect  image.Rectangle
}

func newER(p fromFields, children ...*ExtremalRegion) *ExtremalRegion {
	r := &ExtremalRegion{level: p.level, area: p.area, rect: p.rect}
	if len(children) != 0 {
		r.child = children[0]
	}
	for i, child := range children {
		child.parent = r
		if i+1 < len(children) {
			child.next = children[i+1]
		} else {
			child.next = nil
		}
	}
	return r
}

var testData = []struct {
	input  *image.Gray
	output *ExtremalRegion
}{{
	input:  newGray(1, 1, []uint8{0}),
	output: newER(fromFields{level: 0, area: 1, rect: newRect(0, 0, 1, 1)}),
}, {
	input: newGray(3, 2, []uint8{
		1, 2, 2,
		2, 1, 1}),
	output: newER(fromFields{level: 2, area: 6, rect: newRect(0, 0, 3, 2)},
		newER(fromFields{level: 1, area: 1, rect: newRect(0, 0, 1, 1)}),
		newER(fromFields{level: 1, area: 2, rect: newRect(1, 1, 3, 2)})),
}, {
	input: newGray(2, 3, []uint8{
		3, 3,
		2, 2,
		1, 1}),
	output: newER(fromFields{level: 3, area: 6, rect: newRect(0, 0, 2, 3)},
		newER(fromFields{level: 2, area: 4, rect: newRect(0, 1, 2, 3)},
			newER(fromFields{level: 1, area: 2, rect: newRect(0, 2, 2, 3)}))),
}, {
	input: newGray(3, 3, []uint8{
		3, 1, 3,
		2, 3, 2,
		3, 1, 3}),
	output: newER(fromFields{level: 3, area: 9, rect: newRect(0, 0, 3, 3)},
		newER(fromFields{level: 1, area: 1, rect: newRect(1, 0, 2, 1)}),
		newER(fromFields{level: 2, area: 1, rect: newRect(0, 1, 1, 2)}),
		newER(fromFields{level: 2, area: 1, rect: newRect(2, 1, 3, 2)}),
		newER(fromFields{level: 1, area: 1, rect: newRect(1, 2, 2, 3)})),
}, {
	input: newGray(4, 4, []uint8{
		5, 5, 5, 9,
		4, 1, 2, 1,
		4, 3, 4, 2,
		3, 3, 3, 1}),
	output: newER(fromFields{level: 9, area: 16, rect: newRect(0, 0, 4, 4)},
		newER(fromFields{level: 5, area: 15, rect: newRect(0, 0, 4, 4)},
			newER(fromFields{level: 4, area: 12, rect: newRect(0, 1, 4, 4)},
				newER(fromFields{level: 3, area: 9, rect: newRect(0, 1, 4, 4)},
					newER(fromFields{level: 2, area: 5, rect: newRect(1, 1, 4, 4)},
						newER(fromFields{level: 1, area: 1, rect: newRect(1, 1, 2, 2)}),
						newER(fromFields{level: 1, area: 1, rect: newRect(3, 1, 4, 2)}),
						newER(fromFields{level: 1, area: 1, rect: newRect(3, 3, 4, 4)})))))),
}, {
	input: newGray(5, 5, []uint8{
		0, 0, 0, 0, 0,
		0, 3, 1, 3, 0,
		0, 2, 3, 2, 0,
		0, 3, 1, 3, 0,
		0, 0, 0, 0, 0}).SubImage(image.Rect(1, 1, 4, 4)).(*image.Gray),
	output: newER(fromFields{level: 3, area: 9, rect: newRect(1, 1, 4, 4)},
		newER(fromFields{level: 1, area: 1, rect: newRect(2, 1, 3, 2)}),
		newER(fromFields{level: 2, area: 1, rect: newRect(1, 2, 2, 3)}),
		newER(fromFields{level: 2, area: 1, rect: newRect(3, 2, 4, 3)}),
		newER(fromFields{level: 1, area: 1, rect: newRect(2, 3, 3, 4)})),
}}

func TestExtractERTree(t *testing.T) {
	for i, tc := range testData {
		tree := ExtractERTree(tc.input)
		if !assertERsEqual(t, tc.output, tree) {
			imageBuf, expBuf, actBuf := new(bytes.Buffer), new(bytes.Buffer), new(bytes.Buffer)
			printGray(tc.input, 2, imageBuf)
			printER(tc.output, 2, expBuf)
			printER(tree, 2, actBuf)

			t.Errorf("test case %d:\ninput (min = %v):\n%v\nexpected:\n%v\nactual:\n%v\n",
				i+1, tc.input.Bounds().Min, imageBuf.String(), expBuf.String(), actBuf.String())
		}
	}
}

func printGray(im *image.Gray, indent int, w io.Writer) {
	bounds := im.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for i := 0; i < indent; i++ {
			fmt.Fprintf(w, " ")
		}
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			fmt.Fprintf(w, " %03d", im.At(x, y).(color.Gray).Y)
		}
		fmt.Fprintf(w, "\n")
	}
}

func printER(r *ExtremalRegion, indent int, w io.Writer) {
	for i := 0; i < indent; i++ {
		fmt.Fprintf(w, " ")
	}
	fmt.Fprintf(w, "{level: %v, area: %v, rect: %v}\n", r.level, r.area, r.rect)
	for _, child := range r.Children() {
		printER(child, indent+2, w)
	}
}

type erSorter []*ExtremalRegion

// Len returns the size of erSorter.
func (s erSorter) Len() int { return len(s) }

// Swap changes the place of ERs.
func (s erSorter) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// Less compares two ERs.
func (s erSorter) Less(i, j int) bool {
	pi, pj := s[i].rect.Min, s[j].rect.Min
	if pi.X != pj.X {
		return pi.X < pj.X
	}
	return pi.Y < pj.Y
}

func assertERsEqual(t *testing.T, exp, act *ExtremalRegion) bool {
	equal := true
	if exp.area != act.area {
		equal = false
		t.Errorf("region.area = %d want %d", act.area, exp.area)
	}
	if exp.level != act.level {
		equal = false
		t.Errorf("region.level = %d want %d", act.level, exp.level)
	}
	if exp.rect != act.rect {
		equal = false
		t.Errorf("region.rect = %d want %d", act.rect, exp.rect)
	}

	return equal && assertChildrenEqual(t, exp, act)
}

func assertChildrenEqual(t *testing.T, exp, act *ExtremalRegion) bool {
	ec, ac := erSorter(exp.Children()), erSorter(act.Children())
	sort.Sort(ec)
	sort.Sort(ac)

	if len(ec) != len(ac) {
		t.Errorf("len(region.Children()) = %d want %d", len(ac), len(ec))
		return false
	}
	for i := range ac {
		if ac[i].parent != act {
			t.Errorf("region.parent = <%v> want <%v>", ac[i].parent, act)
			return false
		}
		if !assertERsEqual(t, ec[i], ac[i]) {
			return false
		}
	}
	return true
}
