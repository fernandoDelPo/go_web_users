package user

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type (
	Controller func(w http.ResponseWriter, r *http.Request)

	Endpoints struct {
		Create *Controller
		GetAll *Controller
		Update *Controller
		Delete *Controller
	}

	CreateRequest struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	}
)

func MakeEndpoints(ctx context.Context, s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			GetAllUser(ctx, s, w)
		case http.MethodPost:
			decode := json.NewDecoder(r.Body)
			var req CreateRequest
			err := decode.Decode(&req)
			if err != nil {
				MsgResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
			PostUser(ctx, s, w, req)

		default:
			InvalidMethod(w)
		}
	}
}

func GetAllUser(ctx context.Context, s Service, w http.ResponseWriter) {
	users, err := s.GetAll(ctx)
	if err != nil {
		MsgResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	DataResponse(w, http.StatusOK, users)
}

func PostUser(ctx context.Context, s Service, w http.ResponseWriter, data interface{}) {

	req := data.(CreateRequest)

	if req.FirstName == "" {
		MsgResponse(w, http.StatusBadRequest, "First name is required.")
		return
	}

	if req.LastName == "" {
		MsgResponse(w, http.StatusBadRequest, "Last name is required.")
		return
	}
	if req.Email == "" {
		MsgResponse(w, http.StatusBadRequest, "Email is required.")
		return
	}

	// id := len(users) + 1
	// idString := strconv.Itoa(id)
	// user.ID = idString
	// users = append(users, user)

	user, err := s.Create(ctx, req.FirstName, req.LastName, req.Email)
	if err != nil {
		MsgResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	DataResponse(w, http.StatusCreated, user)

}

func MsgResponse(w http.ResponseWriter, status int, message string) {

	w.WriteHeader(status)
	fmt.Fprintf(w, `{"status": %d, "message": "%s"}`, status, message)
}

func DataResponse(w http.ResponseWriter, status int, users interface{}) {
	value, err := json.Marshal(users)
	if err != nil {
		MsgResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"status": %d, "data": %s}`, status, value)
}

func InvalidMethod(w http.ResponseWriter) {

	status := http.StatusNotFound
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"status": %d, "message": "method doesn't exits"}`, status)
}
