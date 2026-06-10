package metrics

import (
	"context"
	"errors"
	"github.com/prometheus/client_golang/prometheus"
)

type ctxKeyPrometheusType uint

const ctxKeyPrometheus ctxKeyPrometheusType = iota + 1

func GetPrometheusRegistryFromContext(ctx context.Context) (prometheus.Registerer, error) {
	maybePrometheusRegistry, ok := ctx.Value(ctxKeyPrometheus).(prometheus.Registerer)
	if !ok {
		return nil, errors.New("invalid prometheus registry provided")
	}
	return maybePrometheusRegistry, nil
}

func NewContextWithPrometheusRegistry(ctx context.Context, registry prometheus.Registerer) context.Context {
	return context.WithValue(ctx, ctxKeyPrometheus, registry)
}
