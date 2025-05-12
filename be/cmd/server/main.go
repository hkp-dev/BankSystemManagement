package main

import (
	"app/be/internal/hub"
	"app/be/internal/srv/transaction"
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
	transService := transaction.NewTransactionService(appHub)
	transHandler := transaction.NewTransactionHandler(transService)
	http.HandleFunc("/api/user/register", authHandler.Register)
	http.HandleFunc("/api/user/login", authHandler.Login)
	http.HandleFunc("/api/user/2fa/setup-get", authHandler.GetTOTPSetup)
	http.HandleFunc("/api/user/2fa/setup-post", authHandler.VerifyFirstTimeOTP)
	http.HandleFunc("/api/user/2fa/verify", authHandler.VerifyLoginOTP)
	http.HandleFunc("/api/logout", authHandler.Logout)
	http.HandleFunc("/api/account/info", authHandler.GetUsrInfo)
	http.HandleFunc("/api/transactions", transHandler.GetTransactions)
	http.HandleFunc("/api/verify-transaction", transHandler.PostTransaction)
	// Page handlers
	http.HandleFunc("/user/login", serveTemplate("../fe/html/user/login.html"))
	http.HandleFunc("/user/register", serveTemplate("../fe/html/user/register.html"))
	http.HandleFunc("/user/2fa/setup", serveTemplate("../fe/html/user/2fa_setup.html"))
	http.HandleFunc("/user/2fa/verify", serveTemplate("../fe/html/user/2fa_verify_login.html"))
	http.HandleFunc("/user/home", serveTemplate("../fe/html/user/home.html"))
	http.HandleFunc("/user/verify-otp-transaction", serveTemplate("../fe/html/user/2fa_verify_transaction.html"))
	http.HandleFunc("/user/transfer-money", serveTemplate("../fe/html/user/transferMoney.html"))

	fs := http.FileServer(http.Dir("fe"))
	http.Handle("/fe/", http.StripPrefix("/fe/", fs))
	// Page handlers
	log.Println("Server is running on: http://localhost:8080/user/login")
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
