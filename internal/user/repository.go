package user

import (
	"context"
	"errors"
	"log"
	"slices"

	"github.com/fernandoDelPo/go_web_users/internal/domain"
)

type DB struct {
	Users     []domain.User
	MaxUserID uint64
}

type (
	Repository interface {
		Create(ctx context.Context, user *domain.User) error
		GetAll(ctx context.Context) ([]domain.User, error)
		Get(ctx context.Context, id uint64) (*domain.User, error)
		// Update(ctx context.Context, user domain.User) error
		// Delete(ctx context.Context, id uint64) error
	}
	dbRepo struct {
		db  DB
		log *log.Logger
	}
)

func NewDBRepository(db DB, l *log.Logger) Repository {
	return &dbRepo{
		db:  db,
		log: l,
	}
}

func (r *dbRepo) Create(ctx context.Context, user *domain.User) error {
	r.db.MaxUserID++
	user.ID = r.db.MaxUserID
	r.db.Users = append(r.db.Users, *user)
	r.log.Println("Created User with ID ", user.ID)
	return nil
}

func (r *dbRepo) GetAll(ctx context.Context) ([]domain.User, error) {
	r.log.Println("Retrieving all users")
	return r.db.Users, nil
}

func (r *dbRepo) Get(ctx context.Context, id uint64) (*domain.User, error) {
	index := slices.IndexFunc(r.db.Users, func(v domain.User) bool {
		return v.ID == id
	})

	if index < 0 {
		return nil, errors.New("user doesn`t exist")
	}
	return  &r.db.Users[index], nil
}
