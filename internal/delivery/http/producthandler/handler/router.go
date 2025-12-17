package handler

import (
	"fmt"

	"github.com/dona-dllollin/belajar-clean-arch/internal/repository/productrepo"
	"github.com/dona-dllollin/belajar-clean-arch/internal/usecase/imagecase"
	"github.com/dona-dllollin/belajar-clean-arch/internal/usecase/productcase"
	"github.com/dona-dllollin/belajar-clean-arch/pkgs/validation"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

func Routes(
	r chi.Router,
	db *pgx.Conn,
	validator validation.Validation,
	imagePath string,
	stotagePath string,
) {

	productRepository := productrepo.NewProductRepository(db)
	categoryRepository := productrepo.NewCategoryRepsitory(db)
	ProductUseCase := productcase.NewProductService(productRepository, categoryRepository)
	imageService := imagecase.ImageUploadService{
		PublicPath:  imagePath,
		StoragePath: fmt.Sprintf("%s/%s", stotagePath, imagePath),
	}
	productHandler := NewProductHandler(ProductUseCase, validator, &imageService)

	// product
	r.Post("/", productHandler.StoreProduct)

	//categories
	r.Get("/categories", productHandler.ListCategories)
	r.Post("/category", productHandler.CreateCategory)
	r.Get("/category/{categoryId}", productHandler.GetCategory)
	r.Put("/category", productHandler.UpdateCategory)
	r.Delete("/category/{categoryId}", productHandler.DeleteCategory)

}
