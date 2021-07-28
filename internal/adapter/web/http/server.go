package http

import (
	"github.com/gorilla/mux"
	"github.com/shawnritchie/go-video-store/internal/port/driven"
	"net/http"
	"sync"
)

type server struct {
	finder   driven.FilmFinder
	appender driven.FilmAppender
	invoicer driven.FilmInvoicer
	once     sync.Once
	router   *mux.Router
}

func New(finder driven.FilmFinder, appender driven.FilmAppender, invoicer driven.FilmInvoicer) *server {
	return &server{
		finder:   finder,
		appender: appender,
		invoicer: invoicer,
	}
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "hello world"}`))
}
