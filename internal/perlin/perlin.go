package perlin

import (
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/combo23/mimic/models"
)

type Perlin struct {
	models.Movement
	permutation []int
	gradients   []models.Point
}

// generatePermutation creates a random permutation table for Perlin noise
func (p *Perlin) generatePermutation() {
	p.permutation = make([]int, 512) // Double size for safe wrapping
	for i := 0; i < 256; i++ {
		p.permutation[i] = i
	}

	// Fisher-Yates shuffle for first 256 values
	for i := 255; i > 0; i-- {
		j := rand.Intn(i + 1)
		p.permutation[i], p.permutation[j] = p.permutation[j], p.permutation[i]
	}

	// Duplicate the permutation for safe wrapping
	for i := 0; i < 256; i++ {
		p.permutation[256+i] = p.permutation[i]
	}
}

// generateGradients creates random unit vectors for gradient calculation
func (p *Perlin) generateGradients() {
	p.gradients = make([]models.Point, 256)
	for i := range p.gradients {
		angle := rand.Float64() * 2 * math.Pi
		p.gradients[i] = models.Point{
			X: math.Cos(angle),
			Y: math.Sin(angle),
		}
	}
}

// fade applies Ken Perlin's fade function: 6t^5 - 15t^4 + 10t^3
func fade(t float64) float64 {
	return t * t * t * (t*(t*6-15) + 10)
}

// lerp performs linear interpolation between a and b
func lerp(a, b, t float64) float64 {
	return a + t*(b-a)
}

// hash creates a repeatable hash of x and y coordinates
func (p *Perlin) hash(x, y int) int {
	return p.permutation[p.permutation[x&255]+(y&255)&255]
}

// generateControlPoints creates control points along a Perlin noise path
func (p *Perlin) generateControlPoints(opts models.MovementOptions) []models.Point {
	if p.permutation == nil {
		p.generatePermutation()
		p.generateGradients()
	}

	points := make([]models.Point, opts.ControlPoints)
	points[0] = opts.StartPoint
	points[len(points)-1] = opts.EndPoint

	// Calculate the main direction vector
	dirX := opts.EndPoint.X - opts.StartPoint.X
	dirY := opts.EndPoint.Y - opts.StartPoint.Y

	// Use multiple frequencies of Perlin noise
	scale := math.Min(opts.Resolution.X, opts.Resolution.Y) * opts.NoiseLevel

	var wg sync.WaitGroup
	for i := 1; i < len(points)-1; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			// Calculate base position along the direct path
			progress := float64(i) / float64(len(points)-1)
			baseX := opts.StartPoint.X + dirX*progress
			baseY := opts.StartPoint.Y + dirY*progress

			// Generate Perlin noise offset
			t := progress * 4 // Scale time for more variation
			offsetX := p.perlinNoise2D(t, 0.0) * scale
			offsetY := p.perlinNoise2D(t, 1.0) * scale

			points[i] = models.Point{
				X: baseX + offsetX,
				Y: baseY + offsetY,
			}
		}(i)
	}
	wg.Wait()

	return points
}

// perlinNoise2D generates 2D Perlin noise value
func (p *Perlin) perlinNoise2D(x, y float64) float64 {
	// Get grid cell coordinates
	x0 := int(math.Floor(x)) & 255
	y0 := int(math.Floor(y)) & 255

	// Get relative coords within cell
	xf := x - math.Floor(x)
	yf := y - math.Floor(y)

	// Compute fade curves
	u := fade(xf)
	v := fade(yf)

	// Get hashed gradient indices
	aa := p.hash(x0, y0)
	ab := p.hash(x0, y0+1)
	ba := p.hash(x0+1, y0)
	bb := p.hash(x0+1, y0+1)

	// Get gradients
	g00 := p.gradients[aa%256]
	g10 := p.gradients[ba%256]
	g01 := p.gradients[ab%256]
	g11 := p.gradients[bb%256]

	// Calculate dot products
	d00 := g00.X*xf + g00.Y*yf
	d10 := g10.X*(xf-1) + g10.Y*yf
	d01 := g01.X*xf + g01.Y*(yf-1)
	d11 := g11.X*(xf-1) + g11.Y*(yf-1)

	// Interpolate dot products
	x0Interp := lerp(d00, d10, u)
	x1Interp := lerp(d01, d11, u)

	return lerp(x0Interp, x1Interp, v)
}

// interpolatePoints generates points between control points using Perlin noise
func (p *Perlin) interpolatePoints(points []models.Point, numPoints int, noiseLevel float64) []models.Point {
	result := make([]models.Point, numPoints)
	segmentSize := float64(len(points)-1) / float64(numPoints-1)

	for i := 0; i < numPoints; i++ {
		t := float64(i) * segmentSize
		segment := int(t)
		if segment >= len(points)-1 {
			segment = len(points) - 2
		}

		localT := t - float64(segment)

		// Base interpolation between control points
		p1 := points[segment]
		p2 := points[segment+1]

		baseX := lerp(p1.X, p2.X, localT)
		baseY := lerp(p1.Y, p2.Y, localT)

		noiseScale := noiseLevel * 0.1 // Reduce noise scale for smoother movement
		noiseX := p.perlinNoise2D(float64(i)*0.1, 0.0) * noiseScale
		noiseY := p.perlinNoise2D(float64(i)*0.1, 1.0) * noiseScale

		result[i] = models.Point{
			X: baseX + noiseX,
			Y: baseY + noiseY,
		}
	}

	return result
}

// GenerateMovement creates a complete mouse movement path using Perlin noise
func (p *Perlin) GenerateMovement(opts models.MovementOptions) *models.Movement {
	if opts.ControlPoints < 2 {
		opts.ControlPoints = 2
	}
	if opts.NoiseLevel < 0 {
		opts.NoiseLevel = 0
	}
	if opts.NoiseLevel > 1 {
		opts.NoiseLevel = 1
	}

	controlPoints := p.generateControlPoints(opts)

	// Calculate approximate path length
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

	// Generate interpolated points with additional Perlin noise
	points := p.interpolatePoints(controlPoints, numPoints, opts.NoiseLevel*5)

	// Create movement with timing
	movement := models.Movement{
		Points:        points,
		ControlPoints: controlPoints,
	}

	// Calculate timing
	prevPoint := opts.StartPoint

	for i := 0; i < numPoints; i++ {
		point := points[i]

		// Calculate timing based on distance from previous point
		dx := point.X - prevPoint.X
		dy := point.Y - prevPoint.Y
		distance := math.Sqrt(dx*dx + dy*dy)

		// Add some randomness to timing
		timeForSegment := time.Duration(float64(time.Second) * distance / opts.Speed)
		timeForSegment += time.Duration(rand.Float64() * float64(timeForSegment) * 0.1)

		points[i].Timing = timeForSegment

		prevPoint = point
	}

	p.Movement = movement
	return &p.Movement
}
