package server

import (
	"go-service-skeleton/internal/app/controller"
	"net/http"

	"github.com/SintroSecurity/go-libraries/router/response"
	"github.com/ggicci/httpin"

	sensor_readings_collector_pkg "go-service-skeleton/pkg/sensor-readings-collector"
)

func GetClientsHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		input := r.Context().Value(httpin.Input).(*sensor_readings_collector_pkg.GetClientsRequest)
		result, err := controller.GetClients(r.Context(), input)
		if err != nil {
			response.RespondError(r, w, err)
			return
		}
		response.OkJson(r, w, result)
	}
}

func CreateClientHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		input := r.Context().Value(httpin.Input).(*sensor_readings_collector_pkg.CreateClientRequest)
		result, err := controller.CreateClient(r.Context(), input)
		if err != nil {
			response.RespondError(r, w, err)
			return
		}
		response.CreatedJson(r, w, result)
	}
}
