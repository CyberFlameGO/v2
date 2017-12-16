// Copyright 2017 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package scheduler

import (
	"time"

	"github.com/miniflux/miniflux/logger"
	"github.com/miniflux/miniflux/storage"
)

// NewScheduler starts a new scheduler that push jobs to a pool of workers.
func NewScheduler(store *storage.Storage, workerPool *WorkerPool, frequency, batchSize int) {
	go func() {
		c := time.Tick(time.Duration(frequency) * time.Minute)
		for now := range c {
			jobs, err := store.NewBatch(batchSize)
			if err != nil {
				logger.Error("[Scheduler] %v", err)
			} else {
				logger.Debug("[Scheduler:%v] => Pushing %d jobs", now, len(jobs))
				workerPool.Push(jobs)
			}
		}
	}()
}
