package workerPool

import (
	subtitleService "main/services"
	"main/types"
	"sync"
)

// Run creates and manages a pool of worker goroutines to process subtitle extraction jobs.
// Parameters:
//   - jobs: channel of subtitle extraction tasks
//   - workers: number of concurrent workers to spawn
//   - wg: wait group to track job completion
func Run(jobs <-chan types.SubtitleJob, workers int, wg *sync.WaitGroup) {
	for range workers {
		go func() {
			for job := range jobs {
				subtitleService.RunningEmbedSubtitleCheck(job.Location, job.Item)
				wg.Done()
			}
		}()
	}
}
