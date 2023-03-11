package dft

import (
	"bytes"
	"image"
	"image/png"
	"log"
	"strings"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

// Plot y against x
func PlotXY(x, y []float64, title, labelX, labelY, filename string) image.Image {
	p := plot.New()
	p.Title.Text = title
	p.X.Label.Text = labelX
	p.Y.Label.Text = labelY

	data := make(plotter.XYs, len(y))
	for i, v := range y {
		data[i].Y = v
		data[i].X = x[i]
	}

	plotline, err := plotter.NewLine(data)
	if err != nil {
		log.Println(err)
	}
	p.Add(plotline)

	// Save Image
	if err := p.Save(4*vg.Inch, 4*vg.Inch, filename); err != nil {
		log.Println(err)
	}

	img := saveFile(p, filename)
	return img
}

// Plot y, x goes from 0 to len(y)
func PlotY(y []float64, title, labelX, labelY, filename string) image.Image {
	p := plot.New()
	p.Title.Text = title
	p.X.Label.Text = labelX
	p.Y.Label.Text = labelY

	data := make(plotter.XYs, len(y))
	for i, v := range y {
		data[i].Y = v
		data[i].X = float64(i)
	}

	plotline, err := plotter.NewLine(data)
	if err != nil {
		log.Println(err)
	}
	p.Add(plotline)

	// Save Image
	if err := p.Save(4*vg.Inch, 4*vg.Inch, filename); err != nil {
		log.Println(err)
	}

	img := saveFile(p, filename)
	return img
}

// Plot value of map against the key
func PlotMap(y map[string]float64, filename string) image.Image {

	p := plot.New()
	p.Title.Text = "Spectrum"
	p.X.Label.Text = "Keys"
	p.Y.Label.Text = "Values"

	data := plotter.Values{}
	for _, key := range noteName {
		data = append(data, y[key])
	}

	bars, err := plotter.NewBarChart(data, vg.Points(2))
	if err != nil {
		log.Println(err)
	}
	p.Add(bars)
	p.NominalX(filterdNoteName()...)

	img := saveFile(p, filename)
	return img
}

func saveFile(p *plot.Plot, filename string) image.Image {
	if filename == "" {
		// output image.Image if no filename is empty string
		wt, err := p.WriterTo(4*vg.Inch, 4*vg.Inch, "png")
		if err != nil {
			log.Printf("saveFile error %s\n", err)
		}
		var buf bytes.Buffer
		_, err = wt.WriteTo(&buf)
		if err != nil {
			log.Printf("saveFile error %s\n", err)
		}

		img, err := png.Decode(&buf)
		if err != nil {
			log.Printf("saveFile error %s\n", err)
		}
		return img

	} else {
		// Save Image
		if err := p.Save(8*vg.Inch, 6*vg.Inch, filename); err != nil {
			log.Println(err)
		}
		return nil
	}
}

func filterdNoteName() []string {
	names := []string{}
	for _, n := range noteName {
		if strings.Contains(n, "C") && !strings.Contains(n, "s") {
			names = append(names, n)
		} else {
			names = append(names, "")
		}
	}
	return names
}
