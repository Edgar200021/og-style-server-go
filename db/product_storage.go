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
	args := make([]any, 0, 2)
	limit, page := 8, 1

	fmt.Println(params.Page)

	if params.Limit != 0 {
		limit = params.Limit
	}

	if params.Page != 0 {
		page = params.Page
	}

	if len(params.Size) != 0 {
		query += ` WHERE size && ($1)`
		args = append(args, params.Size)
	}

	if len(params.Colors) != 0 {
		if utils.IsContainsSubstring(query, "where") {
			query += ` AND colors && ($2)`
		} else {
			query += ` WHERE colors && ($1)`
		}

		args = append(args, params.Colors)
	}

	if len(params.Brand) != 0 {
		if utils.IsContainsSubstring(query, "where") {
			query += fmt.Sprintf(` AND brand = ANY($%d)`, len(args)+1)
		} else {
			query += ` WHERE brand = ANY($1)`
		}

		args = append(args, params.Brand)
	}

	if params.Name != "" {
		if utils.IsContainsSubstring(query, "where") {
			query += fmt.Sprintf(` AND name like $%d)`, len(args)+1)
		} else {
			query += ` WHERE name like $1`
		}
		args = append(args, "%"+params.Name+"%")
	}

	if params.Category != "" {
		if utils.IsContainsSubstring(query, "where") {
			query += fmt.Sprintf(` AND category = $%d)`, len(args)+1)
		} else {
			query += ` WHERE category = $1`
		}
		args = append(args, params.Category)
	}

	if params.SubCategory != "" {
		if utils.IsContainsSubstring(query, "where") {
			query += fmt.Sprintf(` AND sub_category = $%d)`, len(args)+1)
		} else {
			query += ` WHERE sub_category = $1`
		}
		args = append(args, params.SubCategory)
	}

	query += fmt.Sprintf(` LIMIT $%d OFFSET $%d`, len(args)+1, len(args)+2)
	args = append(args, limit, (page*limit)-limit)

	if err := pgxscan.Select(context.Background(), p.DB, &products, query, args...); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return products, err
		}
	}

	return products, nil
}
func (p *ProductPgStorage) Create(data *types.CreateProduct) error {
	var discountedPrice int

	if data.Discount != 0 {
		discountedPrice = data.Price - ((data.Price * data.Discount) / 100)
	}

	if _, err := p.DB.Query(context.Background(), `INSERT INTO product (name,description,price,discounted_price,discount,images,size,category,sub_category,colors,brand,materials) VALUES ($1, $2, $3, CASE WHEN $4 = 0 THEN NULL ELSE $4 END,CASE WHEN $5 = 0 THEN NULL ELSE $5 END, $6, $7, $8, $9, $10, $11, $12)`, data.Name, data.Description, data.Price, discountedPrice, data.Discount, data.Images, data.Size, data.Category, data.SubCategory, data.Colors, data.Brand, data.Materials); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
func (p *ProductPgStorage) Update(id int, data *types.UpdateProduct) error {

	if _, err := p.DB.Query(context.Background(), `UPDATE product p SET 
                     name=COALESCE(NULLIF($1,''), p.name),
                     description=COALESCE(NULLIF($2,''), p.description),
                     price=COALESCE(NULLIF($3,0), p.price),
                     discounted_price = CASE WHEN $3 != 0 AND $4 != 0 THEN $3 - (($3 * $4) / 100)
                         							  WHEN $3 != 0 AND $4 = 0 THEN $3 - (($3 * p.discount) / 100)
                         							  WHEN $3 = 0 AND $4 != 0 THEN p.price - ((p.price * $4) / 100)
                         							  ELSE p.discounted_price
												END,
                     discount=COALESCE(NULLIF($4,0), p.discount),
                     images=COALESCE(NULLIF($5, '{}'::TEXT[]), p.images),
                     size=COALESCE(NULLIF($6, '{}'::TEXT[]), p.size),
                     category=COALESCE(NULLIF($7,''), p.category),
                     sub_category=COALESCE(NULLIF($8,''), p.sub_category),
                     materials=COALESCE(NULLIF($9, '{}'::TEXT[]), p.materials),
                     colors=COALESCE(NULLIF($10,'{}'::TEXT[]),p.colors),
                     brand=COALESCE(NULLIF($11, 0), p.brand) WHERE id = $12 `, data.Name, data.Description, data.Price, data.Discount, data.Images, data.Size, data.Category, data.SubCategory, data.Materials, data.Colors, data.Brand, id); err != nil {
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
