package worker

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/SintroSecurity/go-libraries/logger"
	"github.com/SintroSecurity/go-libraries/metrics"
)

var _ Worker = &Simple{}

// NewSimple creates a basic implementation of the Worker interface
// that is backed using just the standard library and goroutines.
func NewSimple(ctx context.Context, config *Config, serviceName string) *Simple {
	ctx, cancel := context.WithCancel(ctx)
	l := logger.GetLoggerFromContext(ctx)
	registry, err := metrics.GetPrometheusRegistryFromContext(ctx)
	if err != nil {
		panic(err)
	}

	newWorkerMetrics(serviceName, registry, dflBuckets...)

	return &Simple{
		Config:   config,
		Logger:   l,
		ctx:      ctx,
		cancel:   cancel,
		handlers: map[string]Handler{},
		moot:     &sync.Mutex{},
		started:  false,
	}
}

// Simple is a basic implementation of the Worker interface
// that is backed using just the standard library and goroutines.
type Simple struct {
	Config   *Config
	Logger   *logger.Logger
	ctx      context.Context
	cancel   context.CancelFunc
	handlers map[string]Handler
	moot     *sync.Mutex
	wg       sync.WaitGroup
	started  bool
}

// Register Handler with the worker
func (w *Simple) Register(name string, h Handler) error {
	if name == "" || h == nil {
		return fmt.Errorf("name or handler cannot be empty/nil")
	}

	w.moot.Lock()
	defer w.moot.Unlock()
	if _, ok := w.handlers[name]; ok {
		return fmt.Errorf("handler already mapped for name %s", name)
	}
	w.handlers[name] = h
	return nil
}

// Start the worker
func (w *Simple) Start(ctx context.Context) error {
	// TODO(sio4): #road-to-v1 - define the purpose of Start clearly
	if w.Config.Enabled {
		w.Logger.Info("starting Simple background worker")

		w.moot.Lock()
		defer w.moot.Unlock()

		w.ctx, w.cancel = context.WithCancel(ctx)
		w.started = true
		return nil
	} else {
		w.Logger.Info("Simple background worker is disabled")
		return nil
	}
}

// Stop the worker
func (w *Simple) Stop() error {
	// prevent job submission when stopping
	w.moot.Lock()
	defer w.moot.Unlock()

	w.Logger.Info("stopping Simple background worker")

	w.cancel()

	w.wg.Wait()
	w.Logger.Info("all background jobs stopped completely")
	return nil
}

// Perform a job as soon as possibly using a goroutine.
func (w *Simple) Perform(job Job) error {
	if job.Handler == "" {
		w.Logger.Error("no handler name given", w.Logger.Any("job", job))
		return fmt.Errorf("no handler name given: %s", job)
	}
	enqueuedJobs.WithLabelValues(job.Handler).Inc()

	w.moot.Lock()
	defer w.moot.Unlock()

	if !w.started {
		return fmt.Errorf("worker is not yet started")
	}

	// Perform should not allow a job submission if the worker is not running
	if err := w.ctx.Err(); err != nil {
		return fmt.Errorf("worker is not ready to perform a job: %v", err)
	}

	w.Logger.Debug("performing job", w.Logger.Any("job", job))

	if h, ok := w.handlers[job.Handler]; ok {
		// TODO: consider timeout and/or cancellation
		w.wg.Add(1)
		go func() {
			start := time.Now()
			defer w.wg.Done()
			defer processedJobs.WithLabelValues(job.Handler).Inc()
			err := safeRun(func() error {
				return h(job.Args)
			})

			if err != nil {
				w.Logger.Error("error running job", w.Logger.Err(err), w.Logger.Any("job", job))
			}
			processedJobsLatency.WithLabelValues(job.Handler).Observe(float64(time.Since(start).Nanoseconds()) / 1000000)
			w.Logger.Debug("completed job", w.Logger.Any("job", job))
		}()
		w.wg.Wait()
		return nil
	}

	w.Logger.Error("no handler mapped for name", w.Logger.Any("job", job))
	return fmt.Errorf("no handler mapped for name %s", job.Handler)
}

// safeRun the function safely knowing that if it panics
// the panic will be caught and returned as an error
func safeRun(fn func() error) (err error) {
	defer func() {
		if ex := recover(); ex != nil {
			if e, ok := ex.(error); ok {
				err = e
				return
			}
			err = errors.New(fmt.Sprint(ex))
		}
	}()

	return fn()
}

// PerformAt performs a job at a particular time using a goroutine.
func (w *Simple) PerformAt(job Job, t time.Time) error {
	return w.PerformIn(job, time.Until(t))
}

// PerformIn performs a job after waiting for a specified amount
// using a goroutine.
func (w *Simple) PerformIn(job Job, d time.Duration) error {
	// Perform should not allow a job submission if the worker is not running
	if err := w.ctx.Err(); err != nil {
		return fmt.Errorf("worker is not ready to perform a job: %v", err)
	}

	w.wg.Add(1) // waiting job also should be counted
	go func() {
		defer w.wg.Done()

		for {
			w.moot.Lock()
			if w.started {
				w.moot.Unlock()
				break
			}
			w.moot.Unlock()

			waiting := 100 * time.Millisecond
			time.Sleep(waiting)
			d = d - waiting
		}

		select {
		case <-time.After(d):
			w.Perform(job)
		case <-w.ctx.Done():
			// TODO(sio4): #road-to-v1 - it should be guaranteed to be performed
			w.cancel()
		}
	}()
	return nil
}
