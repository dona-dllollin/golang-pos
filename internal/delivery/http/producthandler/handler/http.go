package handler

import (
	"net/http"
	"strconv"

	"github.com/dona-dllollin/belajar-clean-arch/internal/delivery/http/producthandler/dto"
	"github.com/dona-dllollin/belajar-clean-arch/internal/domain/productModel"
	"github.com/dona-dllollin/belajar-clean-arch/internal/usecase/imagecase"
	"github.com/dona-dllollin/belajar-clean-arch/internal/usecase/productcase"
	"github.com/dona-dllollin/belajar-clean-arch/pkgs/logger"
	"github.com/dona-dllollin/belajar-clean-arch/pkgs/validation"
	errorUtils "github.com/dona-dllollin/belajar-clean-arch/utils/errors"
	"github.com/dona-dllollin/belajar-clean-arch/utils/response"
)

type productHandler struct {
	productService productcase.ProductService
	imageService   imagecase.ImageService
	validator      validation.Validation
}

func NewProductHandler(productService productcase.ProductService, validator validation.Validation, imageService imagecase.ImageService) *productHandler {
	return &productHandler{
		productService: productService,
		imageService:   imageService,
		validator:      validator,
	}
}

func (h *productHandler) StoreProduct(w http.ResponseWriter, r *http.Request) {

	// wajib untuk multipart
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req dto.CreateProductRequest

	// string
	req.Name = r.FormValue("name")
	req.Description = r.FormValue("description")

	// category[] -> []*int64
	categoryValues := r.MultipartForm.Value["category_id"]
	for _, v := range categoryValues {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			http.Error(w, "invalid category_id", http.StatusBadRequest)
			return
		}
		req.CategoryId = append(req.CategoryId, &id)
	}

	// validate
	if err := h.validator.ValidateStruct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	product := productModel.Product{
		Name:        req.Name,
		Description: req.Description,
		CategoryId:  req.CategoryId,
	}

	//	panggil service product dulu supaya dapat id productnya
	id, err := h.productService.CreateProduct(r.Context(), &product)
	if err != nil {
		logger.Error("error:", err)
		errorUtils.WriteHTTPError(w, err)
		return
	}

	// Handle Image
	files := r.MultipartForm.File["images"]

	// panggil image service
	for i, file := range files {
		currentSortOrder := i + 1
		url, err := h.imageService.ImageUpload(r.Context(), file)
		if err != nil {
			logger.Error("gagal upload gambar", err.Error())
			errorUtils.WriteHTTPError(w, err)
			return
		}

		// add in product model
		product.Images = append(product.Images, productModel.ProductImage{
			ProductID: *id,
			URL:       url,
			SortOrder: currentSortOrder,
		})
	}
	response.JSON(w, http.StatusCreated, "success", id)

}
