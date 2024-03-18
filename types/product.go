package types

type CreateProduct struct {
	Name        string   `json:"name" validate:"required,lte=60"`
	Description string   `json:"description" validate:"required,lte=1000"`
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
	Name        string   `json:"name" validate:"lte=60"`
	Description string   `json:"description" validate:"lte=1000"`
	Price       int      `json:"price" validate:"number,min=1000"`
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

type ProductFilters struct {
	Size       []string `json:"size"`
	Colors     []string `json:"colors"`
	MinPrice   int      `json:"minPrice" db:"min_price"`
	MaxPrice   int      `json:"maxPrice" db:"max_price"`
	BrandsId   []int    `json:"brandsId" db:"brands_id"`
	BrandsName []string `json:"brandsName" db:"brands_name"`
}
