package main

import (
	"github.com/shawnritchie/go-video-store/internal/adapter/repository/inmem"
	web "github.com/shawnritchie/go-video-store/internal/adapter/web/http"
	"github.com/shawnritchie/go-video-store/internal/service"
	"log"
	"net/http"
)

func main() {
	catalogue := &inmem.StoreCatalogue{}
	service := service.New(catalogue, catalogue)
	s := web.New(
		service,
		service,
		service,
	)
	log.Fatal(http.ListenAndServe(":8080", s.Router()))
}
