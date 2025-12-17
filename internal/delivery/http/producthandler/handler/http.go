package handler

import (
	"encoding/json"
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
	"github.com/go-chi/chi/v5"
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

// CREATE PRODUCT
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

	// Handle Image
	files := r.MultipartForm.File["images"]

	// panggil image service
	if len(files) > 0 {
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
				URL:       url,
				SortOrder: currentSortOrder,
			})
		}
	}

	//	panggil service product
	id, err := h.productService.CreateProduct(r.Context(), &product)
	if err != nil {
		logger.Error("error:", err)
		errorUtils.WriteHTTPError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, "success", id)

}

// GET ALL CATEGORY
func (h *productHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.productService.ListCategories(r.Context())
	var res []dto.ListCategories
	for _, c := range categories {
		var category dto.ListCategories
		category.ID = c.ID
		category.Name = c.Name
		category.ParentID = c.ParentID
		res = append(res, category)
	}
	if err != nil {
		logger.Error("failed to get list category", err.Error())
		errorUtils.WriteHTTPError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, "success", res)
}

// CREATE CATEGORY
func (h *productHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateCategory

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.Error("failed to decode json body", err.Error())
		errorUtils.WriteHTTPError(w, err)
		return
	}

	category := productModel.Category{
		Name:     req.Name,
		ParentID: req.ParentID,
	}

	id, err := h.productService.CreateCategory(r.Context(), &category)
	if err != nil {
		logger.Error("failed to create category", err.Error())
		errorUtils.WriteHTTPError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, "success", id)

}

// Get Category BY ID
func (s productHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "categoryId")
	categoryId, err := strconv.Atoi(id)
	if err != nil {
		logger.Error(err.Error())
		errorUtils.WriteHTTPError(w, err)
		return
	}
	category, err := s.productService.GetCategory(r.Context(), int64(categoryId))
	if err != nil {
		logger.Error(err.Error())
		errorUtils.WriteHTTPError(w, errorUtils.ErrNotFound)
		return
	}

	response.JSON(w, http.StatusOK, "success", category)
}

// UPDATE CATEGORY
func (s productHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateCategory

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.Error("failed to decode json body", err.Error())
		errorUtils.WriteHTTPError(w, err)
		return
	}
	category := productModel.Category{
		ID:       req.ID,
		Name:     req.Name,
		ParentID: req.ParentID,
	}
	err = s.productService.UpdateCategory(r.Context(), &category)
	if err != nil {
		logger.Error("failed to update category", err.Error())
		errorUtils.WriteHTTPError(w, err)
		return
	}

	response.JSON(w, http.StatusNoContent, "success")

}

// DELETE CATEGORY
func (s productHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "categoryId")
	categoryId, err := strconv.Atoi(id)
	if err != nil {
		logger.Error(err.Error())
		errorUtils.WriteHTTPError(w, err)
		return
	}
	err = s.productService.DeleteCategory(r.Context(), int64(categoryId))
	if err != nil {
		logger.Error(err.Error())
		errorUtils.WriteHTTPError(w, err)
		return
	}

	response.JSON(w, http.StatusNoContent, "success")

}
