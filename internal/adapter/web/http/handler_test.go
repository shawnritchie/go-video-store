package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_MissingContentType(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "catalogue/film", nil)
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	handler(func(writer http.ResponseWriter, request *http.Request) error {
		return nil
	}).ServeHTTP(res, req)

	var errorResponse Error
	unmarshalBody(t, res, &errorResponse)

	if res.Code != http.StatusBadRequest {
		t.Errorf("got status %d but wanted %d", res.Code, http.StatusBadRequest)
	}

	if errorResponse.Detail == "" {
		t.Errorf("received unexpected response %#v", errorResponse)
	}
}

func TestHandler_UnhandledError(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "catalogue/film", nil)
	req.Header.Add("Content-Type", contentType)
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	handler(func(writer http.ResponseWriter, request *http.Request) error {
		return fmt.Errorf("unknown error")
	}).ServeHTTP(res, req)

	var errorResponse Error
	unmarshalBody(t, res, &errorResponse)

	if res.Code != http.StatusInternalServerError {
		t.Errorf("got status %d but wanted %d", res.Code, http.StatusInternalServerError)
	}
}

func unmarshalBody(t *testing.T, w *httptest.ResponseRecorder, res interface{}) {
	reqBody, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Errorf("request body cannot be read : %w", err)
	}

	if err := json.Unmarshal(reqBody, &res); err != nil {
		t.Errorf("Post response cannot be deserialized. %w", err)
	}
}

//
//func TestHandler_UnserializableResponse(t *testing.T) {
//	req, err := http.NewRequest(http.MethodGet,"catalogue/film", nil)
//	req.Header.Add("Content-Type", contentType)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	res := httptest.NewRecorder()
//	handler(func(writer http.ResponseWriter, request *http.Request) error {
//		return NewClientError(nil, http.StatusBadRequest, *(*string)(unsafe.Pointer(nil)))
//	}).ServeHTTP(res, req)
//
//	var errorResponse Error
//	unmarshalBody(res, &errorResponse)
//
//	if res.Code != http.StatusInternalServerError {
//		t.Errorf("got status %d but wanted %d", res.Code, http.StatusInternalServerError)
//	}
//}

func toJSON(request interface{}) io.Reader {
	stuff, _ := json.Marshal(request)
	return bytes.NewReader(stuff)
}
