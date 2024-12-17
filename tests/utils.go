package tests

import (
	"fmt"
	"image/color"

	"github.com/combo23/mimic/models"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

type VisualizationOptions struct {
	ShowControlPoints bool
	ShowGrid          bool
	ShowArrows        bool
	Title             string
	OutputPath        string
	Width             float64 // in inches
	Height            float64 // in inches
}

// Visualize creates a plot of the mouse movement path
func visualize(movement models.Movement, controlPoints []models.Point, opts VisualizationOptions) error {
	// Create a new plot
	p := plot.New()

	// Set plot title and labels
	p.Title.Text = opts.Title
	p.X.Label.Text = "X Position (pixels)"
	p.Y.Label.Text = "Y Position (pixels)"

	// Create main path line
	pathXYs := make(plotter.XYs, len(movement.Points))
	for i, point := range movement.Points {
		pathXYs[i].X = point.X
		pathXYs[i].Y = point.Y
	}

	// Create line plot for the path
	pathLine, err := plotter.NewLine(pathXYs)
	if err != nil {
		return fmt.Errorf("error creating path line: %v", err)
	}
	pathLine.Color = color.RGBA{B: 255, A: 255}
	pathLine.Width = vg.Points(1)

	// Add the path line to the plot
	p.Add(pathLine)

	// Add control points if enabled
	if opts.ShowControlPoints && controlPoints != nil {
		controlXYs := make(plotter.XYs, len(controlPoints))
		for i, point := range controlPoints {
			controlXYs[i].X = point.X
			controlXYs[i].Y = point.Y
		}

		// Create scatter plot for control points
		controlScatter, err := plotter.NewScatter(controlXYs)
		if err != nil {
			return fmt.Errorf("error creating control points scatter: %v", err)
		}
		controlScatter.Color = color.RGBA{R: 255, A: 255}
		controlScatter.Radius = vg.Points(3)
		p.Add(controlScatter)

		// Add control point labels
		for i, xy := range controlXYs {
			label, err := plotter.NewLabels(plotter.XYLabels{
				XYs:    []plotter.XY{{X: xy.X, Y: xy.Y}},
				Labels: []string{fmt.Sprintf("CP%d", i)},
			})
			if err != nil {
				return fmt.Errorf("error creating control point label: %v", err)
			}
			p.Add(label)
		}
	}

	// Add grid if enabled
	if opts.ShowGrid {
		p.Add(plotter.NewGrid())
	}

	// Save the plot
	if err := p.Save(vg.Length(opts.Width)*vg.Inch,
		vg.Length(opts.Height)*vg.Inch, opts.OutputPath); err != nil {
		return fmt.Errorf("error saving plot: %v", err)
	}

	return nil
}
