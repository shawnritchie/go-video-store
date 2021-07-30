package http

import (
	"github.com/gorilla/mux"
	"github.com/shawnritchie/go-video-store/internal/port/driven"
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

//Step 1. Only single Method per interface definition
//func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Content-Type", "application/json")
//	w.WriteHeader(http.StatusOK)
//	w.Write([]byte(`{"message": "hello world"}`))
//}
//
//func main() {
//	s := &server{}
//	http.Handle("/", s)    //
//	log.Fatal(http.ListenAndServe(":8080", nil))
//}

/*
func home(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    switch r.Method {
    case "GET":
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"message": "get called"}`))
    case "POST":
        w.WriteHeader(http.StatusCreated)
        w.Write([]byte(`{"message": "post called"}`))
    case "PUT":
        w.WriteHeader(http.StatusAccepted)
        w.Write([]byte(`{"message": "put called"}`))
    case "DELETE":
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"message": "delete called"}`))
    default:
        w.WriteHeader(http.StatusNotFound)
        w.Write([]byte(`{"message": "not found"}`))
    }
}

func main() {
    http.HandleFunc("/", home)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
*/
