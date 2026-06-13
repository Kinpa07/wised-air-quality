package server

import (
	"context"
	"errors"
	sensor_readings_collector_pkg "go-service-skeleton/pkg/sensor-readings-collector"
	"net/http"
	"strings"

	"github.com/SintroSecurity/go-libraries/router"
	"github.com/SintroSecurity/go-libraries/router/response"
	"github.com/ggicci/httpin"
	httpin_core "github.com/ggicci/httpin/core"
	httpin_integration "github.com/ggicci/httpin/integration"
	"github.com/go-chi/chi/v5"
)

func init() {
	httpin_integration.UseGochiURLParam("path", chi.URLParam)

	httpin_core.RegisterErrorHandler(func(w http.ResponseWriter, r *http.Request, err error) {
		var invalidFieldError *httpin_core.InvalidFieldError
		if errors.As(err, &invalidFieldError) {
			response.UnprocessableEntityJson(r, w, map[string]interface{}{
				"error": "field not provided or invalid",
				"field": strings.ToLower(invalidFieldError.Field),
			})
			return
		} else {
			response.InternalServerErrorJson(r, w, "internal server error")
		}
	})
}

func InitRoutingTable(ctx context.Context, r *router.Router) {
	r.With(httpin.NewInput(sensor_readings_collector_pkg.GetClientsRequest{})).NewRoute(http.MethodGet, "v1", "clients", GetClientsHandler())
	r.With(httpin.NewInput(sensor_readings_collector_pkg.CreateClientRequest{}), router.NewValidator(ctx)).NewRoute(http.MethodPost, "v1", "clients", CreateClientHandler())
	r.With(httpin.NewInput(sensor_readings_collector_pkg.CreateReadingRequest{}), router.NewValidator(ctx)).NewRoute(http.MethodPost, "v1", "clients/{client_id}:read", CreateReadingHandler())
	r.With(httpin.NewInput(sensor_readings_collector_pkg.GetReadingsRequest{})).NewRoute(http.MethodGet, "v1", "clients/{client_id}/readings", GetReadingHandler())
	r.NewRoute(http.MethodGet, "v1", "stations", GetStationsHandler())
	r.NewRoute(http.MethodGet, "v1", "stats", GetStatsHandler())
}
