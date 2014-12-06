package mser

import (
	"image"
	"testing"
)

type regionTree struct {
	level, area int
	point       image.Point
	children    []*regionTree
}

var testData = []struct {
	input  image.Gray
	output *regionTree
}{{
	input:  image.Gray{Pix: []uint8{}},
	output: nil,
}, {
	input: image.Gray{
		Pix:    []uint8{0},
		Stride: 1,
		Rect:   image.Rectangle{Min: image.Point{0, 0}, Max: image.Point{1, 1}}},
	output: &regionTree{level: 0, area: 1, point: image.Point{0, 0}},
}, {
	input: image.Gray{
		Pix: []uint8{
			1, 2, 2,
			2, 1, 1,
		},
		Stride: 3,
		Rect:   image.Rectangle{Min: image.Point{0, 0}, Max: image.Point{3, 2}}},
	output: &regionTree{level: 2, area: 6, point: image.Point{0, 1},
		children: []*regionTree{
			&regionTree{level: 1, area: 1, point: image.Point{0, 0}},
			&regionTree{level: 1, area: 2, point: image.Point{1, 1}}}},
}, {
	input: image.Gray{
		Pix: []uint8{
			3, 3,
			2, 2,
			1, 1,
		},
		Stride: 2,
		Rect:   image.Rectangle{Min: image.Point{0, 0}, Max: image.Point{2, 3}}},
	output: &regionTree{level: 3, area: 6, point: image.Point{0, 0},
		children: []*regionTree{
			&regionTree{level: 2, area: 4,
				point: image.Point{0, 1},
				children: []*regionTree{
					&regionTree{level: 1, area: 2,
						point: image.Point{0, 2}}}}}},
}, {
	input: image.Gray{
		Pix: []uint8{
			3, 1, 3,
			2, 3, 2,
			3, 1, 3,
		},
		Stride: 3,
		Rect:   image.Rectangle{Min: image.Point{0, 0}, Max: image.Point{3, 3}}},
	output: &regionTree{level: 3, area: 9, point: image.Point{0, 0},
		children: []*regionTree{
			&regionTree{level: 1, area: 1, point: image.Point{1, 0}},
			&regionTree{level: 1, area: 1, point: image.Point{1, 2}},
			&regionTree{level: 2, area: 1, point: image.Point{0, 1}},
			&regionTree{level: 2, area: 1, point: image.Point{2, 1}}}},
}, {
	input: image.Gray{
		Pix: []uint8{
			5, 5, 5, 9,
			4, 1, 2, 1,
			4, 3, 4, 2,
			3, 3, 3, 1,
		},
		Stride: 4,
		Rect:   image.Rectangle{Min: image.Point{0, 0}, Max: image.Point{4, 4}}},
	output: &regionTree{
		level: 9, area: 16, point: image.Point{3, 0},
		children: []*regionTree{
			&regionTree{level: 5, area: 15, point: image.Point{0, 0},
				children: []*regionTree{
					&regionTree{level: 4, area: 12, point: image.Point{0, 1},
						children: []*regionTree{
							&regionTree{level: 3, area: 9, point: image.Point{1, 1},
								children: []*regionTree{
									&regionTree{level: 2, area: 5, point: image.Point{2, 1},
										children: []*regionTree{
											&regionTree{level: 1, area: 1, point: image.Point{1, 1}},
											&regionTree{level: 1, area: 1, point: image.Point{3, 1}},
											&regionTree{level: 1, area: 1, point: image.Point{3, 3}}}}}}}}}}}},
}, {
	input: image.Gray{
		Pix: []uint8{
			0, 0, 0, 0, 0,
			0, 3, 1, 3, 0,
			0, 2, 3, 2, 0,
			0, 3, 1, 3, 0,
			0, 0, 0, 0, 0,
		},
		Stride: 5,
		Rect:   image.Rectangle{Min: image.Point{1, 1}, Max: image.Point{4, 4}}},
	output: &regionTree{level: 3, area: 9, point: image.Point{1, 1},
		children: []*regionTree{
			&regionTree{level: 1, area: 1, point: image.Point{2, 1}},
			&regionTree{level: 1, area: 1, point: image.Point{2, 3}},
			&regionTree{level: 2, area: 1, point: image.Point{1, 2}},
			&regionTree{level: 2, area: 1, point: image.Point{3, 2}}}},
}}

func TestBuildERTree(t *testing.T) {
	for _, td := range testData {
		assertERTree(t, td.input, BuildERTree(td.input), td.output)
	}
}

func contains(r *ExtremalRegion, p image.Point, im image.Gray) bool {
	visited := map[image.Point]struct{}{}
	var top image.Point
	que := []image.Point{r.Point()}
	for len(que) != 0 {
		que, top = que[1:], que[0]
		if _, ok := visited[top]; ok {
			continue
		}
		if int(im.Pix[im.PixOffset(top.X, top.Y)]) > r.Level() {
			continue
		}
		visited[top] = struct{}{}

		if top == p {
			return true
		}

		ds := []struct{ dx, dy int }{{+1, +0}, {-1, +0}, {+0, +1}, {+0, -1}}
		for _, d := range ds {
			np := image.Point{top.X + d.dx, top.Y + d.dy}
			if np.In(im.Bounds()) {
				que = append(que, np)
			}
		}
	}
	return false
}

func assertERTree(t *testing.T, im image.Gray, actual *ExtremalRegion, expected *regionTree) {
	if expected == nil && actual == nil {
		return
	}
	if actual.area != expected.area {
		t.Errorf("region.area = %v want %v where region = %+v, input = %+v",
			actual.area, expected.area, actual, im)
		return
	}
	if actual.level != expected.level {
		t.Errorf("region.level = %v want %v where region = %+v, input = %+v",
			actual.level, expected.level, actual, im)
		return
	}
	if len(actual.Children()) != len(expected.children) {
		t.Errorf("len(region.Children()) = %v want %v where region = %+v, input = %+v",
			len(actual.Children()), len(expected.children), actual, im)
		return
	}
	if !contains(actual, expected.point, im) {
		t.Errorf("region contains(%v, input) = false want true where region = %+v, input = %+v",
			expected.point, actual, im)
		return
	}

	for _, expChild := range expected.children {
		found := false
		for _, actChild := range actual.Children() {
			if contains(actChild, expChild.point, im) {
				assertERTree(t, im, actChild, expChild)
				found = true
				break
			}
		}
		if !found {
			t.Errorf("region corresponds to <%#v> not found in children of %#v",
				expChild, actual)
			return
		}
	}
}
