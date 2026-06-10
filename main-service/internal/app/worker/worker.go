package worker

import (
	"context"
	"time"

	"github.com/SintroSecurity/go-libraries/worker"
)

const (
	testHandler = "test"
)

func CreateAndRegisterHandlers(ctx context.Context, config *worker.Config, serviceName string) error {
	w := worker.NewSimple(ctx, config, serviceName)

	//Test background handler
	err := w.Register(testHandler, func(args worker.Args) error {
		return nil
	})
	if err != nil {
		return err
	}

	if err := w.Start(ctx); err != nil {
		return err
	}

	go func() {
		for {
			start := time.Now()
			w.Perform(worker.Job{
				Handler: testHandler,
			})
			<-time.After(1*time.Minute - time.Since(start))
		}
	}()

	return nil
}
