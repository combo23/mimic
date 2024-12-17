package models

import (
	"fmt"
	"time"
)

type Point struct {
	X, Y   float64
	Timing time.Duration
}

func (p *Point) String() string {
	return fmt.Sprintf("(%f, %f) took %v", p.X, p.Y, p.Timing)
}

type MovementOptions struct {
	StartPoint    Point   // Starting coordinates
	EndPoint      Point   // Ending coordinates
	NoiseLevel    float64 // Amount of randomness (0.0 to 1.0)
	Resolution    Point   // Screen resolution (width x height)
	ControlPoints int     // Number of control points for Bezier curve (min 2)
	Speed         float64 // Movement speed in pixels per second
}

// Movement represents a complete mouse movement path
type Movement struct {
	Points        []Point
	ControlPoints []Point
}

type EffectType int

const (
	HesitateEffect EffectType = iota
	AccelerateEffect
)
