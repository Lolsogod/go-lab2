package executor

import (
	decomposition "awesomeProject/internal/decompose"
	"awesomeProject/internal/eratosthenes"
	"awesomeProject/internal/pool"
	"awesomeProject/internal/sequential"
	"fmt"
)

var algorithms = map[string]func(baseLimit, upperLimit int, numWorkers int) []int{
	"default": func(baseLimit, upperLimit int, numWorkers int) []int {
		return eratosthenes.Interval(baseLimit, upperLimit)
	},
	"decompose":       decomposition.Decompose,
	"decompose_prime": decomposition.DecomposePrime,
	"pool":            pool.Run,
	"sequential":      sequential.Run,
}

func Execute(algName string, baseLimit, upperLimit int, numWorkers int) []int {
	algorithm, exists := algorithms[algName]

	if !exists {
		fmt.Println("Unrecognized algorithm: " + algName)
	}
	return algorithm(baseLimit, upperLimit, numWorkers)
}
