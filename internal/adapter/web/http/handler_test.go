package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
	unmarshalBody(res, &errorResponse)

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
	unmarshalBody(res, &errorResponse)

	if res.Code != http.StatusInternalServerError {
		t.Errorf("got status %d but wanted %d", res.Code, http.StatusInternalServerError)
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
