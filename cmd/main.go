package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/fernandoDelPo/go_web_users/internal/domain"
	"github.com/fernandoDelPo/go_web_users/internal/user"
)

func main() {

	server := http.NewServeMux()

	db := user.DB{
		Users: []domain.User{{
			ID:        1,
			FirstName: "fernando",
			LastName:  "Del Pozzi",
			Email:     "fernandodelpozzi@example.com",
		}, {
			ID:        2,
			FirstName: "Leticia",
			LastName:  "Caceres",
			Email:     "leticia@example.com",
		}, {
			ID:        3,
			FirstName: "Francesca",
			LastName:  "Del Pozzi",
			Email:     "fran@example.com",
		}, {
			ID:        4,
			FirstName: "Maxima",
			LastName:  "Perez",
			Email:     "mperez@example.com",
		}},
		MaxUserID: 4,
	}

	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	repo := user.NewDBRepository(db, logger)
	service := user.NewService(logger, repo)

	ctx := context.Background()

	server.HandleFunc("/users", user.MakeEndpoints(ctx, service))
	fmt.Println("Starting server on port :8080")
	log.Fatal(http.ListenAndServe(":8080", server))

}
