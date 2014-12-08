package main

import (
	"fmt"
	"image"

	"github.com/seikichi/tampopo/mser"
)

func main() {
	// im := &image.Gray{
	// 	Pix:    []uint8{1},
	// 	Stride: 4,
	// 	Rect:   image.Rect(0, 0, 1, 1)}
	im := &image.Gray{
		Pix: []uint8{
			1, 1, 9, 1,
			1, 2, 9, 9,
			9, 9, 3, 1,
			9, 9, 1, 1,
		},
		Stride: 4,
		Rect:   image.Rect(0, 0, 4, 4)}
	tree := mser.ExtractERTree(im)
	printERTree(tree, 0)

	fmt.Println("!!!!!!!!")

	forest := mser.ExtractMSERForest(im, mser.Params{
		Delta:        2,
		MinArea:      0.2,
		MaxArea:      1.1,
		MaxVariation: 0.9,
		MinDiversity: 0.0,
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
