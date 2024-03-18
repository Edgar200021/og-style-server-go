package processors

import (
	"fmt"
	"mime/multipart"
	"og-style/db"
	"og-style/models"
	"og-style/services"
	"og-style/types"
)

type ProductProcessor interface {
	Get(id int) (models.Product, error)
	GetAll(params types.GetProductsParams) ([]*models.Product, error)
	Create(data *types.CreateProduct) error
	Update(id int, data *types.UpdateProduct) error
	Delete(id int) error
	UploadImage(file multipart.File) (string, error)
	GetFilters(category string) (types.ProductFilters, error)
}

type ProductPgProcessor struct {
	ProductStorage db.ProductStorage
	ImageUploader  services.ImageUploaderService
}

func (p *ProductPgProcessor) Get(id int) (models.Product, error) {

	product, err := p.ProductStorage.Get(id)
	if err != nil {
		return product, err
	}

	if product.ID == 0 {
		return product, fmt.Errorf("продукт с ID %d не существует", id)
	}

	return product, nil
}
func (p *ProductPgProcessor) GetAll(params types.GetProductsParams) ([]*models.Product, error) {

	products, err := p.ProductStorage.GetAll(params)
	if err != nil {
		return products, err
	}

	return products, nil
}
func (p *ProductPgProcessor) Create(data *types.CreateProduct) error {
	if err := p.ProductStorage.Create(data); err != nil {
		return err
	}
	return nil
}
func (p *ProductPgProcessor) Update(id int, data *types.UpdateProduct) error {

	product, err := p.Get(id)
	if err != nil {
		return err
	}

	if product.ID == 0 {
		return fmt.Errorf("продукт с ID %d не существует", id)
	}

	err = p.ProductStorage.Update(id, data)
	if err != nil {
		return err
	}

	return nil
}
func (p *ProductPgProcessor) Delete(id int) error {
	product, err := p.Get(id)
	if err != nil {
		return err
	}

	if product.ID == 0 {
		return fmt.Errorf("продукт с ID %d не существует", id)
	}

	err = p.ProductStorage.Delete(id)
	if err != nil {
		return err
	}

	return nil
}

func (p *ProductPgProcessor) UploadImage(file multipart.File) (string, error) {
	if imgUrl, err := p.ImageUploader.Upload(file); err != nil {
		return "", err
	} else {
		return imgUrl, nil
	}
}

func (p *ProductPgProcessor) GetFilters(category string) (types.ProductFilters, error) {
	if filters, err := p.ProductStorage.GetFilters(category); err != nil {
		return filters, err
	} else {
		return filters, nil
	}
}
