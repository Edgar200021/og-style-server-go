package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CartStorage interface {
	Create(userId int) error
	Delete(userId int) error
}

type CartPgStorage struct {
	DB *pgxpool.Pool
}

func (c *CartPgStorage) Create(userId int) error {
	if _, err := c.DB.Query(context.Background(), `INSERT INTO cart (user_id) VALUES ($1)`, userId); err != nil {
		return err
	}
	return nil
}
func (c *CartPgStorage) Delete(userId int) error {
	if _, err := c.DB.Query(context.Background(), `DELETE FROM cart WHERE user_id = $1`, userId); err != nil {
		return err
	}
	return nil
}
