package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/fernandoDelPo/go_web_users/internal/user"
	"github.com/fernandoDelPo/go_web_users/pkg/bootstrapt"
	"github.com/fernandoDelPo/go_web_users/pkg/handler"
)

func main() {

	server := http.NewServeMux()

	db := bootstrapt.NewDB()
	logger := bootstrapt.NewLogger()

	repo := user.NewDBRepository(db, logger)
	service := user.NewService(logger, repo)

	ctx := context.Background()

	handler.NewUserHTPPServer(ctx, server,user.MakeEndpoints(ctx, service))

	fmt.Println("Starting server on port :8080")
	log.Fatal(http.ListenAndServe(":8080", server))

}
