package user

import (
	"context"
	"github.com/fernandoDelPo/go_web_users/internal/domain"
	"log"
)

type (
	Service interface {
		Create(ctx context.Context, firstName, lastName, email string) (*domain.User, error)
		GetAll(ctx context.Context) ([]domain.User, error)
		Get(ctx context.Context, id uint64) (*domain.User, error)
		Update(ctx context.Context, id uint64, firstName, lastName, email *string) error
	}

	service struct {
		log  *log.Logger
		repo Repository
	}
)

func NewService(l *log.Logger, repo Repository) Service {
	return &service{
		log:  l,
		repo: repo,
	}
}

func (s service) Create(ctx context.Context, firstName, lastName, email string) (*domain.User, error) {
	log.Printf("test create service")
	user := &domain.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	s.log.Println("Created User with ID ", user.ID)
	return user, nil
}

func (s service) GetAll(ctx context.Context) ([]domain.User, error) {
	users, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	s.log.Println(len(users), " users were fetched")
	return users, nil
}

func (s service) Get(ctx context.Context, id uint64) (*domain.User, error) {
	user, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s service) Update(ctx context.Context, id uint64, firstName, lastName, email *string) error {
	if err := s.repo.Update(ctx, id, firstName, lastName, email); err != nil {
		return err
	}
	return nil
}