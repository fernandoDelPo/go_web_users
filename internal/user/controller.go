package user

import (
	"context"
	"errors"

	"github.com/fernandoDelPo/go-web-users-response/response"
)

type (
	Controller func(ctx context.Context, request interface{}) (interface{}, error)

	Endpoints struct {
		Create Controller
		GetAll Controller
		Get    Controller
		Update Controller
		Delete Controller
	}

	GetReq struct {
		ID uint64
	}
	CreateRequest struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	}

	UpdateRequest struct {
		ID        uint64
		FirstName *string `json:"first_name"`
		LastName  *string `json:"last_name"`
		Email     *string `json:"email"`
	}
)

func MakeEndpoints(ctx context.Context, s Service) Endpoints {
	return Endpoints{
		Create: makeCreateEndpoint(ctx, s),
		GetAll: makeGetAllEndpoint(ctx, s),
		Get:    makeGetEndpoint(ctx, s),
		Update: makeUpdateEndpoint(ctx, s),
	}
}

func makeCreateEndpoint(ctx context.Context, s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateRequest)

		if req.FirstName == "" {
			return nil, response.BadRequest(ErrFirstNameRequired.Error())
		}

		if req.LastName == "" {
			return nil, response.BadRequest(ErrLastNameRequired.Error())
		}

		// if req.Email == "" {
		// 	return nil, errors.New("field email is required")
		// }

		// id := len(users) + 1
		// idString := strconv.Itoa(id)
		// user.ID = idString
		// users = append(users, user)

		user, err := s.Create(ctx, req.FirstName, req.LastName, req.Email)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.Created("success", user), nil

	}
}

func makeGetAllEndpoint(ctx context.Context, s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		users, err := s.GetAll(ctx)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("success", users), nil
	}
}

func makeGetEndpoint(ctx context.Context, s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetReq)
		user, err := s.Get(ctx, req.ID)
		if err != nil {
			if errors.As(err, &ErrNotFound{}){
				return  nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("success", user), nil
	}
}

func makeUpdateEndpoint(ctx context.Context, s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateRequest)

		if req.FirstName != nil && *req.FirstName == "" {
			return nil, response.BadRequest(ErrFirstNameRequired.Error())
		}

		if req.LastName != nil && *req.LastName == "" {
			return nil, response.BadRequest(ErrLastNameRequired.Error())
		}

		if err := s.Update(ctx, req.ID, req.FirstName, req.LastName, req.Email); err != nil {
			if errors.As(err, &ErrNotFound{}){
				return  nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}
		return response.OK("success", nil), nil
	}
}
