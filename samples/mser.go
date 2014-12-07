package main

import (
	"fmt"
	"image"

	"github.com/seikichi/tampopo/mser"
)

func main() {
	im := &image.Gray{
		Pix: []uint8{
			1, 1, 9, 1,
			1, 2, 9, 9,
			9, 9, 3, 1,
			9, 9, 1, 1,
		},
		Stride: 4,
		Rect: image.Rectangle{
			Min: image.Point{0, 0},
			Max: image.Point{4, 4}}}
	tree := mser.ExtractERTree(im)
	printERTree(tree, 0)

	fmt.Println("!!!!!!!!")

	forest := mser.ExtractMSERForest(im, mser.Params{
		Delta:        2,
		MinArea:      0.2,
		MaxArea:      0.5,
		MaxVariation: 0.5,
		MinDiversity: 0.33,
	})
	for _, tree := range forest {
		printERTree(tree, 0)
	}
}

func printERTree(r *mser.ExtremalRegion, index int) {
	for i := 0; i < index; i++ {
		fmt.Print(" ")
	}
	fmt.Printf("Region{level: %v, area: %v, bounds: %v}\n", r.Level(), r.Area(), r.Bounds())
	for _, child := range r.Children() {
		printERTree(child, index+2)
	}
}
