package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fernandoDelPo/go_web_users/internal/user"
	"github.com/fernandoDelPo/go_web_users/pkg/transport"
)

func NewUserHTPPServer(ctx context.Context, router *http.ServeMux, endpoints user.Endpoints) {

	router.HandleFunc("/users", UserServer(ctx, endpoints))

}

func UserServer(ctx context.Context, endpoints user.Endpoints) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		tran := transport.New(w, r, ctx)
		switch r.Method {
		case http.MethodGet:
			tran.Server(
				transport.Endpoint(endpoints.GetAll),
				decodeGetAllUsers,
				encodeResponse,
				encondeError)
			return
		}
		InvalidMethod(w)
	}
}

func decodeGetAllUsers(ctx context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	data, err := json.Marshal(response)
	if err != nil {
		return err
	}

	status := http.StatusOK

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(w, `{"status": %d, "data": %s}`, status, data)
	return  nil

}

func encondeError(_ context.Context, err error, w http.ResponseWriter) {
	status := http.StatusInternalServerError
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(w, `{"status": %d, "message":"%s"}`, status, err.Error())
}

// func MsgResponse(w http.ResponseWriter, status int, message string) {

// 	w.WriteHeader(status)
// 	fmt.Fprintf(w, `{"status": %d, "message": "%s"}`, status, message)
// }

// func DataResponse(w http.ResponseWriter, status int, users interface{}) {
// 	value, err := json.Marshal(users)
// 	if err != nil {
// 		MsgResponse(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}
// 	w.WriteHeader(status)
// 	fmt.Fprintf(w, `{"status": %d, "data": %s}`, status, value)
//}

func InvalidMethod(w http.ResponseWriter) {

	status := http.StatusNotFound
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"status": %d, "message": "method doesn't exits"}`, status)
}
