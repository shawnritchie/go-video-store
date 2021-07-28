package http

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//Any error to be returned to the client
type ClientError interface {
	error
	ResponseBody() ([]byte, error)
	ResponseHeaders() (int, map[string]string)
}

type Error struct {
	Cause  error  `json:"-"`
	Detail string `json:"detail"`
	Status int    `json:"-"`
}

var (
	contentType = "application/json"
	httpHeader  = map[string]string{
		"Content-Type": contentType,
	}
	TypeClientError *Error
)

func (e *Error) Error() string {
	if e.Cause == nil {
		return e.Detail
	}
	return e.Detail + " : " + e.Cause.Error()
}

func (e *Error) ResponseBody() ([]byte, error) {
	body, err := json.Marshal(e)
	if err != nil {
		return nil, fmt.Errorf("error while parsing response body: %v", err)
	}
	return body, nil
}

func (e *Error) ResponseHeaders() (int, map[string]string) {
	return e.Status, httpHeader
}

func NewClientError(err error, status int, detail string) error {
	return &Error{
		Cause:  err,
		Detail: detail,
		Status: status,
	}
}

type handler func(http.ResponseWriter, *http.Request) error

func (fn handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	ct := r.Header.Get("Content-Type")
	if ct != contentType {
		err = NewClientError(nil,
			http.StatusBadRequest,
			fmt.Sprintf("Bad Request: Content-Type must be set to %q", contentType),
		)
	} else {
		err = fn(w, r)
	}

	if err == nil {
		return
	}

	clientError, ok := err.(ClientError)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, err := clientError.ResponseBody()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	status, headers := clientError.ResponseHeaders()
	for k, v := range headers {
		w.Header().Set(k, v)
	}
	w.WriteHeader(status)
	w.Write(body)
}

func setHeaders(w http.ResponseWriter) {
	for k, v := range httpHeader {
		w.Header().Set(k, v)
	}
}
