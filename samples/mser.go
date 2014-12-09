package main

import (
	"flag"

	"image"
	"image/color"

	"github.com/disintegration/imaging"
	"github.com/seikichi/tampopo/mser"
)

var delta = flag.Int("delta", 2, "DELTA parameter of the MSER algorithm.")
var maxArea = flag.Float64("maxArea", 1.0, "Maximum area relative to the image area.")
var minArea = flag.Float64("minArea", 0.0, "Minimum area relative to the image area.")
var maxVariation = flag.Float64("maxVariation", 1.0, "Maximum variation of the regions.")
var minDiversity = flag.Float64("minDiversity", 0.0, "Minimum diversity of the regions.")

var twoPass = flag.Bool("twoPass", true, "Extract MSERs to the inversed image.")

func main() {
	flag.Parse()

	if flag.NArg() != 2 {
		flag.Usage()
		return
	}

	input, output := flag.Arg(0), flag.Arg(1)

	im, _ := imaging.Open(input)

	bounds := im.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	gray := image.NewGray(image.Rect(0, 0, w, h))
	inv := image.NewGray(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			oldColor := im.At(bounds.Min.X+x, bounds.Min.Y+y)
			grayColor := color.GrayModel.Convert(oldColor).(color.Gray)
			gray.Set(x, y, grayColor)

			grayColor.Y = 255 - grayColor.Y
			inv.Set(x, y, grayColor)
		}
	}

	rgba := im.(*image.NRGBA)

	params := mser.Params{
		Delta:        *delta,
		MinArea:      *minArea,
		MaxArea:      *maxArea,
		MaxVariation: *maxVariation,
		MinDiversity: *minDiversity,
	}

	forest := mser.ExtractMSERForest(gray, params)
	for _, tree := range forest {
		drawRegions(rgba, tree, &color.RGBA{255, 0, 0, 255})
	}

	if *twoPass {
		forest = mser.ExtractMSERForest(inv, params)
		for _, tree := range forest {
			drawRegions(rgba, tree, &color.RGBA{0, 0, 255, 255})
		}
	}

	imaging.Save(rgba, output)
}

func countTree(r *mser.ExtremalRegion) int {
	count := 1
	for _, child := range r.Children() {
		count += countTree(child)
	}
	return count
}

func drawRegions(rgba *image.NRGBA, r *mser.ExtremalRegion, c *color.RGBA) {
	drawRect(rgba, r.Bounds(), c)
	for _, child := range r.Children() {
		drawRegions(rgba, child, c)
	}
}

func drawRect(rgba *image.NRGBA, rect image.Rectangle, c *color.RGBA) {
	for x := rect.Min.X; x < rect.Max.X; x++ {
		rgba.Set(x, rect.Min.Y, c)
		rgba.Set(x, rect.Max.Y, c)
	}
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		rgba.Set(rect.Min.X, y, c)
		rgba.Set(rect.Max.X, y, c)
	}
}
