package server

import (
	"go-service-skeleton/internal/app/controller"
	"net/http"

	"github.com/SintroSecurity/go-libraries/router/response"
)

func GetStationsHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := controller.GetStations(r.Context())
		if err != nil {
			response.RespondError(r, w, err)
			return
		}

		response.OkJson(r, w, result)
	}
}
