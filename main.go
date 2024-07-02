package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/Eddy150893/blog-sockets/handlers"
	"github.com/Eddy150893/blog-sockets/middleware"
	"github.com/Eddy150893/blog-sockets/server"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	PORT := os.Getenv("PORT")
	JWT_SECRET := os.Getenv("JWT_SECRET")
	DATABASE_URL := os.Getenv("DATABASE_URL")

	s, err := server.NewServer(context.Background(), &server.Config{
		Port:        PORT,
		JWTSecret:   JWT_SECRET,
		DatabaseUrl: DATABASE_URL,
	})

	if err != nil {
		log.Fatal(err)
	}

	s.Start(BindRoutes)
}

func BindRoutes(s server.Server, r *mux.Router) {

	//Este subrouter servira para aplicar el middleware
	//Esto es util por que si quisieramos aplicar
	//diferentes middlewares a diferentes rutas
	//podriamos crear sub routers
	api := r.PathPrefix("/api/v1").Subrouter()

	api.Use(middleware.CheckAuthMiddleware(s))
	r.HandleFunc("/", handlers.HomeHandler(s)).Methods(http.MethodGet)
	r.HandleFunc("/signup", handlers.SignUpHandler(s)).Methods(http.MethodPost)
	r.HandleFunc("/login", handlers.LoginHandler(s)).Methods(http.MethodPost)
	api.HandleFunc("/me", handlers.MeHandler(s)).Methods(http.MethodGet)
	api.HandleFunc("/posts", handlers.InsertPostHandler(s)).Methods(http.MethodPost)
	r.HandleFunc("/posts/{id}", handlers.GetPostByIdHandler(s)).Methods(http.MethodGet)
	api.HandleFunc("/posts/{id}", handlers.UpdatePostHandler(s)).Methods(http.MethodPut)
	api.HandleFunc("/posts/{id}", handlers.DeletePostHandler(s)).Methods(http.MethodDelete)
	r.HandleFunc("/posts", handlers.ListPostHandler(s)).Methods(http.MethodGet)
	r.HandleFunc("/ws", s.Hub().HandleWebSocket)
}
