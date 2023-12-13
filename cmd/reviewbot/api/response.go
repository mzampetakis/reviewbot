package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"reviewbot/app"
)

func JSON(w http.ResponseWriter, status int, data any) error {
	return JSONWithHeaders(w, status, data, nil)
}

func JSONWithHeaders(w http.ResponseWriter, status int, data any, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

// Ok encodes to JSON and writes the provided response (if any) along with the httpStatus.
func Ok(w http.ResponseWriter, response interface{}, httpStatus int) {
	buf := new(bytes.Buffer)
	if response != nil {
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(response)
		if err != nil {

			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	if buf.Len() > 0 {
		// Remove the extra newline json.Encoder.Encode() adds.
		w.Write(bytes.TrimRight(buf.Bytes(), "\n"))
	}
}

// BadRequestError writes the provided error along with a 400 http status.
func BadRequestError(w http.ResponseWriter, err error) {
	Error(w, err, http.StatusBadRequest)
}

// NotFoundError writes the provided error along with a 404 http status.
func NotFoundError(w http.ResponseWriter, err error) {
	Error(w, err, http.StatusNotFound)
}

// ServerError writes the provided error along with a 500 http status.
func ServerError(w http.ResponseWriter, err error) {
	Error(w, err, http.StatusInternalServerError)
}

// Error writes the provided error along with the provided http status.
func Error(w http.ResponseWriter, err error, httpStatus int) {
	// Malformed request error.
	if e, ok := err.(*malformedRequestError); ok {
		err = &app.Error{
			Msg: e.msg,
			Err: e.err,
		}
		httpStatus = e.status
	}

	// 5xx
	if httpStatus >= http.StatusInternalServerError {
		http.Error(w, http.StatusText(httpStatus), httpStatus)
		return
	}

	apiError := ErrorWrapper{
		Error: ApiError{
			Message: app.ErrorMessage(err),
		},
	}

	res, e := json.Marshal(apiError)
	if e != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	w.Write(res)
}
