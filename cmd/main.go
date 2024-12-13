package main

import (
	"awesomeProject/internal/benchmark"
	"awesomeProject/internal/config"
	"awesomeProject/internal/executor"
	"fmt"
	"github.com/joho/godotenv"
	"log"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Ошибка загрузки конфигурации:", err)
	}

	if cfg.BenchMode {
		fmt.Printf("Starting benchmark with base limit: %d, upper limit: %d\n", cfg.BaseLimit, cfg.UpperLimit)
		fmt.Printf("Testing with worker counts: %v\n\n", benchmark.WorkerCounts)

		stats := benchmark.RunBenchmark(cfg.BaseLimit, cfg.UpperLimit)
		benchmark.PrintResults(stats)

		if cfg.SaveCsv {
			if err := benchmark.SaveResultsToCSV(stats, "benchmark_results.csv"); err != nil {
				fmt.Printf("Error saving results: %v\n", err)
			}
		}
	} else {
		var results = executor.Execute(cfg.AlgName, cfg.BaseLimit, cfg.UpperLimit, cfg.ThreadsAmount)
		fmt.Printf("Results: %v\n", len(results))
	}
}
