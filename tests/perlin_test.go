package tests

import (
	"testing"
	"time"

	"github.com/combo23/mimic"
	"github.com/combo23/mimic/models"
)

func TestPerlin(t *testing.T) {
	opts := models.MovementOptions{
		StartPoint:    models.Point{X: 0, Y: 0},
		EndPoint:      models.Point{X: 1920, Y: 1080},
		NoiseLevel:    1,
		Resolution:    models.Point{X: 1920, Y: 1080},
		ControlPoints: 50,
		Speed:         1000,
	}

	mimic := mimic.NewMimic(mimic.PerlinAlgorithm)
	mimic.GenerateMovement(opts)
	mimic.AddHesitation(0.1, 100*time.Millisecond)
	movement := mimic.AddAcceleration(0.8, 1.2)

	err := visualize(*movement, movement.ControlPoints, perlinVisualizationOptions)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Bezier curve generated successfully")
}

var perlinVisualizationOptions = VisualizationOptions{
	ShowControlPoints: true,
	ShowGrid:          true,
	ShowArrows:        true,
	Title:             "Perlin Noise Mouse Movement",
	OutputPath:        "perlin.png",
	Width:             10,
	Height:            6,
}
