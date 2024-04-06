package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/fernandoDelPo/go-web-users-response/response"
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
		case http.MethodPatch:
			switch pathSize {
			case 4:
				end = endpoints.Update
				deco = decoUpdateUser
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

func encodeResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {

	r := resp.(response.Response)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(r.StatusCode())

	return json.NewEncoder(w).Encode(resp)

}

func encondeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp := err.(response.Response)

	w.WriteHeader(resp.StatusCode())
	
	_ = json.NewEncoder(w).Encode(resp)
}

func decoUpdateUser(ctx context.Context, r *http.Request) (interface{}, error) {
	var req user.UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("error decoding request body:  '%v'", err.Error())
	}

	params := ctx.Value("params").(map[string]string)

	id, err := strconv.ParseUint(params["userID"], 10, 64)

	if err != nil {
		return nil, err
	}
	req.ID = id

	return req, nil

}

func InvalidMethod(w http.ResponseWriter) {

	status := http.StatusNotFound
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"status": %d, "message": "method doesn't exits"}`, status)
}
