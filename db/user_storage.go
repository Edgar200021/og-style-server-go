package db

import (
	"context"
	"errors"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"og-style/models"
	"og-style/types"
	"time"
)

type UserStorage interface {
	Get(id int) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Create(data *types.CreateUser) (int, error)
	//Update(data *types.UpdateUser) error
	UpdatePassword(userId int, password string) error
	UpdatePasswordExpires(userId int, passwordResetExpires time.Time) error
	DeletePasswordResetExpires(userId int) error
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
func (u *UserPgStorage) UpdatePassword(userId int, password string) error {

	if _, err := u.DB.Query(context.Background(), "UPDATE users SET password = $1 WHERE id = $2", password, userId); err != nil {
		return err
	}

	return nil
}
func (u *UserPgStorage) UpdatePasswordExpires(userId int, passwordResetExpires time.Time) error {
	if _, err := u.DB.Query(context.Background(), `UPDATE users SET password_reset_expires = $1 WHERE id = $2`, passwordResetExpires, userId); err != nil {
		return err
	}

	return nil
}

//func (u *UserPgStorage) Update(data *types.UpdateUser) error {
//
//	if _,err := u.DB.Query(context.Background(), `UPDATE users Set email = COALESCE($1, email), password = COALESCE($2, password), name = COALESCE($3, name), avatar = COALESCE($4, avatar)`, data.Email,data.Password,data.Name,data.Email); err != nil {
//		return err
//	}
//
//	return nil
//}

func (u *UserPgStorage) DeletePasswordResetExpires(userId int) error {

	if _, err := u.DB.Query(context.Background(), `UPDATE users SET password_reset_expires = null WHERE id = $1`, userId); err != nil {
		return err
	}

	return nil
}
