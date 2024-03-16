package models

import "github.com/jackc/pgx/v5/pgtype"

type Product struct {
	ID              int               `json:"id" db:"id"`
	Name            string            `json:"name" db:"name"`
	Description     string            `json:"description" db:"description"`
	Price           int               `json:"price" db:"price"`
	DiscountedPrice *int              `json:"discountedPrice,omitempty" db:"discounted_price"`
	Discount        *int              `json:"discount,omitempty" db:"discount"`
	Images          []string          `json:"images" db:"images"`
	Size            []string          `json:"size" db:"size"`
	Category        string            `json:"-" db:"category"`
	SubCategory     string            `json:"-" db:"sub_category"`
	Materials       pgtype.JSONBCodec `json:"materials" db:"materials"`
	Colors          []string          `json:"colors" db:"colors"`
	Brand           int               `json:"-" db:"colors"`
}
