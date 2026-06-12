package server

import (
	"go-service-skeleton/internal/app/controller"
	"net/http"

	"github.com/SintroSecurity/go-libraries/router/response"
	"github.com/ggicci/httpin"

	sensor_readings_collector_pkg "go-service-skeleton/pkg/sensor-readings-collector"
)

func CreateReadingHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		input := r.Context().Value(httpin.Input).(*sensor_readings_collector_pkg.CreateReadingRequest)
		result, err := controller.CreateReading(r.Context(), input)
		if err != nil {
			response.RespondError(r, w, err)
			return
		}

		if result.Created {
			response.CreatedJson(r, w, result)
			return
		}

		response.OkJson(r, w, result)
	}
}
