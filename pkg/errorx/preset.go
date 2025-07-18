package errorx

import "net/http"

func BadRequest(message string) Error {
	return New(CodeBadRequest.Int(), http.StatusBadRequest, message)
}

func Unauthorized(message string) Error {
	return New(CodeUnauthorized.Int(), http.StatusUnauthorized, message)
}

func Internal(message string) Error {
	return New(CodeInternal.Int(), http.StatusInternalServerError, message)
}
