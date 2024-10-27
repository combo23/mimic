package mimic

import (
	"time"

	"github.com/combo23/mimic/internal/bezier"
	"github.com/combo23/mimic/types"
)

type AlgorithmType int

const (
	BezierAlgorithm AlgorithmType = iota
)

type Mimic interface {
	GenerateMovement(opts types.MovementOptions) *types.Movement
	AddHesitation(noiseLevel float64, duration time.Duration) *types.Movement
	AddAcceleration(startSpeed, endSpeed float64) *types.Movement
}

func NewMimic(algorithm AlgorithmType) Mimic {
	switch algorithm {
	case BezierAlgorithm:
		return &bezier.Bezier{}
	default:
		return nil
	}
}
