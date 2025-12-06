package action

import "net/http"

// Response .
func Response(err error, bytes []byte, statusCode int) (error, []byte, int) {
	return err, bytes, statusCode
}

// ResponseInternalServerError .
func ResponseInternalServerError(err error, bytes []byte) (error, []byte, int) {
	return err, bytes, http.StatusInternalServerError
}

// ResponseOK .
func ResponseOK(bytes []byte) (error, []byte, int) {
	return nil, bytes, http.StatusOK
}

// ResponseNotFound .
func ResponseNotFound(bytes []byte) (error, []byte, int) {
	return nil, bytes, http.StatusNotFound
}

// ResponseForbidden .
func ResponseForbidden(bytes []byte) (error, []byte, int) {
	return nil, bytes, http.StatusForbidden
}

// ResponseBadRequest .
func ResponseBadRequest(bytes []byte) (error, []byte, int) {
	return nil, bytes, http.StatusBadRequest
}
