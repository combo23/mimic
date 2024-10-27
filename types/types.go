package types

import "time"

type Point struct {
	X, Y float64
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
	Timing        []time.Duration
	ControlPoints []Point
}

type EffectType int

const (
	HesitateEffect EffectType = iota
	AccelerateEffect
)
