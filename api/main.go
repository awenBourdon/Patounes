package main

import (
	"fmt"
	"log"
	"net/http"
	"patoune-api/config/database"
	"patoune-api/database/sqlc"
	"patoune-api/graph"
	"patoune-api/internal/auth"
	user "patoune-api/internal/users"
	authhandlers "patoune-api/server/auth"
	"patoune-api/server/cors"
	ratelimit "patoune-api/server/rate-limit"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const (
	registerRateLimit = 10
	loginRateLimit    = 5
)

func main() {
	fmt.Println("Démarrage...")

	err := database.Connect()
	if err != nil {
		log.Fatalf("Erreur lors de la connexion à la base de données: %v", err)
	}
	defer database.Close()

	if err := ratelimit.InitRateLimiter("localhost:6379"); err != nil {
		log.Fatalf("Erreur Redis: %v", err)
	}
	defer ratelimit.CloseRedis()

	queries := sqlc.New(database.DB)

	userRepo := user.NewUserRepository(queries)

	userService := user.NewUserService(userRepo)
	authService := auth.NewAuthService(queries)

	resolver := &graph.Resolver{
		UserService: userService,
	}

	authHandlers := authhandlers.NewAuthHandlers(authService)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"message": "API Patounes en cours d'exécution"}`)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status": "ok"}`)
	})

	http.HandleFunc("/test-db", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var version string
		err := database.DB.QueryRow("SELECT version()").Scan(&version)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error": "%s"}`, err.Error())
			return
		}

		fmt.Fprintf(w, `{"postgres_version": "%s"}`, version)
	})

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: resolver,
	}))

	http.Handle("/auth/register",
		ratelimit.RateLimitMiddleware(registerRateLimit)(http.HandlerFunc(authHandlers.Register)))

	http.Handle("/auth/login",
		ratelimit.RateLimitMiddleware(loginRateLimit)(http.HandlerFunc(authHandlers.Login)))

	http.HandleFunc("/auth/logout", authHandlers.Logout)
	http.HandleFunc("/auth/me", authHandlers.Me)

	http.Handle("/graphql", srv)
	http.Handle("/sandbox", playground.ApolloSandboxHandler("GraphQL", "/graphql"))

	port := ":3001"
	fmt.Printf("Serveur lancé sur http://localhost%s\n", port)
	fmt.Printf("Test ok: http://localhost%s\n", port)

	log.Fatal(http.ListenAndServe(port, cors.CORSMiddleware(http.DefaultServeMux)))
}
