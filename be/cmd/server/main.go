package main

import (
	"app/be/internal/hub"
	"app/be/internal/srv/user"
	"context"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	ctx := context.Background()
	appHub, _ := hub.New(ctx)
	authService := user.NewAuthService(appHub)
	authHandler := user.NewAuthHandler(authService)
	http.HandleFunc("/api/user/register", authHandler.Register)
	http.HandleFunc("/api/user/login", authHandler.Login)
	http.HandleFunc("/api/user/2fa/setup-get", authHandler.GetTOTPSetup)
	http.HandleFunc("/api/user/2fa/setup-post", authHandler.VerifyFirstTimeOTP)
	http.HandleFunc("/api/user/2fa/verify", authHandler.VerifyLoginOTP)
	http.HandleFunc("/api/logout", authHandler.Logout)

	// Page handlers
	http.HandleFunc("/user/login", serveTemplate("../fe/html/user/login.html"))
	http.HandleFunc("/user/register", serveTemplate("../fe/html/user/register.html"))
	http.HandleFunc("/user/2fa/setup", serveTemplate("../fe/html/user/2fa_setup.html"))
	http.HandleFunc("/user/2fa/verify", serveTemplate("../fe/html/user/2fa_verify.html"))

	fs := http.FileServer(http.Dir("fe"))
	http.Handle("/fe/", http.StripPrefix("/fe/", fs))

	// Page handlers
	log.Println("Server is running on: http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func serveTemplate(tmplPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		absPath, err := filepath.Abs(tmplPath)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles(absPath)
		if err != nil {
			http.Error(w, "Template not found", http.StatusNotFound)
			return
		}

		data := map[string]interface{}{
			"Title": "Banking System",
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			return
		}
	}
}
