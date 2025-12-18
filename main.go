package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"go-microservice/handlers"
	"go-microservice/metrics"
	"go-microservice/services"
	"go-microservice/utils"
)

func main() {
	userService := services.NewUserService()
	userHandler := handlers.NewUserHandler(userService)

	r := mux.NewRouter()

	r.Use(utils.RateLimitMiddleware)
	r.Use(metrics.MetricsMiddleware)

	r.HandleFunc("/api/users", userHandler.GetAllUsers).Methods("GET")
	r.HandleFunc("/api/users/{id}", userHandler.GetUser).Methods("GET")
	r.HandleFunc("/api/users", userHandler.CreateUser).Methods("POST")
	r.HandleFunc("/api/users/{id}", userHandler.UpdateUser).Methods("PUT")
	r.HandleFunc("/api/users/{id}", userHandler.DeleteUser).Methods("DELETE")

	r.Handle("/metrics", metrics.Handler())

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	log.Println("Server starting on :8081")
	if err := http.ListenAndServe(":8081", r); err != nil {
		log.Fatal(err)
	}
}
