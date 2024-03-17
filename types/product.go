package types

type CreateProduct struct {
	Name        string   `json:"name" validate:"required,lte=40"`
	Description string   `json:"description" validate:"required,lte=400"`
	Price       int      `json:"price" validate:"required,number,min=1000"`
	Discount    int      `json:"discount,omitempty" validate:"omitempty,number,min=1,max=99"`
	Images      []string `json:"images" validate:"required,len=4,dive"`
	Size        []string `json:"size" validate:"required,dive"`
	Category    string   `json:"category" validate:"required,oneof=одежда обувь"`
	SubCategory string   `json:"subCategory" validate:"required"`
	Materials   []string `json:"materials" validate:"required,dive"`
	Colors      []string `json:"colors"  validate:"required,dive,hexcolor"`
	Brand       int      `json:"brand" validate:"required,min=1,max=6"`
}

type UpdateProduct struct {
	Name        string   `json:"name" validate:"lte=40"`
	Description string   `json:"description" validate:"lte=200"`
	Price       int      `json:"price" validate:"numeric,min=0"`
	Discount    int      `json:"discount,omitempty" validate:"omitempty,number,min=1,max=99"`
	Images      []string `json:"images" validate:"dive,len=4"`
	Size        []string `json:"size" validate:"dive"`
	Category    string   `json:"-" validate:"oneof=одежда обувь"`
	SubCategory string   `json:"-" `
	Materials   []string `json:"materials" validate:"dive"`
	Colors      []string `json:"colors"  validate:"dive,hexcolor"`
	Brand       int      `json:"brand" validate:"min=1,max=6"`
}

type GetProductsParams struct {
	Name        string   `json:"name,omitempty" validate:"omitempty"`
	Category    string   `json:"category,omitempty" validate:"omitempty,oneof=одежда обувь"`
	SubCategory string   `json:"subCategory,omitempty" validate:"required_with=Category"`
	Brand       []int    `json:"brand,omitempty" validate:"omitempty,dive,min=1,max=6"`
	Limit       int      `json:"limit,omitempty" validate:"omitempty,min=1"`
	Page        int      `json:"page,omitempty" validate:"omitempty,min=1"`
	Size        []string `json:"size,omitempty" validate:"omitempty,dive"`
	Colors      []string `json:"colors,omitempty" validate:"omitempty,dive,hexcolor"`
}
