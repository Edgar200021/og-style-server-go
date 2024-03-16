package types

import "github.com/jackc/pgx/v5/pgtype"

type CreateProduct struct {
	Name        string            `json:"name" validate:"required,lte=40"`
	Description string            `json:"description" validate:"required,lte=200"`
	Price       int               `json:"price" validate:"required,numeric,min=0"`
	Discount    *int              `json:"discount,omitempty" validate:"number,min=1,max=99"`
	Images      []string          `json:"images" validate:"required,dive,len=4"`
	Size        []string          `json:"size" validate:"required,dive"`
	Category    string            `json:"-" validate:"required,oneof=одежда обувь"`
	SubCategory string            `json:"-" validate:"required"`
	Materials   pgtype.JSONBCodec `json:"materials" validate:"required,json"`
	Colors      []string          `json:"colors"  validate:"required,dive,hexcolor"`
	Brand       int               `json:"brand" validate:"required,min=1,max=6"`
}

type UpdateProduct struct {
	Name        string            `json:"name" validate:"lte=40"`
	Description string            `json:"description" validate:"lte=200"`
	Price       int               `json:"price" validate:"numeric,min=0"`
	Discount    *int              `json:"discount,omitempty" validate:"number,min=1,max=99"`
	Images      []string          `json:"images" validate:"dive,len=4"`
	Size        []string          `json:"size" validate:"dive"`
	Category    string            `json:"-" validate:"oneof=одежда обувь"`
	SubCategory string            `json:"-" `
	Materials   pgtype.JSONBCodec `json:"materials" validate:"json"`
	Colors      []string          `json:"colors"  validate:"dive,hexcolor"`
	Brand       int               `json:"brand" validate:"min=1,max=6"`
}

type GetProductsParams struct {
	Name        string   `json:"name"`
	Category    string   `json:"category" validate:"oneof=одежда обувь"`
	SubCategory string   `json:"sub_category" validate:"required_with=Category"`
	Brand       []string `json:"brand" validate:"dive,min=1,max=6"`
	Limit       int      `json:"limit" validate:"min=1"`
	Page        int      `json:"page" validate:"min=1"`
	Size        []string `json:"size,omitempty" validate:"dive"`
	Colors      []string `json:"colors" validate:"dive,hexcolor"`
}
