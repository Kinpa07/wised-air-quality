package sensor_readings_collector_pkg

import (
	"time"

	"github.com/SintroSecurity/go-libraries/router"
	"github.com/go-playground/validator/v10"
)

// clockSkewAllowance caps how far ahead of now a timestamp may be: a far-future
// measured_at would pin itself as the permanent latest, since the cache only
// advances to newer readings and nothing real ever beats it.
const clockSkewAllowance = 5 * time.Minute

func notFuture(fl validator.FieldLevel) bool {
	t, ok := fl.Field().Interface().(time.Time)
	if !ok {
		return false
	}
	return !t.After(time.Now().Add(clockSkewAllowance))
}

func init() {
	router.AddCustomValidator("not_future", notFuture)
}
