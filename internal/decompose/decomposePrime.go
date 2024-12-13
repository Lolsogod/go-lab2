package decomposition

import (
	"awesomeProject/internal/eratosthenes"
	"sync"
)

func DecomposePrime(baseLimit, upperLimit int, numWorkers int) []int {
	if upperLimit <= baseLimit {
		return []int{}
	}

	basePrimes := eratosthenes.Default(baseLimit)
	if len(basePrimes) == 0 {
		return []int{}
	}

	results := make([]bool, upperLimit-baseLimit+1)
	for i := range results {
		results[i] = true
	}

	var wg sync.WaitGroup

	effectiveWorkers := min(numWorkers, len(basePrimes))
	if effectiveWorkers == 0 {
		return []int{}
	}

	primesPerWorker := len(basePrimes) / effectiveWorkers
	if primesPerWorker == 0 {
		primesPerWorker = 1
		effectiveWorkers = len(basePrimes)
	}

	for i := 0; i < effectiveWorkers; i++ {
		start := i * primesPerWorker
		end := start + primesPerWorker
		if i == effectiveWorkers-1 {
			end = len(basePrimes)
		}

		if start >= len(basePrimes) {
			break
		}

		wg.Add(1)
		go func(primes []int) {
			defer wg.Done()
			sieveWithPrimes(baseLimit+1, upperLimit, primes, results)
		}(basePrimes[start:end])
	}

	wg.Wait()

	var primes []int
	for i := range results {
		if results[i] {
			primes = append(primes, baseLimit+1+i)
		}
	}

	return primes
}

func sieveWithPrimes(start, end int, primes []int, results []bool) {
	if len(primes) == 0 {
		return
	}

	for _, prime := range primes {
		if prime == 0 {
			continue
		}

		firstMultiple := ((start + prime - 1) / prime) * prime
		if firstMultiple < start {
			firstMultiple += prime
		}

		for j := firstMultiple; j <= end; j += prime {
			if j > start && (j-start) < len(results) {
				results[j-start] = false
			}
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
