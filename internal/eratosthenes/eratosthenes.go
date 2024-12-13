package eratosthenes

import (
	"math"
)

// Default Классический алгоритм решета Эратосфена
func Default(limit int) []int {
	isPrime := make([]bool, limit+1)
	for i := 2; i <= limit; i++ {
		isPrime[i] = true
	}

	for i := 2; i*i <= limit; i++ {
		if isPrime[i] {
			for j := i * i; j <= limit; j += i {
				isPrime[j] = false
			}
		}
	}

	var primes []int
	for i := 2; i <= limit; i++ {
		if isPrime[i] {
			primes = append(primes, i)
		}
	}
	return primes
}

func Interval(start, end int) []int {
	baseLimit := int(math.Sqrt(float64(end)))
	basePrimes := Default(baseLimit)

	isPrime := make([]bool, end-start+1)
	for i := range isPrime {
		isPrime[i] = true
	}

	for _, prime := range basePrimes {
		firstMultiple := max(prime*prime, ((start+prime-1)/prime)*prime)
		for j := firstMultiple; j <= end; j += prime {
			if j >= start {
				isPrime[j-start] = false
			}
		}
	}

	var result []int
	for i := 0; i <= end-start; i++ {
		if isPrime[i] && (i+start) > 1 {
			result = append(result, i+start)
		}
	}
	return result
}
