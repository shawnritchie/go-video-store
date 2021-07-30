package http

import (
	"github.com/gorilla/mux"
	"net/http"
)

/*
curl -X GET http://localhost:8080/catalogue/film?name=Loki -H "Content-Type: application/json"

curl -X POST http://localhost:8080/catalogue/film/new -H "Content-Type: application/json" -d '{"name":"Loki", "director":"Marvel"}'
curl -X POST http://localhost:8080/catalogue/film/regular -H "Content-Type: application/json" -d '{"name":"Black Widow", "director":"Marvel"}'
curl -X POST http://localhost:8080/catalogue/film/old -H "Content-Type: application/json" -d '{"name":"Morbius", "director":"Marvel"}'

curl -X POST http://localhost:8080/store/return -H "Content-Type: application/json" -d '{"return":[{"name": "Loki", "days": 1}]}'
*/

func (s *server) Router() (r *mux.Router) {
	s.once.Do(func() {
		r = mux.NewRouter()
		//r.HandleFunc("/catalogue/film/new", s.addNewFilm).Methods(http.MethodPost)
		//r.Handle("/catalogue/film/regular", handler(s.addRegularFilm)).Methods(http.MethodPost)
		//r.Handle("/catalogue/film/old", handler(s.addOldFilm)).Methods(http.MethodPost)
		r.Handle("/catalogue/film/{release}", handler(s.addFilm)).Methods(http.MethodPost)

		r.Handle("/catalogue/film", handler(s.findFilm)).Methods(http.MethodGet)
		r.Handle("/store/return", handler(s.processReturn)).Methods(http.MethodPost)
		s.router = r
	})
	return s.router
}
