# mimic

Mimic is a Go library for generating human-like mouse movements based on various algorithms.

## Installation

```bash
go get github.com/combo23/mimic
```

## Usage

~~~go
package main

import (
	"fmt"
	"time"

	"github.com/combo23/mimic"
	"github.com/combo23/mimic/models"
)

func main() {
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

	fmt.Println(movement)
}
~~~

## Supported algorithms

-[Bézier curve](https://en.wikipedia.org/wiki/B%C3%A9zier_curve)

![Bézier curve visualization](img/bezier.png)

-[Perlin Noise](https://en.wikipedia.org/wiki/Perlin_noise)

![Perlin noise visualization](img/perlin.png)

## Contributing

Contributions are welcome! For feature requests and bug reports please submit an issue!