package tests

import (
	"testing"
	"time"

	"github.com/combo23/mimic"
	"github.com/combo23/mimic/types"
)

func TestPerlin(t *testing.T) {
	opts := types.MovementOptions{
		StartPoint:    types.Point{X: 0, Y: 0},
		EndPoint:      types.Point{X: 1920, Y: 1080},
		NoiseLevel:    1,
		Resolution:    types.Point{X: 1920, Y: 1080},
		ControlPoints: 20,
		Speed:         1000,
	}

	mimic := mimic.NewMimic(mimic.PerlinAlgorithm)
	mimic.GenerateMovement(opts)
	mimic.AddHesitation(0.1, 100*time.Millisecond)
	movement := mimic.AddAcceleration(0.8, 1.2)

	err := visualize(*movement, movement.ControlPoints, defaultVisualizationOptions())
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Bezier curve generated successfully")
}