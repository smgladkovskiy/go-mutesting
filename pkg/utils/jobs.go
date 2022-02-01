package utils

import (
	"runtime"

	"github.com/smgladkovskiy/go-mutesting/pkg/models"
)

func MaxJobs() int {
	maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()

	if maxProcs < numCPU {
		return maxProcs
	}

	return numCPU
}

func GetJobs(opts models.Options) int {
	jobs := opts.Exec.Jobs
	if jobs <= 0 {
		jobs = 1
	}

	if jobs > MaxJobs() {
		jobs = MaxJobs()
	}

	return jobs
}
