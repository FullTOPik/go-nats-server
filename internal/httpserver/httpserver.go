package httpserver

import (
	"context"
	"github.com/go-chi/chi/v5"
	"go_nats-streaming_pg/internal/db"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type Server struct {
	router *chi.Mux
	cache  *db.Cache
	server *http.Server
	data   *db.OrderDto
}

func NewServer(cache *db.Cache) *Server {
	server := Server{}
	server.Init(cache)
	return &server
}

func (s *Server) Init(cache *db.Cache) {
	s.cache = cache
	s.data = &db.OrderDto{}
	s.router = chi.NewRouter()

	s.router.Get("/", s.StartHandler)
	s.router.Route("/{orderId}", func(r chi.Router) {
		r.Get("/", s.GetOrder)
	})

	s.StartServer()
}

func (s *Server) StartServer() {
	s.server = &http.Server{
		Addr:    ":4001",
		Handler: s.router,
	}

	go func() {
		if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("ListenAndServe error %+v", err)
			return
		}
	}()
}

func (s *Server) FinishServer() {
	if err := s.server.Shutdown(context.Background()); err != nil {
		panic(err)
	}

	log.Println("Server successfully shutdown")
}

func (s *Server) StartHandler(w http.ResponseWriter, r *http.Request) {
	templ, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Printf("Error to parse template, %+v", err)
		http.Error(w, "Server error", 500)
	}

	w.WriteHeader(200)
	err = templ.ExecuteTemplate(w, "index.html", s.data)
	if err != nil {
		log.Printf("Error to execute template, %+v", err)
		return
	}
}

func (s *Server) GetOrder(w http.ResponseWriter, r *http.Request) {
	orderIDString := chi.URLParam(r, "orderId")
	orderID, err := strconv.ParseInt(orderIDString, 10, 64)
	if err != nil {
		log.Printf("Error to convert orderID, %+v", err)
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}

	log.Printf("Find order (%d)", orderID)
	s.data, err = s.cache.GetOrderInCache(orderID)

	if err != nil {
		log.Printf("Error to find order, %+v", err)
		http.Error(w, http.StatusText(404), http.StatusNotFound)
		return
	}

	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Printf("Error to parse html file, %+v", err)
		http.Error(w, "Server error", 500)
		return
	}

	w.WriteHeader(200)
	t.ExecuteTemplate(w, "index.html", s.data)
	if err != nil {
		log.Printf("Template error, %+v", err)
	}
}
