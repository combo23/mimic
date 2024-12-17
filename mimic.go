package mimic

import (
	"github.com/combo23/mimic/internal/bezier"
	"github.com/combo23/mimic/internal/perlin"
	"github.com/combo23/mimic/models"
)

type AlgorithmType int

const (
	BezierAlgorithm AlgorithmType = iota
	PerlinAlgorithm
)

type Mimic interface {
	GenerateMovement(opts models.MovementOptions) *models.Movement
}

func NewMimic(algorithm AlgorithmType) Mimic {
	switch algorithm {
	case BezierAlgorithm:
		return &bezier.Bezier{}
	case PerlinAlgorithm:
		return &perlin.Perlin{}
	default:
		return nil
	}
}
