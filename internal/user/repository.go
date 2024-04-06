package user

import (
	"context"
	"database/sql"
	"log"

	"github.com/fernandoDelPo/go_web_users/internal/domain"
)

type (
	Repository interface {
		Create(ctx context.Context, user *domain.User) error
		GetAll(ctx context.Context) ([]domain.User, error)
		Get(ctx context.Context, id uint64) (*domain.User, error)
		Update(ctx context.Context, id uint64, firstName, lastName, email *string) error
		// Delete(ctx context.Context, id uint64) error
	}
	dbRepo struct {
		db  *sql.DB
		log *log.Logger
	}
)

func NewDBRepository(db *sql.DB, l *log.Logger) Repository {
	return &dbRepo{
		db:  db,
		log: l,
	}
}

func (r *dbRepo) Create(ctx context.Context, user *domain.User) error {
	sqlQ := "INSERT INTO users (first_name, last_name, email) VALUES (?,?,?)"
	res, err := r.db.Exec(sqlQ, user.FirstName, user.LastName, user.Email)

	if err != nil {
		r.log.Printf("Error creating the user %s", err.Error())
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		r.log.Println("Could not get the ID of the created user")
		return err
	}

	user.ID = uint64(id)
	r.log.Println("Created User with ID ", user.ID)

	return nil
}

func (r *dbRepo) GetAll(ctx context.Context) ([]domain.User, error) {
	var users []domain.User
	sqlQ := "SELECT id, first_name, last_name, email FROM users"

	rows, err := r.db.Query(sqlQ)

	if err != nil {
		r.log.Println("Failed to retrieve all users from DB : ", err.Error())
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email); err != nil {
			r.log.Printf("Failed to scan row: %v\n", err)
			return nil, err
		}
		users = append(users, u)
	}
	r.log.Println(len(users), " Users retrieved successfully!")

	return users, nil
}

func (r *dbRepo) Get(ctx context.Context, id uint64) (*domain.User, error) {
	// index := slices.IndexFunc(r.db.Users, func(v domain.User) bool {
	// 	return v.ID == id
	// })

	// if index < 0 {
	// 	return nil, ErrNotFound{id}
	// }
	return nil, nil
}

func (r *dbRepo) Update(ctx context.Context, id uint64, firstName, lastName, email *string) error {
	// user, err := r.Get(ctx, id)
	// if err != nil {
	// 	return err
	// }

	// if firstName != nil {
	// 	user.FirstName = *firstName
	// }

	// if lastName != nil {
	// 	user.LastName = *lastName
	// }

	// if email != nil {
	// 	user.Email = *email
	// }

	// r.log.Printf("Updated User %d: %s %s", user.ID, user.FirstName, user.LastName)
	return nil
}
