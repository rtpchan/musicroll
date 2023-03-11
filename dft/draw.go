package dft

// Draw an image of a keyboard, and representation of Fourier Transform

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"strings"

	// _ "image/png"
	"log"
	"math"
	"os"
)

var keyPosLeft map[string]float64
var keyPosMid map[string]float64
var imageWidth int

func initDraw() {
	keyPosLeft = make(map[string]float64)
	keyPosMid = make(map[string]float64)
	for i, k := range noteName {
		keyPosLeft[k] = keyPositionLeft[i]
		if strings.Contains(k, "s") {
			// black key
			keyPosMid[k] = keyPositionLeft[i] + 0.5/123.
		} else {
			// white key
			keyPosMid[k] = keyPositionLeft[i] + ((123./52.)/123.)/2.0
		}
	}
	imageWidth = 800 // pixel
}

func DrawEmptyPiano() image.Image {
	var img *image.RGBA
	f, err := os.Open("keyboard.png")
	if err != nil {
		img = image.NewRGBA(image.Rect(0, 0, 800, 600))
	}
	defer f.Close()
	return img

}

func DrawEmpty(height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 800, height))
	return img
}

// draw piano with data in value of key
func DrawPiano(data map[string]float64) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 800, 600))

	f, err := os.Open("keyboard.png")
	if err != nil {
		log.Println("Error loading image")
	}
	defer f.Close()

	imgF, err := png.Decode(f)
	if err != nil {
		log.Printf("Error decoding image: %s", err)
	}
	draw.Draw(img, img.Bounds(), imgF, imgF.Bounds().Min, draw.Src)

	maxValue := 0.0
	for k := range data {
		maxValue = math.Max(float64(data[k]), float64(maxValue))
	}
	// log.Printf("max coefficient %0.2f", maxValue)
	keyboardPixel := 780.
	for _, k := range noteName {
		x := keyPosMid[k]*keyboardPixel + 10 // 800 pixel
		value := uint8(data[k] / maxValue * 255)
		// draw small box representing the
		temp := image.NewRGBA(image.Rect(0, 0, 10, 50))
		colour := color.RGBA{255, 255 - value, 255 - value, 255}
		draw.Draw(temp, temp.Bounds(), &image.Uniform{colour}, image.Point{}, draw.Src)
		// draw
		b := image.Rect(int(x)-5, 450, int(x)+5, 490)
		draw.Draw(img, b, temp, image.Point{}, draw.Src)
	}

	return img
}

func DrawOnImage(dst image.Image, src image.Image, pt image.Point) image.Image {
	dstB := dst.Bounds()
	dstDraw := image.NewRGBA(image.Rect(0, 0, dstB.Dx(), dstB.Dy()))
	draw.Draw(dstDraw, dstDraw.Bounds(), dst, dstB.Min, draw.Src)
	r := image.Rectangle{pt, pt.Add(src.Bounds().Size())}
	draw.Draw(dstDraw, r, src, src.Bounds().Min, draw.Src)
	return dstDraw
}

func DrawNewStripe(value map[string]float64, maxValue float64, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, imageWidth, height))
	localMax := 0.0
	for k := range value {
		localMax = math.Max(float64(value[k]), localMax)
	}
	localMax = math.Max(localMax, maxValue)
	for _, k := range noteName {
		x := keyPosMid[k] * float64(imageWidth) // 800 pixel
		// v := uint8(value[k] / localMax * 255)  // linear value
		v := uint8(math.Pow((value[k]/localMax), 3) * 255.)
		temp := image.NewRGBA(image.Rect(0, 0, 10, 10))
		var colour color.RGBA
		if strings.Contains(k, "s") {
			colour = color.RGBA{0, 0, v, v}
		} else {
			colour = color.RGBA{v, 0, 0, v}
		}
		draw.Draw(temp, temp.Bounds(), &image.Uniform{colour}, image.Point{}, draw.Over)
		// draw
		b := image.Rect(int(x)-5, 0, int(x)+5, height)
		draw.Draw(img, b, temp, image.Point{}, draw.Over)
	}

	return img
}

func DrawSpectrum(s []map[string]float64, height int) image.Image {

	imageHeight := len(s) * height
	bg := DrawEmpty(imageHeight) // 10 pixel * num of strips

	for i, spec := range s {
		// log.Printf("inside drawingspectrum, %d, %f", i, spec["C4"])
		img := DrawNewStripe(spec, 0.001, 10)
		bg = DrawOnImage(bg, img, image.Point{0, imageHeight - ((i + 1) * 10)})
		// err = Export(fmt.Sprintf("test%02d.png", i), img)
		// if err != nil {
		// 	t.Errorf("export png : %v", err)
		// }
	}
	return bg
}

func Export(name string, m image.Image) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}

	err = png.Encode(f, m)
	if err != nil {
		return err
	}
	return nil
}

var keyPositionLeft = []float64{
	0,
	0.0170731707317073,
	0.0192307692307692,
	0.0384615384615385,
	0.0518605378361476,
	0.0576923076923077,
	0.074624765478424,
	0.0769230769230769,
	0.0961538461538462,
	0.109412132582864,
	0.115384615384615,
	0.130550343964978,
	0.134615384615385,
	0.151688555347092,
	0.153846153846154,
	0.173076923076923,
	0.186475922451532,
	0.192307692307692,
	0.209240150093809,
	0.211538461538462,
	0.230769230769231,
	0.244027517198249,
	0.25,
	0.265165728580363,
	0.269230769230769,
	0.286303939962477,
	0.288461538461538,
	0.307692307692308,
	0.321091307066917,
	0.326923076923077,
	0.343855534709193,
	0.346153846153846,
	0.365384615384615,
	0.378642901813634,
	0.384615384615385,
	0.399781113195747,
	0.403846153846154,
	0.420919324577861,
	0.423076923076923,
	0.442307692307692,
	0.455706691682302,
	0.461538461538462,
	0.478470919324578,
	0.480769230769231,
	0.5,
	0.513258286429018,
	0.519230769230769,
	0.534396497811132,
	0.538461538461538,
	0.555534709193246,
	0.557692307692308,
	0.576923076923077,
	0.590322076297686,
	0.596153846153846,
	0.613086303939962,
	0.615384615384615,
	0.634615384615385,
	0.647873671044403,
	0.653846153846154,
	0.669011882426517,
	0.673076923076923,
	0.69015009380863,
	0.692307692307692,
	0.711538461538462,
	0.724937460913071,
	0.730769230769231,
	0.747701688555347,
	0.75,
	0.769230769230769,
	0.782489055659787,
	0.788461538461538,
	0.803627267041901,
	0.807692307692308,
	0.824765478424015,
	0.826923076923077,
	0.846153846153846,
	0.859552845528455,
	0.865384615384615,
	0.882317073170732,
	0.884615384615385,
	0.903846153846154,
	0.917104440275172,
	0.923076923076923,
	0.938242651657286,
	0.942307692307692,
	0.9593808630394,
	0.961538461538462,
	0.980769230769231,
}
