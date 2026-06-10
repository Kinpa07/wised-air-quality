package router

import (
	"context"
	"net/http"
	"reflect"
	"strings"

	"github.com/SintroSecurity/go-libraries/logger"
	"github.com/SintroSecurity/go-libraries/router/response"
	"github.com/ggicci/httpin"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func NewValidator(ctx context.Context) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			// Here we read the incoming struct, then we proceed to test it against validator
			// If error(s) has been found, they'll be returned with 400 HTTP status code
			l := logger.GetLoggerFromContext(ctx)

			// Fetch inputStruct from httpin middleware
			inputStruct := r.Context().Value(httpin.Input)

			l.Debug("unmarshaled input data", l.Any("struct", inputStruct))

			err := validate.Struct(inputStruct)
			if err != nil {
				l := logger.GetLoggerFromContext(ctx)
				// this check is only needed when your code could produce
				// an invalid value for validation such as interface with nil
				// value most including myself do not usually have code like this.
				if _, ok := err.(*validator.InvalidValidationError); ok {
					//Log and return empty 400 bad request error
					l.Info("got error validating input struct", l.Any("type", reflect.TypeOf(inputStruct)), l.Err(err))
					response.BadRequestJson(r, rw, &response.Error{
						Code:    response.ErrorCodeBadRequest,
						Message: response.ErrorMessageBadRequest,
						Details: nil,
					})
					return
				}

				type Details struct {
					Field string      `json:"field"`
					Value interface{} `json:"value"`
				}

				details := make([]Details, 0, len(err.(validator.ValidationErrors)))

				inputType := reflect.TypeOf(inputStruct)
				if inputType.Kind() == reflect.Ptr {
					inputType = inputType.Elem()
				}

				for _, fieldErr := range err.(validator.ValidationErrors) {
					jsonField := getJSONFieldName(inputType, fieldErr.StructNamespace())
					var errorDetails Details

					errorDetails.Field = jsonField
					errorDetails.Value = fieldErr.Value()

					details = append(details, errorDetails)
				}

				if len(details) > 0 {
					response.BadRequestJson(r, rw, &response.Error{
						Code:    response.ErrorCodeBadRequest,
						Message: response.ErrorMessageBadRequest,
						Details: details,
					})
					return
				}

				return
			}

			next.ServeHTTP(rw, r)
		})
	}
}

func AddCustomValidator(tag string, validator func(sl validator.FieldLevel) bool) error {
	return validate.RegisterValidation(tag, validator, true)
}

// Extracts the JSON field name from a nested struct path
func getJSONFieldName(t reflect.Type, structNamespace string) string {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	parts := strings.Split(structNamespace, ".")

	// Remove the root struct name (e.g. "RegisterUserRequest")
	if len(parts) > 1 {
		parts = parts[1:]
	} else {
		return ""
	}

	for i, part := range parts {
		if t.Kind() != reflect.Struct {
			return ""
		}

		field, ok := t.FieldByName(part)
		if !ok {
			return ""
		}

		if i == len(parts)-1 {
			// Last part, return JSON field name
			jsonTag := field.Tag.Get("json")
			jsonName := strings.Split(jsonTag, ",")[0]
			if jsonName == "" || jsonName == "-" {
				return field.Name
			}
			return jsonName
		}

		// Go deeper
		t = field.Type
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
	}

	return ""
}
