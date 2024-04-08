package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/fernandoDelPo/go_web_users/internal/user"
	"github.com/fernandoDelPo/go_web_users/pkg/bootstrapt"
	"github.com/fernandoDelPo/go_web_users/pkg/handler"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	db, err := bootstrapt.NewDB()
	if err != nil {
		log.Fatalf("Error creating database connection: %s", err)
	}

	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping the DB: %v\n", err)
	} else {
		log.Println("Connected to the Database")
	}

	logger := bootstrapt.NewLogger()

	repo := user.NewDBRepository(db, logger)
	service := user.NewService(logger, repo)

	ctx := context.Background()

	h := handler.NewUserHTPPServer(user.MakeEndpoints(ctx, service))

	port := os.Getenv("PORT")
	fmt.Println("Starting server on port: ", port)

	address := fmt.Sprintf("127.0.0.1:%s", port)

	srv := &http.Server{
		Addr: address, 
		Handler: accessControl(h),
	}
	log.Fatal(srv.ListenAndServe())

}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS, HEAD, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept,Authorization,Cache-Control,Content-Type,DNT,If-Modified-Since,Keep-Alive,Origin,User-Agent,X-Requested-With")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
