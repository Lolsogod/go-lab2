package pool

import (
	"awesomeProject/internal/eratosthenes"
	"sync"
)

type Task struct {
	Prime int
	Start int
	End   int
}

func Run(baseLimit, upperLimit int, numWorkers int) []int {
	if upperLimit <= baseLimit {
		return []int{}
	}

	basePrimes := eratosthenes.Default(baseLimit)
	results := make([]bool, upperLimit-baseLimit)
	for i := range results {
		results[i] = true
	}

	taskChan := make(chan Task, len(basePrimes))
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(taskChan, &wg, results, baseLimit+1)
	}

	for _, prime := range basePrimes {
		taskChan <- Task{
			Prime: prime,
			Start: baseLimit + 1,
			End:   upperLimit,
		}
	}
	close(taskChan)

	wg.Wait()

	var primes []int
	for i := range results {
		if results[i] {
			primes = append(primes, baseLimit+1+i)
		}
	}

	return primes
}

func worker(tasks <-chan Task, wg *sync.WaitGroup, results []bool, offset int) {
	defer wg.Done()

	for task := range tasks {
		firstMultiple := ((task.Start + task.Prime - 1) / task.Prime) * task.Prime
		if firstMultiple < task.Start {
			firstMultiple += task.Prime
		}

		for j := firstMultiple; j <= task.End; j += task.Prime {
			if j > task.Start {
				results[j-offset] = false
			}
		}
	}
}
