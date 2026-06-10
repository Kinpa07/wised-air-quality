package router

import (
	"fmt"
	"github.com/go-chi/cors"
	"net/http"
	"time"

	"github.com/SintroSecurity/go-libraries/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Router struct {
	chi.Router
	mainRouter chi.Router
	listenUrl  string
}

const RequestIDHeader = "X-Request-ID"

func NewRouter(config *Config, loggerConfig *logger.Config, prometheusRegistry prometheus.Registerer, serviceName string, middlewares ...func(next http.Handler) http.Handler) *Router {
	r := chi.NewRouter()
	InitGoRuntimeMetrics(prometheusRegistry)
	r.Route("/", func(r chi.Router) {
		r.Handle("/metrics", promhttp.HandlerFor(
			prometheusRegistry.(*prometheus.Registry),
			promhttp.HandlerOpts{
				EnableOpenMetrics: true,
			},
		))
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			render.Status(r, 200)
			render.JSON(w, r, "OK")
		})
	})

	requestTimeMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestId := middleware.GetReqID(r.Context())
			ctx := r.Context()
			log := logger.New(ctx, loggerConfig)
			log = log.With(log.String("request_id", requestId))
			ctx = logger.NewContextWithLogger(r.Context(), log)
			scheme := "http"
			if r.TLS != nil {
				scheme = "https"
			}
			log.Info("request started", log.String("method", r.Method), log.String("URL", fmt.Sprintf("%s://%s%s %s\"", scheme, r.Host, r.RequestURI, r.Proto)), log.Int64("content_length", r.ContentLength))

			startTime := time.Now()
			rw := NewResponseWriter(w)
			rw.Header().Set(RequestIDHeader, requestId)
			next.ServeHTTP(rw, r.WithContext(ctx))
			duration := time.Since(startTime)
			log.Info("request completed", log.String("method", r.Method),
				log.String("URL", fmt.Sprintf("%s://%s%s %s\"", scheme, r.Host, r.RequestURI, r.Proto)),
				log.Int64("duration", duration.Milliseconds()), log.Int("status", rw.Status()), log.Int("size", rw.Size()))
		})
	}

	return &Router{
		mainRouter: r,
		Router: r.Route(fmt.Sprintf("/api/%s", serviceName), func(r chi.Router) {
			// Conditionally apply CORS middleware if overridden in config
			if config.CORS != nil {
				r.Use(cors.Handler(NewCORSOptions(config.CORS)))
			}
			r.Use(NewPrometheusPatternMiddleware(serviceName, prometheusRegistry))
			r.Use(middleware.RequestID)
			r.Use(requestTimeMiddleware)
			r.Use(middlewares...)
		}),
		listenUrl: config.ListenURL,
	}
}

func (r *Router) NewMainRoute(method string, pattern string, handler func(w http.ResponseWriter, r *http.Request)) {
	switch method {
	case http.MethodGet:
		r.mainRouter.Get(pattern, handler)
		break
	case http.MethodPost:
		r.mainRouter.Post(pattern, handler)
		break
	case http.MethodPatch:
		r.mainRouter.Patch(pattern, handler)
		break
	case http.MethodPut:
		r.mainRouter.Put(pattern, handler)
		break
	case http.MethodDelete:
		r.mainRouter.Delete(pattern, handler)
		break
	default:
		panic("method not implemented")
	}
}

func (r *Router) With(middlewares ...func(http.Handler) http.Handler) *Router {
	return &Router{
		Router:     r.Router.With(middlewares...),
		mainRouter: r.mainRouter,
		listenUrl:  r.listenUrl,
	}
}

func (r *Router) NewRoute(method string, version string, pattern string, handler func(w http.ResponseWriter, r *http.Request)) {
	fullPattern := fmt.Sprintf("/%s/%s", version, pattern)
	switch method {
	case http.MethodGet:
		r.Get(fullPattern, handler)
		break
	case http.MethodPost:
		r.Post(fullPattern, handler)
		break
	case http.MethodPatch:
		r.Patch(fullPattern, handler)
		break
	case http.MethodPut:
		r.Put(fullPattern, handler)
		break
	case http.MethodDelete:
		r.Delete(fullPattern, handler)
		break
	default:
		panic("method not implemented")
	}
}

func (r *Router) MountOnMainRouter(pattern string, subRouter chi.Router) {
	r.mainRouter.Mount(pattern, subRouter)
}

func (r *Router) Start() {
	err := http.ListenAndServe(r.listenUrl, r.mainRouter)
	if err != nil {
		panic(err)
	}
}
