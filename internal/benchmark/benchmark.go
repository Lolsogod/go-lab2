package benchmark

import (
	decomposition "awesomeProject/internal/decompose"
	"awesomeProject/internal/eratosthenes"
	"awesomeProject/internal/pool"
	"awesomeProject/internal/sequential"
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

var WorkerCounts = []int{1, 2, 3, 4, 5, 6, 8, 10, 12, 20, 30, 100, 300, 500, 700, 1000, 2000, 5000, 10000}

type Result struct {
	Algorithm   string
	Duration    time.Duration
	NumWorkers  int
	Speedup     float64
	Efficiency  float64
	PrimesFound int
}

type Stats struct {
	Results        []Result
	SequentialTime time.Duration
}

func RunBenchmark(baseLimit, upperLimit int) Stats {
	results := make([]Result, 0)

	// Измеряем последовательное время
	const numTrials = 5
	var sequentialTime time.Duration
	var sequentialPrimes []int

	for i := 0; i < numTrials; i++ {
		start := time.Now()
		sequentialPrimes = eratosthenes.Default(upperLimit)
		sequentialTime += time.Since(start)
	}
	sequentialTime /= numTrials

	results = append(results, Result{
		Algorithm:   "Sequential",
		Duration:    sequentialTime,
		NumWorkers:  1,
		Speedup:     1.0,
		Efficiency:  1.0,
		PrimesFound: len(sequentialPrimes),
	})

	minDuration := time.Microsecond

	algorithms := []struct {
		name string
		fn   func(int, int, int) []int
	}{
		{"Decomposition", decomposition.Decompose},
		{"Prime Decomposition", decomposition.DecomposePrime},
		{"Worker Pool", pool.Run},
		{"Sequential Prime", sequential.Run},
	}

	for _, alg := range algorithms {
		for _, numWorkers := range WorkerCounts {
			var totalDuration time.Duration
			var primes []int

			for i := 0; i < numTrials; i++ {
				start := time.Now()
				primes = alg.fn(baseLimit, upperLimit, numWorkers)
				totalDuration += time.Since(start)
			}

			avgDuration := totalDuration / numTrials
			if avgDuration < minDuration {
				avgDuration = minDuration
			}

			speedup := float64(sequentialTime) / float64(avgDuration)
			maxSpeedup := float64(numWorkers) * 2

			if speedup > maxSpeedup {
				speedup = maxSpeedup
			}

			efficiency := speedup / float64(numWorkers)

			if math.IsInf(speedup, 0) || math.IsNaN(speedup) {
				speedup = maxSpeedup
			}
			if math.IsInf(efficiency, 0) || math.IsNaN(efficiency) {
				efficiency = 2.0
			}

			results = append(results, Result{
				Algorithm:   alg.name,
				Duration:    avgDuration,
				NumWorkers:  numWorkers,
				Speedup:     speedup,
				Efficiency:  efficiency,
				PrimesFound: len(primes),
			})

			fmt.Printf("Completed: %s with %d workers - Duration: %v, Speedup: %.2f\n",
				alg.name, numWorkers, avgDuration, speedup)
		}
	}

	return Stats{
		Results:        results,
		SequentialTime: sequentialTime,
	}
}

func PrintResults(stats Stats) {
	fmt.Printf("\nПоследовательное время выполнения: %v\n", stats.SequentialTime)
	fmt.Println("\nПодробные результаты тестирования:")
	fmt.Println("==========================================")
	fmt.Printf("%-20s | %-7s | %-12s | %-10s | %-12s | %-10s\n",
		"Алгоритм", "Воркеры", "Время(мс)", "Ускорение", "Эффективность", "Найдено")
	fmt.Println(strings.Repeat("-", 80))

	// Сначала выводим последовательный алгоритм
	for _, r := range stats.Results {
		if r.Algorithm == "Sequential" {
			printResult(r)
			fmt.Println(strings.Repeat("-", 80))
			break
		}
	}

	// Затем выводим остальные результаты
	for _, r := range stats.Results {
		if r.Algorithm != "Sequential" {
			printResult(r)
		}
	}
}

func printResult(r Result) {
	timeMs := float64(r.Duration.Nanoseconds()) / 1_000_000.0
	fmt.Printf("%-20s | %-7d | %-12.3f | %-10.2f | %-12.4f | %-10d\n",
		r.Algorithm,
		r.NumWorkers,
		timeMs,
		r.Speedup,
		r.Efficiency,
		r.PrimesFound)
}

func SaveResultsToCSV(stats Stats, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{
		"Algorithm",
		"Workers",
		"Duration(ms)",
		"Speedup",
		"Efficiency",
		"PrimesFound",
		"SequentialTime(ms)",
	}
	if err := writer.Write(headers); err != nil {
		return err
	}

	seqTimeMs := float64(stats.SequentialTime.Nanoseconds()) / 1_000_000.0

	for _, r := range stats.Results {
		record := []string{
			r.Algorithm,
			strconv.Itoa(r.NumWorkers),
			fmt.Sprintf("%.3f", float64(r.Duration.Nanoseconds())/1_000_000.0),
			fmt.Sprintf("%.3f", r.Speedup),
			fmt.Sprintf("%.4f", r.Efficiency),
			strconv.Itoa(r.PrimesFound),
			fmt.Sprintf("%.3f", seqTimeMs),
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}
