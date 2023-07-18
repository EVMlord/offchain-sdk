package worker

import (
	"fmt"
	"sync"

	"cosmossdk.io/log"
)

// worker is a worker thread that executes jobs.
type worker struct {
	id uint32
	// Gets jobs fed to it.
	newExecutor chan (Executor)
	// Feeds results onto a channel.
	newRes chan (Resultor)
	// Notify the worker to stop.
	stop chan struct{}
	// logger represents our logger
	logger log.Logger
	// wg is used to wait for the worker to stop.
	wg *sync.WaitGroup
}

// NewWorker creates a new worker.
func newWorker(
	id uint32,
	newExecutor chan Executor,
	newRes chan Resultor,
	logger log.Logger,
	wg *sync.WaitGroup,
) *worker {
	return &worker{
		id:          id,
		logger:      logger,
		newExecutor: newExecutor,
		newRes:      newRes,
		wg:          wg,
	}
}

// Logger returns the logger for the worker.
func (w *worker) Logger() log.Logger {
	return w.logger.With("namespace", fmt.Sprintf("worker-%d", w.id))
}

// Start starts the worker.
func (w *worker) Start() {
	// Manage stopping the worker.
	w.stop = make(chan struct{})
	// Add the worker to the wait group.
	// On exiting, remove the worker from the wait group.
	defer func() {
		close(w.stop)
		w.wg.Done()
	}()

	w.Logger().Info("starting")
	for {
		select {
		case <-w.stop:
			w.Logger().Info("stopping worker")
			return
		case executor, ok := <-w.newExecutor:
			if !ok {
				w.Logger().Error("worker stopped because of error")
				return
			}
			w.Logger().Info("executing job")
			w.newRes <- executor.Execute()
			continue
		case <-w.newRes:
			w.Logger().Info("received result")
			continue
		}
	}
}

// Stop stops the worker.
func (w *worker) Stop() {
	w.Logger().Info("triggering worker to stop")
	w.stop <- struct{}{}
}
