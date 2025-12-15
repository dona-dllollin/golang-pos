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
	ProductUseCase := productcase.NewProductService(*productRepository)
	imageService := imagecase.ImageUploadService{
		PublicPath:  imagePath,
		StoragePath: fmt.Sprintf("%s/%s", stotagePath, imagePath),
	}
	productHandler := NewProductHandler(ProductUseCase, validator, &imageService)

	r.Post("/", productHandler.StoreProduct)

}
