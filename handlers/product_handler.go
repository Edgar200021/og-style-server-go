package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"og-style/processors"
	"og-style/types"
	"og-style/utils"
	"strconv"
	"strings"
	"sync"
)

const (
	maxFileSize    = 1024 * 1024
	imageFieldName = "image"
)

type ProductHandler struct {
	ProductProcessor processors.ProductProcessor
}

func (p *ProductHandler) Get(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.BadRequestError(w, err)
		return
	}

	if product, err := p.ProductProcessor.Get(id); err != nil {
		utils.BadRequestError(w, err)
	} else {
		utils.SendJSON(w, product, http.StatusOK)
	}

}
func (p *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {

	m, err := p.transformUrlParams(r.URL.Query())
	if err != nil {
		utils.BadRequestError(w, err)
		return
	}

	var getProductsParams types.GetProductsParams

	if err := mapstructure.Decode(m, &getProductsParams); err != nil {
		utils.BadRequestError(w, err)
		return
	}

	if errors := utils.ValidateStruct(getProductsParams); errors != nil {
		utils.SendValidatonErrors(w, errors)
		return
	}

	if products, err := p.ProductProcessor.GetAll(getProductsParams); err != nil {
		utils.BadRequestError(w, err)
	} else {
		utils.SendJSON(w, products, http.StatusOK)
	}

}
func (p *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {

	var body types.CreateProduct

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.BadRequestError(w, err)
		return
	}

	if err := utils.ValidateStruct(body); err != nil {
		utils.SendValidatonErrors(w, err)
		return
	}

	if err := p.ProductProcessor.Create(&body); err != nil {
		utils.BadRequestError(w, err)
		return
	}

	utils.SendJSON(w, "success", http.StatusCreated)
}
func (p *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.BadRequestError(w, err)
		return
	}

	var updateProduct types.UpdateProduct
	if err := json.NewDecoder(r.Body).Decode(&updateProduct); err != nil {
		utils.BadRequestError(w, err)
		return
	}

	if err := p.ProductProcessor.Update(id, &updateProduct); err != nil {
		utils.BadRequestError(w, err)
		return
	}

	utils.SendJSON(w, "success", http.StatusOK)
}
func (p *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.BadRequestError(w, err)
		return
	}

	if err := p.ProductProcessor.Delete(id); err != nil {
		utils.BadRequestError(w, err)
		return
	}

	utils.SendJSON(w, "success", http.StatusOK)
}

func (p *ProductHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(maxFileSize); err != nil {
		utils.BadRequestError(w, err)
		return
	}

	if len(r.MultipartForm.File[imageFieldName]) != 4 {
		utils.BadRequestError(w, errors.New("требуемое количество картинок-4"))
		return
	}

	imgUrls := make([]string, len(r.MultipartForm.File[imageFieldName]))
	mu, wg := sync.Mutex{}, sync.WaitGroup{}
	var err error

	addErr := func(e error) {
		mu.Lock()
		err = errors.Join(err, e)
		mu.Unlock()
	}

	for i, fileHeader := range r.MultipartForm.File[imageFieldName] {
		wg.Add(1)
		go func() {
			defer wg.Done()
			file, err := fileHeader.Open()
			if err != nil {
				addErr(err)
				return
			}
			defer file.Close()

			imgUrl, uploadErr := p.ProductProcessor.UploadImage(file)
			if uploadErr != nil {
				addErr(err)
				return
			}

			imgUrls[i] = imgUrl
		}()
	}

	wg.Wait()

	if err != nil {
		utils.InternalServerError(w, errors.New("что-то пошло не так"))
		return
	}

	utils.SendJSON(w, imgUrls, http.StatusOK)
}
func (p *ProductHandler) GetFilters(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")

	if category == "" || (category != "одежда" && category != "обувь") {
		utils.BadRequestError(w, errors.New("категория должно быть один из вариантов одежда,обувь"))
		return
	}

	if filters, err := p.ProductProcessor.GetFilters(category); err != nil {
		fmt.Println(err)
		utils.InternalServerError(w, errors.New("Что-то пошло не так.Повторите попытку чуть позже"))
	} else {
		utils.SendJSON(w, filters, http.StatusOK)
	}

}
func (p *ProductHandler) transformUrlParams(params map[string][]string) (*map[string]any, error) {
	m := make(map[string]any, len(params))

	for key, val := range params {
		switch key {
		case "page", "limit":
			num, err := strconv.Atoi(val[0])
			if err != nil {
				return nil, fmt.Errorf("%s должно быть целым числом", key)
			}
			m[key] = num
		case "colors", "size":
			m[key] = strings.Split(strings.Join(val, ","), ",")
		case "brand":
			splitedArr := strings.Split(strings.Join(val, ","), ",")
			slice := make([]int, 0, len(splitedArr))
			for _, value := range splitedArr {
				num, err := strconv.Atoi(value)
				if err != nil {
					fmt.Println(err)
					return nil, fmt.Errorf("%s должно быть массив целых чисел", key)
				}
				slice = append(slice, num)
			}

			m[key] = slice
		default:
			m[key] = val[0]
		}
	}

	return &m, nil
}
