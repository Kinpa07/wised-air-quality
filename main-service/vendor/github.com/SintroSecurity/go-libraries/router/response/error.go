package response

import (
	"fmt"
	"net/http"
)

type ErrorCode string
type ErrorMessage string

type Error struct {
	Code    ErrorCode    `json:"code"`
	Message ErrorMessage `json:"message,omitempty"`
	Details interface{}  `json:"details,omitempty"`
}

var (
	ErrorCodeInternal              ErrorCode = "internalError"
	ErrorCodeBadRequest            ErrorCode = "badRequest"
	ErrorCodeUnauthorized          ErrorCode = "unauthorized"
	ErrorCodeForbidden             ErrorCode = "forbidden"
	ErrorCodePaymentRequired       ErrorCode = "paymentRequired"
	ErrorCodeNotFound              ErrorCode = "notFound"
	ErrorCodeConflict              ErrorCode = "conflict"
	ErrorCodeUnprocessableEntity   ErrorCode = "unprocessableEntity"
	ErrorCodeRequestEntityTooLarge ErrorCode = "requestEntityTooLarge"
	ErrorCodeResourceLimitReached  ErrorCode = "resourceLimitReached"
)

var (
	ErrorMessageInternalServerError   ErrorMessage = "Internal server error"
	ErrorMessageBadRequest            ErrorMessage = "Bad request"
	ErrorMessageUnauthorized          ErrorMessage = "Unauthorized"
	ErrorMessageForbidden             ErrorMessage = "Forbidden"
	ErrorMessagePaymentRequired       ErrorMessage = "Payment Required"
	ErrorMessageNotFound              ErrorMessage = "Not found"
	ErrorMessageConflict              ErrorMessage = "Conflict"
	ErrorMessageUnprocessableEntity   ErrorMessage = "Unprocessable entity"
	ErrorMessageRequestEntityTooLarge ErrorMessage = "Request entity too large"
	ErrorMessageResourceLimitReached  ErrorMessage = "Resource limit reached"
)

func NewResponseError(code ErrorCode, message ErrorMessage) *Error {
	return &Error{Code: code, Message: message}
}

func NewResponseErrorWithDetails(code ErrorCode, message ErrorMessage, details interface{}) *Error {
	return &Error{Code: code, Message: message, Details: details}
}

func (e *Error) Error() string {
	s := fmt.Sprintf("%s: %s", e.Code, e.Message)
	return s
}

func RespondError(r *http.Request, w http.ResponseWriter, e *Error) {
	if e != nil {
		switch e.Code {
		case ErrorCodeBadRequest:
			BadRequestJson(r, w, e)
		case ErrorCodeUnauthorized:
			UnauthorizedJson(r, w, e)
		case ErrorCodeForbidden:
			ForbiddenJson(r, w, e)
		case ErrorCodePaymentRequired:
			PaymentRequiredJson(r, w, e)
		case ErrorCodeNotFound:
			NotFoundJson(r, w, e)
		case ErrorCodeConflict:
			ConflictJson(r, w, e)
		case ErrorCodeUnprocessableEntity:
			UnprocessableEntityJson(r, w, e)
		case ErrorCodeRequestEntityTooLarge:
			RequestEntityTooLargeJson(r, w, e)
		case ErrorCodeResourceLimitReached:
			ResourceLimitReachedJson(r, w, e)
		default:
			InternalServerErrorJson(r, w, e)
		}
	}
}
