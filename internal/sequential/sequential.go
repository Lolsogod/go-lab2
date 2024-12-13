package sequential

import (
	"awesomeProject/internal/eratosthenes"
	"sync"
	"sync/atomic"
)

func Run(baseLimit, upperLimit int, numWorkers int) []int {
	if upperLimit <= baseLimit {
		return []int{}
	}

	basePrimes := eratosthenes.Default(baseLimit)
	results := make([]bool, upperLimit-baseLimit)
	for i := range results {
		results[i] = true
	}

	var currentPrimeIndex int32 = 0
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				index := atomic.AddInt32(&currentPrimeIndex, 1) - 1
				if int(index) >= len(basePrimes) {
					return
				}

				prime := basePrimes[index]
				sieveWithPrime(baseLimit+1, upperLimit, prime, results)
			}
		}()
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

func sieveWithPrime(start, end, prime int, results []bool) {
	firstMultiple := ((start + prime - 1) / prime) * prime
	if firstMultiple < start {
		firstMultiple += prime
	}

	for j := firstMultiple; j <= end; j += prime {
		if j > start {
			results[j-start] = false
		}
	}
}
