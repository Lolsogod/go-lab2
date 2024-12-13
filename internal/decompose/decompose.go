package decomposition

import (
	"awesomeProject/internal/eratosthenes"
	"sync"
)

func Decompose(baseLimit, upperLimit int, numWorkers int) []int {
	if upperLimit <= baseLimit {
		return []int{}
	}

	basePrimes := eratosthenes.Default(baseLimit)

	resultChan := make(chan []int, numWorkers)
	var wg sync.WaitGroup

	step := (upperLimit - baseLimit) / numWorkers
	if step == 0 {
		step = 1
	}

	for i := 0; i < numWorkers; i++ {
		start := baseLimit + 1 + i*step
		end := start + step - 1
		if i == numWorkers-1 {
			end = upperLimit
		}
		wg.Add(1)
		go worker(start, end, basePrimes, resultChan, &wg)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	var results []int
	for primes := range resultChan {
		results = append(results, primes...)
	}

	return results
}

func worker(start, end int, basePrimes []int, resultChan chan<- []int, wg *sync.WaitGroup) {
	defer wg.Done()
	localPrimes := make([]int, 0, end-start+1)

	for i := start; i <= end; i++ {
		if isPrime(i, basePrimes) {
			localPrimes = append(localPrimes, i)
		}
	}

	resultChan <- localPrimes
}

func isPrime(n int, basePrimes []int) bool {
	for _, prime := range basePrimes {
		if prime*prime > n {
			break
		}
		if n%prime == 0 {
			return false
		}
	}
	return true
}
