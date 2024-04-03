package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/fernandoDelPo/go_web_users/internal/user"
	"github.com/fernandoDelPo/go_web_users/pkg/transport"
)

func NewUserHTPPServer(ctx context.Context, router *http.ServeMux, endpoints user.Endpoints) {

	router.HandleFunc("/users/", UserServer(ctx, endpoints))

}

func UserServer(ctx context.Context, endpoints user.Endpoints) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		url := r.URL.Path
		log.Println(r.Method, ": ", url)

		path, pathSize := transport.Clean(url)

		params := make(map[string]string)

		if pathSize == 4 && path[2] != "" {
			params["userID"] = path[2]
		}

		tran := transport.New(w, r, context.WithValue(ctx, "params", params))

		var end user.Controller
		var deco func(ctx context.Context, r *http.Request) (interface{}, error)

		switch r.Method {
		case http.MethodGet:
			switch pathSize {
			case 3:
				end = endpoints.GetAll
				deco = decodeGetAllUsers
			case 4:
				end = endpoints.Get
				deco = decodeGetUser
			}

		case http.MethodPost:
			switch pathSize {
			case 3:
				end = endpoints.Create
				deco = decodeCreateUser
			}
		}
		if end != nil && deco != nil {
			tran.Server(
				transport.Endpoint(end),
				deco,
				encodeResponse,
				encondeError)
		} else {
			InvalidMethod(w)
		}

	}
}

func decodeGetUser(ctx context.Context, r *http.Request) (interface{}, error) {
	params := ctx.Value("params").(map[string]string)

	id, err := strconv.ParseUint(params["userID"], 10, 64)

	if err != nil {
		return nil, err
	}

	return user.GetReq{
		ID: id,
	}, nil
}

func decodeGetAllUsers(ctx context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeCreateUser(ctx context.Context, r *http.Request) (interface{}, error) {
	var req user.CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("error decoding request body:  '%v'", err.Error())
	}
	return req, nil
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
	return nil

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
