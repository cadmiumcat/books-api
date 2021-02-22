package initialiser

import (
	dpHttp "github.com/ONSdigital/dp-net/http"
	"github.com/cadmiumcat/books-api/api"
	"github.com/cadmiumcat/books-api/interfaces"
	"github.com/gorilla/mux"
	"net/http"
)

type Service struct {
	Server interfaces.HTTPServer
	Router *mux.Router
	API    *api.API
}

func GetHTTPServer(bindAddr string, router http.Handler) interfaces.HTTPServer {
	httpServer := dpHttp.NewServer(bindAddr, router)
	return httpServer
}
