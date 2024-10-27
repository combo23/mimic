package bezier

import (
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/combo23/mimic/types"
)

type Bezier struct {
	types.Movement
}

// generateControlPoints creates random control points for the Bezier curve
func generateControlPoints(opts types.MovementOptions) []types.Point {
	points := make([]types.Point, opts.ControlPoints)
	points[0] = opts.StartPoint
	points[len(points)-1] = opts.EndPoint

	// Calculate the main direction vector
	dirX := opts.EndPoint.X - opts.StartPoint.X
	dirY := opts.EndPoint.Y - opts.StartPoint.Y

	// Generate intermediate control points
	for i := 1; i < len(points)-1; i++ {
		// Calculate base position along the direct path
		progress := float64(i) / float64(len(points)-1)
		baseX := opts.StartPoint.X + dirX*progress
		baseY := opts.StartPoint.Y + dirY*progress

		// Add noise based on noise level
		maxOffset := math.Min(opts.Resolution.X, opts.Resolution.Y) * opts.NoiseLevel
		offsetX := (rand.Float64()*2 - 1) * maxOffset
		offsetY := (rand.Float64()*2 - 1) * maxOffset

		points[i] = types.Point{
			X: baseX + offsetX,
			Y: baseY + offsetY,
		}
	}

	return points
}

func bezierCurveSequential(t float64, points []types.Point) types.Point {
	n := len(points) - 1
	var x, y float64
	for i := 0; i <= n; i++ {
		coeff := float64(binomialCoeff(n, i))
		basis := coeff * math.Pow(t, float64(i)) * math.Pow(1-t, float64(n-i))
		x += points[i].X * basis
		y += points[i].Y * basis
	}
	return types.Point{X: x, Y: y}
}

// Improved parallel implementation
func bezierCurveParallel(t float64, points []types.Point) types.Point {
	n := len(points) - 1

	// For small number of points, use sequential version
	if n < 4 {
		return bezierCurveSequential(t, points)
	}

	var wg sync.WaitGroup
	results := make(chan types.Point, n+1)

	// Pre-calculate common values
	tPow := make([]float64, n+1)
	oneMinusTPow := make([]float64, n+1)
	tPow[0] = 1
	oneMinusTPow[0] = 1

	for i := 1; i <= n; i++ {
		tPow[i] = tPow[i-1] * t
		oneMinusTPow[i] = oneMinusTPow[i-1] * (1 - t)
	}

	// Launch workers
	for i := 0; i <= n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			coeff := float64(binomialCoeff(n, i))
			basis := coeff * tPow[i] * oneMinusTPow[n-i]

			results <- types.Point{
				X: points[i].X * basis,
				Y: points[i].Y * basis,
			}
		}(i)
	}

	// Close results channel after all goroutines complete
	go func() {
		wg.Wait()
		close(results)
	}()

	// Sum up results
	var finalPoint types.Point
	for result := range results {
		finalPoint.X += result.X
		finalPoint.Y += result.Y
	}

	return finalPoint
}

// binomialCoeff calculates the binomial coefficient (n choose k)
func binomialCoeff(n, k int) int {
	if k == 0 || k == n {
		return 1
	}

	if k > n {
		return 0
	}

	return binomialCoeff(n-1, k-1) + binomialCoeff(n-1, k)
}

// GenerateMovement creates a complete mouse movement path
func (b *Bezier) GenerateMovement(opts types.MovementOptions) *types.Movement {
	if opts.ControlPoints < 2 {
		opts.ControlPoints = 2
	}
	if opts.NoiseLevel < 0 {
		opts.NoiseLevel = 0
	}
	if opts.NoiseLevel > 1 {
		opts.NoiseLevel = 1
	}

	// Generate control points
	controlPoints := generateControlPoints(opts)

	// Calculate path length (approximate)
	totalDistance := 0.0
	for i := 0; i < len(controlPoints)-1; i++ {
		dx := controlPoints[i+1].X - controlPoints[i].X
		dy := controlPoints[i+1].Y - controlPoints[i].Y
		totalDistance += math.Sqrt(dx*dx + dy*dy)
	}

	// Calculate number of points based on distance and speed
	numPoints := int(totalDistance/opts.Speed*60) + 1
	if numPoints < 2 {
		numPoints = 2
	}

	// Generate points along the curve
	movement := types.Movement{
		Points:        make([]types.Point, numPoints),
		Timing:        make([]time.Duration, numPoints),
		ControlPoints: controlPoints,
	}

	prevPoint := opts.StartPoint
	totalTime := time.Duration(0)

	for i := 0; i < numPoints; i++ {
		t := float64(i) / float64(numPoints-1)
		point := bezierCurveParallel(t, controlPoints)

		// Add subtle noise to each point
		if opts.NoiseLevel > 0 {
			microNoise := opts.NoiseLevel * 0.1
			point.X += (rand.Float64()*2 - 1) * microNoise
			point.Y += (rand.Float64()*2 - 1) * microNoise
		}

		// Calculate timing based on distance from previous point
		dx := point.X - prevPoint.X
		dy := point.Y - prevPoint.Y
		distance := math.Sqrt(dx*dx + dy*dy)

		// Add some randomness to timing
		timeForSegment := time.Duration(float64(time.Second) * distance / opts.Speed)
		timeForSegment += time.Duration(rand.Float64() * float64(timeForSegment) * 0.1) // 10% random variation

		totalTime += timeForSegment
		movement.Points[i] = point
		movement.Timing[i] = totalTime

		prevPoint = point
	}

	b.Movement = movement
	return &b.Movement
}

// AddHesitation adds random pauses to the movement
func (b *Bezier) AddHesitation(probability float64, maxPause time.Duration) *types.Movement {
	m := &b.Movement

	for i := 1; i < len(m.Timing)-1; i++ {
		if rand.Float64() < probability {
			pause := time.Duration(rand.Float64() * float64(maxPause))
			for j := i; j < len(m.Timing); j++ {
				m.Timing[j] += pause
			}
		}
	}
	return &b.Movement
}

// AddAcceleration modifies timing to simulate acceleration/deceleration
func (b *Bezier) AddAcceleration(startSpeed, endSpeed float64) *types.Movement {
	m := &b.Movement

	totalTime := m.Timing[len(m.Timing)-1]
	for i := 0; i < len(m.Timing); i++ {
		progress := float64(i) / float64(len(m.Timing)-1)
		speedFactor := startSpeed + (endSpeed-startSpeed)*progress
		m.Timing[i] = time.Duration(float64(totalTime) * speedFactor)
	}

	return &b.Movement
}
