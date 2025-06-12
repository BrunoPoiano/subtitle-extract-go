package main

import (
	"fmt"
	"main/config"
	workerPool "main/pool"
	subtitleService "main/services"
	"main/types"
	"main/utils"
	"runtime"
	"sync"
)

// main is the entry point of the program that orchestrates the subtitle extraction process.
// It sets up worker pools for parallel processing and waits for all jobs to complete.
func main() {

	var wg sync.WaitGroup
	cpus_available := runtime.NumCPU()

	if cpus_available > 1 {
		cpus_available--
	}

	utils.GenerateLogs(fmt.Sprintf("INFO | using %d workers", cpus_available))
	jobs := make(chan types.SubtitleJob)

	// Initialize worker pool to process subtitle extraction jobs concurrently
	workerPool.Run(jobs, cpus_available, &wg)

	// Search for video files and enqueue extraction jobs
	subtitleService.GetItemsFromFolder(config.RootFolder, jobs, &wg)

	close(jobs) // Signal that no more jobs will be added

	wg.Wait() // Wait for all extraction jobs to complete
}
