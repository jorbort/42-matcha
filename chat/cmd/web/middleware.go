package main

import (
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func commonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // Your main service origin
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true") // If you need to send cookies

		// Your existing security headers
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; "+
				"style-src 'self' fonts.googleapis.com; "+
				"font-src fonts.gstatic.com; "+
				"connect-src 'self' http://localhost:3000 ws://localhost:3001; "+ // Add WebSocket connection
				"img-src 'self' data:; "+ // If you need to handle images
				"script-src 'self' 'unsafe-inline' https://unpkg.com") // For htmx if loading from CDN
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		w.Header().Set("Server", "Go")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			remoteAddr = r.RemoteAddr
			proto      = r.Proto
			method     = r.Method
			url        = r.URL.RequestURI()
		)
		log.Printf("received request from %s %s %s %s", remoteAddr, proto, method, url)
		next.ServeHTTP(w, r)
	})
}

func authHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access-token")
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		accesToken := cookie.Value
		str, err := jwt.Parse(accesToken, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET_KEY")), nil
		})
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if !str.Valid {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
