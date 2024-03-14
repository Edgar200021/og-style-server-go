package db

import (
	"context"
	"errors"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"og-style/models"
	"og-style/types"
)

type UserStorage interface {
	Get(id int) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Create(data *types.CreateUser) (int, error)
}

type UserPgStorage struct {
	DB *pgxpool.Pool
}

func (u *UserPgStorage) Get(id int) (*models.User, error) {
	var user models.User

	if err := pgxscan.Get(context.Background(), u.DB, &user, `SELECT * FROM users WHERE id = $1`, id); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
	}

	return &user, nil
}
func (u *UserPgStorage) GetByEmail(email string) (*models.User, error) {
	var user models.User

	if err := pgxscan.Get(context.Background(), u.DB, &user, `SELECT * FROM users WHERE email = $1`, email); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
	}

	return &user, nil
}

func (u *UserPgStorage) Create(data *types.CreateUser) (int, error) {
	var userId int

	if err := u.DB.QueryRow(context.Background(), "INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id ", data.Email, data.Password).Scan(&userId); err != nil {
		return 0, err
	}

	return userId, nil
}
