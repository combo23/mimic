package tests

import (
	"fmt"
	"testing"

	"github.com/combo23/mimic"
	"github.com/combo23/mimic/models"
)

func TestBezier(t *testing.T) {
	opts := models.MovementOptions{
		StartPoint:    models.Point{X: 0, Y: 0},
		EndPoint:      models.Point{X: 1920, Y: 1080},
		NoiseLevel:    1,
		Resolution:    models.Point{X: 1920, Y: 1080},
		ControlPoints: 20,
		Speed:         1000,
	}

	mimic := mimic.NewMimic(mimic.BezierAlgorithm)
	movement := mimic.GenerateMovement(opts)

	err := visualize(*movement, movement.ControlPoints, bezierVisualizationOptions)
	if err != nil {
		t.Fatal(err)
	}

	for _, point := range movement.Points {
		t.Log(point.String())
	}
	fmt.Println(len(movement.Points))

	t.Log("Bezier curve generated successfully")
}

var bezierVisualizationOptions = VisualizationOptions{
	ShowControlPoints: true,
	ShowGrid:          true,
	ShowArrows:        true,
	Title:             "Bezier Curve Mouse Movement",
	OutputPath:        "bezier.png",
	Width:             10,
	Height:            6,
}
