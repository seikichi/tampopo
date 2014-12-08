package main

import (
	"flag"
	"fmt"

	"image"
	"image/color"

	"github.com/disintegration/imaging"
	"github.com/seikichi/tampopo/mser"
)

var maxArea = flag.Float64("maxArea", 1.0,
	"Maximum area of any stable region relative to the image domain area.")

var minArea = flag.Float64("minArea", 0.0,
	"Minimum area of any stable region relative to the image domain area.")

func main() {
	flag.Parse()

	if flag.NArg() != 2 {
		return
	}

	input, output := flag.Arg(0), flag.Arg(1)

	im, _ := imaging.Open(input)

	bounds := im.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	gray := image.NewGray(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			oldColor := im.At(bounds.Min.X+x, bounds.Min.Y+y)
			grayColor := color.GrayModel.Convert(oldColor)
			gray.Set(x, y, grayColor)
		}
	}

	forest := mser.ExtractMSERForest(gray, mser.Params{
		MaxArea: *maxArea,
		MinArea: *minArea,
	})

	rgba := im.(*image.NRGBA)

	count := 0
	for _, tree := range forest {
		count += countTree(tree)
		drawRegions(rgba, tree)
	}
	fmt.Println(count)

	imaging.Save(rgba, output)
}

func countTree(r *mser.ExtremalRegion) int {
	count := 1
	for _, child := range r.Children() {
		count += countTree(child)
	}
	return count
}

func drawRegions(rgba *image.NRGBA, r *mser.ExtremalRegion) {
	drawRect(rgba, r.Bounds())
	for _, child := range r.Children() {
		drawRegions(rgba, child)
	}
}

func drawRect(rgba *image.NRGBA, rect image.Rectangle) {
	for x := rect.Min.X; x < rect.Max.X; x++ {
		rgba.Set(x, rect.Min.Y, color.RGBA{255, 0, 0, 255})
		rgba.Set(x, rect.Max.Y, color.RGBA{255, 0, 0, 255})
	}
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		rgba.Set(rect.Min.X, y, color.RGBA{255, 0, 0, 255})
		rgba.Set(rect.Max.X, y, color.RGBA{255, 0, 0, 255})
	}
}
