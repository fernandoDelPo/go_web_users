package user

import (
	"context"
	"log"
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
