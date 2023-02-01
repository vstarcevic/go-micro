package main

import (
	bm "broker/middleware"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	// specify who is allowed to connect
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Use(bm.RateLimiterMiddleware(app.Redis))

	mux.Use(middleware.Heartbeat("/ping"))

	mux.Post("/", app.Broker)

	mux.Post("/log-grpc", app.LogViaGRPC)

	mux.Post("/handle", app.HandleSubmission)

	return mux
}

// func CustomMiddleware() func(http.Handler) http.Handler {
// 	f := func(h http.Handler) http.Handler {
// 		fn := func(w http.ResponseWriter, r *http.Request) {

// 			cookie := http.Cookie{Name: "mycookie", Value: "myvalue"}
// 			if r.Method == "GET" || r.Method == "POST" {
// 				fmt.Println("hey you hit my custom middleware")
// 				w.Header().Add("HeyHeyHey", "HeyHeyHey2") // NOTE THIS LINE

// 				w.Write([]byte("end of response"))

// 				return
// 			}
// 			h.ServeHTTP(w, r)
// 		}
// 		return http.HandlerFunc(fn)
// 	}
// 	return f
// }
