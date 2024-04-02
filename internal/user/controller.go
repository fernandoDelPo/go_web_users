package user

import (
	"context"
	"errors"
)

type (
	Controller func(ctx context.Context, request interface{}) (interface{}, error)

	Endpoints struct {
		Create Controller
		GetAll Controller
		Update Controller
		Delete Controller
	}

	CreateRequest struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	}
)

func MakeEndpoints(ctx context.Context, s Service) Endpoints {
	return Endpoints{
		Create: makeCreateEndpoint(ctx, s),
		GetAll: makeGetAllEndpoint(ctx, s),
	}
}

func makeCreateEndpoint(ctx context.Context, s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateRequest)

		if req.FirstName == "" {
			return nil, errors.New("field first name is required")
		}

		if req.LastName == "" {
			return nil, errors.New("field last name is required")
		}
		if req.Email == "" {
			return nil, errors.New("field email is required")
		}

		// id := len(users) + 1
		// idString := strconv.Itoa(id)
		// user.ID = idString
		// users = append(users, user)

		user, err := s.Create(ctx, req.FirstName, req.LastName, req.Email)
		if err != nil {
			return nil, errors.New("cannot create user")
		}

		return user, nil

	}
}

func makeGetAllEndpoint(ctx context.Context, s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		users, err := s.GetAll(ctx)
		if err != nil {
			return nil, err
		}

		return users, nil
	}

}
