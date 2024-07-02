package server

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/Eddy150893/blog-sockets/database"
	"github.com/Eddy150893/blog-sockets/repository"
	"github.com/Eddy150893/blog-sockets/websocket"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Config struct {
	Port        string
	JWTSecret   string
	DatabaseUrl string
}

type Server interface {
	Config() *Config
	Hub() *websocket.Hub //Ahora el server tendra un websocket
}

type Broker struct {
	config *Config
	router *mux.Router
	hub    *websocket.Hub //El broker que implementa server necesita una nueva propiedad
}

func (b *Broker) Config() *Config {
	return b.config
}

func (b *Broker) Hub() *websocket.Hub {
	return b.hub
}

func NewServer(ctx context.Context, config *Config) (*Broker, error) {
	if config.Port == "" {
		return nil, errors.New("port is required")
	}

	if config.JWTSecret == "" {
		return nil, errors.New("secret is required")
	}

	if config.DatabaseUrl == "" {
		return nil, errors.New("database url is required")
	}

	broker := &Broker{
		config: config,
		router: mux.NewRouter(),
		hub:    websocket.NewHub(),
	}
	return broker, nil
}

func (b *Broker) Start(binder func(s Server, r *mux.Router)) {
	b.router = mux.NewRouter()
	binder(b, b.router)
	handler := cors.Default().Handler(b.router)
	repo, err := database.NewPostgresRepository(b.config.DatabaseUrl)
	if err != nil {
		log.Fatal(err)
	}
	go b.hub.Run()
	repository.SetRepository(repo)
	log.Println("Starting server on port", b.Config().Port)
	//antes de la implementacion de cors(linea del handler)
	//se colocaba el router
	// if err := http.ListenAndServe(b.config.Port, b.router); err != nil {
	// 	log.Fatal("ListenAndServe: ", err)
	// }
	if err := http.ListenAndServe(b.config.Port, handler); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
