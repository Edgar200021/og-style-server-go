package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"og-style/models"
	"og-style/types"
	"og-style/utils"
	"strings"
)

type ProductStorage interface {
	Get(id int) (models.Product, error)
	GetAll(params types.GetProductsParams) ([]*models.Product, error)
	Create(data *types.CreateProduct) error
	Update(id int, data *types.UpdateProduct) error
	Delete(id int) error
}

type ProductPgStorage struct {
	DB *pgxpool.Pool
}

func (p *ProductPgStorage) Get(id int) (models.Product, error) {
	var product models.Product

	if err := pgxscan.Get(context.Background(), p.DB, &product, `SELECT * FROM product WHERE id = $1`, id); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return product, err
		}
	}

	return product, nil
}
func (p *ProductPgStorage) GetAll(params types.GetProductsParams) ([]*models.Product, error) {
	products := []*models.Product{}
	query := `SELECT * FROM product`
	args := make([]any, 2)
	limit, page := 8, 1

	if params.Limit != 0 {
		limit = params.Limit
	}

	if params.Page != 0 {
		page = params.Page
	}

	if len(params.Size) != 0 {
		query += ` WHERE size in ($1)`
		args = append(args, strings.Join(params.Size, ", "))
	}

	if len(params.Colors) != 0 {
		if utils.IsContainsSubstring(query, "where") {
			query += ` AND colors in ($2)`
		} else {
			query += ` WHERE colors in ($1)`
		}

		args = append(args, strings.Join(params.Colors, ", "))
	}

	if len(params.Brand) != 0 {
		if utils.IsContainsSubstring(query, "where") {
			query += fmt.Sprintf(` AND brand in ($%d)`, len(args)+1)
		} else {
			query += ` WHERE brand in ($1)`
		}

		args = append(args, strings.Join(params.Size, ", "))
	}

	if params.Name != "" {
		if utils.IsContainsSubstring(query, "where") {
			query += fmt.Sprintf(` AND name like %d)`, len(args)+1)
		} else {
			query += ` WHERE name like $1`
		}
		args = append(args, "%"+params.Name+"%")
	}

	if params.Category != "" {
		if utils.IsContainsSubstring(query, "where") {
			query += fmt.Sprintf(` AND category = %d)`, len(args)+1)
		} else {
			query += ` WHERE category = $1`
		}
		args = append(args, params.Category)
	}

	if params.SubCategory != "" {
		if utils.IsContainsSubstring(query, "where") {
			query += fmt.Sprintf(` AND sub_category = %d)`, len(args)+1)
		} else {
			query += ` WHERE sub_category = $1`
		}
		args = append(args, params.SubCategory)
	}

	query += fmt.Sprintf(` LIMIT %d OFFSET %d`, len(args)+1, len(args)+2)
	args = append(args, limit, (page*limit)-page)

	if err := pgxscan.Select(context.Background(), p.DB, &products, query, args...); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return products, err
		}
	}

	return products, nil
}
func (p *ProductPgStorage) Create(data *types.CreateProduct) error {
	if _, err := p.DB.Query(context.Background(), `INSERT INTO product (name, description, price, discount, images, size, category, sub_category, materials, colors, brand)  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`, data.Name, data.Description, data.Price, data.Discount, data.Images, data.Size, data.Category, data.SubCategory, data.Materials, data.Colors, data.Brand); err != nil {
		return err
	}

	return nil
}
func (p *ProductPgStorage) Update(id int, data *types.UpdateProduct) error {

	if _, err := p.DB.Query(context.Background(), `UPDATE product SET name=COALESCE($1, name), description=COALESCE($2, description), price=COALESCE($3, price), discount=COALESCE($4, discount), images=COALESCE($5, images), size=COALESCE($6, size), category=COALESCE($7, category), sub_category=COALESCE($8, sub_category), materials=COALESCE($9, materials), colors=COALESCE($10, colors), brand=COALESCE($11, brand) WHERE id = $12  `, data.Name, data.Description, data.Price, data.Discount, data.Images, data.Size, data.Category, data.SubCategory, data.Materials, data.Colors, data.Brand, id); err != nil {
		return err
	}

	return nil
}
func (p *ProductPgStorage) Delete(id int) error {
	if _, err := p.DB.Query(context.Background(), `DELETE FROM product WHERE id = $1`, id); err != nil {
		return err
	}
	return nil
}
