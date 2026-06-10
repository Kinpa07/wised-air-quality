package response

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/go-chi/render"
)

func Ok(r *http.Request, w http.ResponseWriter) {
	render.Status(r, http.StatusOK)
}

func OkJson(r *http.Request, w http.ResponseWriter, data interface{}) {
	render.Status(r, http.StatusOK)
	JSON(w, r, data, true)
}

func OkJsonWithoutEscape(r *http.Request, w http.ResponseWriter, data interface{}) {
	render.Status(r, http.StatusOK)
	JSON(w, r, data, false)
}

func Created(r *http.Request, w http.ResponseWriter) {
	render.Status(r, http.StatusCreated)
}

func CreatedJson(r *http.Request, w http.ResponseWriter, data interface{}) {
	render.Status(r, http.StatusCreated)
	JSON(w, r, data, true)
}

func Accepted(r *http.Request, w http.ResponseWriter) {
	render.Status(r, http.StatusAccepted)
}

func AcceptedJson(r *http.Request, w http.ResponseWriter, data interface{}) {
	render.Status(r, http.StatusAccepted)
	JSON(w, r, data, true)
}

func NoContent(r *http.Request, w http.ResponseWriter) {
	render.Status(r, http.StatusNoContent)
}

func NoContentJson(r *http.Request, w http.ResponseWriter, data interface{}) {
	render.Status(r, http.StatusNoContent)
	JSON(w, r, data, true)
}

func BadRequest(r *http.Request, w http.ResponseWriter) {
	render.Status(r, http.StatusBadRequest)
}

func BadRequestJson(r *http.Request, w http.ResponseWriter, data interface{}) {
	render.Status(r, http.StatusBadRequest)
	JSON(w, r, data, true)
}

func Unauthorized(r *http.Request, w http.ResponseWriter) {
	render.Status(r, http.StatusUnauthorized)
}

func UnauthorizedJson(r *http.Request, w http.ResponseWriter, data interface{}) {
	render.Status(r, http.StatusUnauthorized)
	JSON(w, r, data, true)
}

func Forbidden(r *http.Request, w http.ResponseWriter) {
	render.Status(r, http.StatusForbidden)
}

func ForbiddenJson(r *http.Request, w http.ResponseWriter, data interface{}) {
	render.Status(r, http.StatusForbidden)
	JSON(w, r, data, true)
}

func PaymentRequired(r *http.Request, w http.ResponseWriter, data interface{}) {
	render.Status(r, http.StatusPaymentRequired)
}

func PaymentRequiredJson(r *http.Request, w http.ResponseWriter, data interface{}) {
	render.Status(r, http.StatusPaymentRequired)
	JSON(w, r, data, true)
}

func NotFound(r *http.Request, w http.ResponseWriter) {
	render.Status(r, http.StatusNotFound)
}

func NotFoundJson(r *http.Request, w http.ResponseWriter, data interface{}) {
	render.Status(r, http.StatusNotFound)
	JSON(w, r, data, true)
}

func Conflict(r *http.Request, w http.ResponseWriter) {
	render.Status(r, http.StatusConflict)
}

func ConflictJson(r *http.Request, w http.ResponseWriter, data interface{}) {
	render.Status(r, http.StatusConflict)
	JSON(w, r, data, true)
}

func UnprocessableEntity(r *http.Request, w http.ResponseWriter) {
	render.Status(r, http.StatusUnprocessableEntity)
}

func UnprocessableEntityJson(r *http.Request, w http.ResponseWriter, data interface{}) {
	render.Status(r, http.StatusUnprocessableEntity)
	JSON(w, r, data, true)
}

func RequestEntityTooLarge(r *http.Request, w http.ResponseWriter) {
	render.Status(r, http.StatusRequestEntityTooLarge)
}

func RequestEntityTooLargeJson(r *http.Request, w http.ResponseWriter, data interface{}) {
	render.Status(r, http.StatusRequestEntityTooLarge)
	JSON(w, r, data, true)
}

func ResourceLimitReached(r *http.Request, w http.ResponseWriter) {
	render.Status(r, http.StatusLoopDetected) //Resource Limit Reached when over HTTP
}

func ResourceLimitReachedJson(r *http.Request, w http.ResponseWriter, data interface{}) {
	render.Status(r, http.StatusLoopDetected) //Resource Limit Reached when over HTTP
	JSON(w, r, data, true)
}

func InternalServerError(r *http.Request, w http.ResponseWriter) {
	render.Status(r, http.StatusInternalServerError)
}

func InternalServerErrorJson(r *http.Request, w http.ResponseWriter, data interface{}) {
	render.Status(r, http.StatusInternalServerError)
	JSON(w, r, data, true)
}

// JSON encodes a struct to JSON,
// automatically converting any field implementing Marshaller.
func JSON(w http.ResponseWriter, r *http.Request, v interface{}, escapeHTML bool) {
	buf := &bytes.Buffer{}

	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(escapeHTML)
	if err := enc.Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if status, ok := r.Context().Value(render.StatusCtxKey).(int); ok {
		w.WriteHeader(status)
	}
	w.Write(buf.Bytes())
}
