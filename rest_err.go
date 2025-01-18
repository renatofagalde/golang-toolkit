package toolkit

import (
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/renatofagalde/golang-toolkit/context_manager"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

var logger Logger

type RestErr struct {
	Message string  `json:"message"`
	Err     string  `json:"error"`
	Code    int     `json:"code"`
	Causes  []Cause `json:"causes"`
}

type Cause struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type Causes []Cause

func NewCauses(causes ...Cause) Causes {
	return Causes(causes)
}

// construtor
func (t *RestErr) NewRestErr(message, err string, code int, causes []Cause) *RestErr {
	return &RestErr{
		Message: message,
		Err:     err,
		Code:    code,
		Causes:  causes,
	}
}

func (r *RestErr) Error() string {
	return r.Message
}

func (t *RestErr) NewBadRequestError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "bad_request",
		Code:    http.StatusBadRequest,
	}
}

func (t *RestErr) NewBadRequestValidationError(message string, causes []Cause) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "bad_request",
		Code:    http.StatusBadRequest,
		Causes:  causes,
	}
}

func (t *RestErr) NewInternalServerError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "internal_server_error",
		Code:    http.StatusInternalServerError,
	}
}

func (t *RestErr) NewForbiddenError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "forbidden",
		Code:    http.StatusForbidden,
	}
}

func (t *RestErr) NewNotFoundError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "notFound",
		Code:    http.StatusNotFound,
	}
}

func (t *RestErr) NewUnauthorizedRequestError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "unauthorized",
		Code:    http.StatusUnauthorized,
	}
}

func (t *RestErr) NewConflictError(message string, causes []Cause) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "conflict",
		Code:    http.StatusConflict,
		Causes:  causes,
	}
}

func (t *RestErr) NewSystemResourceError(message string, causes []Cause) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "system_resource_error",
		Code:    http.StatusInternalServerError, // 500 Internal Server Error
		Causes:  causes,
	}
}

func (t *RestErr) NewTransactionError(message string, causes []Cause) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "transaction_error",
		Code:    http.StatusConflict, // 409 Conflict
		Causes:  causes,
	}
}

func (t *RestErr) NewSQLSyntaxError(message string, causes []Cause) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "sql_syntax_error",
		Code:    http.StatusBadRequest, // 400 Bad Request
		Causes:  causes,
	}
}

func (t *RestErr) NewConcurrencyError(message string, causes []Cause) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "concurrency_error",
		Code:    http.StatusConflict, // 409 Conflict
		Causes:  causes,
	}
}

func (t *RestErr) NewStorageSpaceError(message string, causes []Cause) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "storage_space_error",
		Code:    http.StatusInsufficientStorage, // 507 Insufficient Storage
		Causes:  causes,
	}
}

func (t *RestErr) NewDataIntegrityError(message string, causes []Cause) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "data_integrity_error",
		Code:    http.StatusConflict, // 409 Conflict
		Causes:  causes,
	}
}

func (t *RestErr) HandlePgError(pgErr *pgconn.PgError) *RestErr {

	journey, requestID := context_manager.Give()

	logger.Warn("Database error",
		zap.String("stage", "repository"),
		zap.String("journey", journey),
		zap.String("requestID", requestID),
		zap.String("pg_error_code", pgErr.Code),
		zap.String("pg_error_message", pgErr.Message),
		zap.String("pg_error_detail", pgErr.Detail),
		zap.String("pg_error_where", pgErr.Where),
	)

	cause := Cause{
		Field:   pgErr.Code,
		Message: pgErr.Message,
	}
	causes := []Cause{cause}

	switch pgErr.Code {
	case "23505": // Unique violation
		return (&RestErr{}).NewConflictError("Duplicate key error while inserting user", causes)
	case "40001": // Serialization failure (deadlock)
		return (&RestErr{}).NewTransactionError("Transaction deadlock detected", causes)
	case "22001": // String data, right truncation
		return (&RestErr{}).NewSQLSyntaxError("String data is too long", causes)
	default:
		// Generic internal server error for other codes
		errorCode, _ := strconv.Atoi(pgErr.Code)
		return (&RestErr{}).NewRestErr(
			"Unexpected database error occurred",
			pgErr.Message,
			errorCode,
			causes,
		)
	}
}
