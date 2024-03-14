package db

import (
	"context"
	"errors"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"og-style/models"
)

type TokenStorage interface {
	Get(userId int) (*models.Token, error)
	Create(userId int, token string) error
}

type TokenPgStorage struct {
	DB *pgxpool.Pool
}

func (t *TokenPgStorage) Get(userId int) (*models.Token, error) {
	var token models.Token

	if err := pgxscan.Get(context.Background(), t.DB, &token, `SELECT * FROM token WHERE user_id = $1`, userId); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
	}

	return &token, nil
}
func (t *TokenPgStorage) Create(userId int, token string) error {
	if _, err := t.DB.Query(context.Background(), "INSERT INTO token (user_id, refresh_token) VALUES ($1, $2)", userId, token); err != nil {
		return err
	}

	return nil
}
