package main

import (
	"entrytask/internal/http"
	ziherpc "entrytask/internal/rpc"
	"log"
	"net/http"
)

func main() {
	client, err := ziherpc.Dial("tcp", httpserver.ServerUrl)
	if err != nil {
		log.Fatal("Dial error:", err)
	}
	httpserver.SetClient(client)
	server := http.NewServeMux()
	// Create a new HTTP server
	server.HandleFunc("/api/login", httpserver.LoginHandler)
	server.HandleFunc("/api/profile", httpserver.AuthMiddleware(httpserver.ProfileHandler))
	server.HandleFunc("/api/upload", httpserver.AuthMiddleware(httpserver.UploadHandler))
	server.HandleFunc("/api/download", httpserver.AuthMiddleware(httpserver.DownloadHandler))
	server.HandleFunc("/api/nickname", httpserver.AuthMiddleware(httpserver.UpdateNicknameHandler))
	corsHandler := corsMiddleware(server)
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}

// Middleware function to handle CORS headers
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow requests from any origin
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Allow specific HTTP methods
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		// Allow specific headers
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
