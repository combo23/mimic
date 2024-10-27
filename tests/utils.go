package tests

type VisualizationOptions struct {
	ShowControlPoints bool
	ShowGrid          bool
	ShowArrows        bool
	Title             string
	OutputPath        string
	Width             float64 // in inches
	Height            float64 // in inches
}

// defaultVisualizationOptions returns default visualization settings
func defaultVisualizationOptions() VisualizationOptions {
	return VisualizationOptions{
		ShowControlPoints: true,
		ShowGrid:          true,
		ShowArrows:        true,
		Title:             "Mouse movements visualization",
		OutputPath:        "mouse_movement.png",
		Width:             10,
		Height:            6,
	}
}
