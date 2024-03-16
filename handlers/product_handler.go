package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"og-style/processors"
	"og-style/types"
	"og-style/utils"
	"strconv"
)

type ProductHandler struct {
	ProductProcessor processors.ProductProcessor
}

func (p *ProductHandler) Get(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		utils.BadRequestError(w, err)
		return
	}

	if product, err := p.ProductProcessor.Get(int(id)); err != nil {
		utils.BadRequestError(w, err)
	} else {
		utils.SendJSON(w, product, http.StatusOK)
	}

}

func (p *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {

	m := make(map[string]any)

	for key, val := range r.URL.Query() {
		if key == "colors" || key == "size" || key == "brand" {
			m[key] = val
		} else {
			m[key] = val[0]
		}
	}

	fmt.Println(m)

	var getProductsParams types.GetProductsParams
	if err := mapstructure.Decode(m, &getProductsParams); err != nil {
		utils.BadRequestError(w, err)
		return
	}

	fmt.Println(getProductsParams)

	utils.SendJSON(w, "ok", http.StatusOK)

}

func (p *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {

	var body types.CreateProduct

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.BadRequestError(w, err)
		return
	}

	if err := p.ProductProcessor.Create(&body); err != nil {
		utils.BadRequestError(w, err)
		return
	}

	utils.SendJSON(w, "success", http.StatusCreated)
}
