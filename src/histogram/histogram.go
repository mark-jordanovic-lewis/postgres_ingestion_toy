package histogram

import (
	"fmt"
	"image/color"

	plot "gonum.org/v1/plot"
	plotter "gonum.org/v1/plot/plotter"
	vg "gonum.org/v1/plot/vg"
)

func BuildHistogram(filename string, data []float64) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Could not generate Histogram: %v\n\n\t%v", filename, r)
		}
	}()
	values := make(plotter.Values, len(data))
	for i := range values {
		values[i] = data[i]
	}

	canvas, err := plot.New()
	if err != nil {
		panic(err)
	}
	canvas.Title.Text = fmt.Sprintf("Binned times for %v", filename)
	canvas.X.Label.Text = "ingested rows per s"
	canvas.Y.Label.Text = "frequency"

	hist, err := plotter.NewHist(values, 20)
	if err != nil {
		panic(err)
	}

	// hist.Normalize(1)
	canvas.Add(hist)

	_, _, ymax, ymin := hist.DataRange()
	mu := mean(data)
	xyer := plotter.XYs{{X: mu, Y: ymin}, {X: mu, Y: ymax}}
	meanline, err := plotter.NewLine(xyer)
	if err != nil {
		panic(err)
	}
	meanline.Color = color.RGBA{G: 169, A: 255}
	canvas.Legend.Top = true
	canvas.Legend.Left = false
	canvas.Legend.Add(fmt.Sprintf("mean rps: %v", int(mu)), meanline)
	canvas.Add(meanline)
	if err := canvas.Save(4*vg.Inch, 4*vg.Inch, filename); err != nil {
		panic(err)
	}
}

func mean(ts []float64) (mu float64) {
	for _, t := range ts {
		mu += t
	}
	mu /= float64(len(ts))
	return
}
