package main

import (
	"entrytask/backend/http"
	ziherpc "entrytask/backend/rpc"
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
	server.HandleFunc("/api/users/login", httpserver.LoginHandler)
	server.HandleFunc("/api/users/profile", httpserver.AuthMiddleware(httpserver.ProfileHandler))
	server.HandleFunc("/api/users/avatar", httpserver.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			httpserver.DownloadHandler(w, r)
		case "POST":
			httpserver.UploadHandler(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}))
	server.HandleFunc("/api/users/nickname", httpserver.AuthMiddleware(httpserver.UpdateNicknameHandler))
	corsHandler := corsMiddleware(server)
	log.Println("http server is running")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}

// Middleware function to handle CORS headers
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow requests from any origin
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Allow specific HTTP methods
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT")

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
