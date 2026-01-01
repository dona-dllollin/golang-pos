package productrepo

import (
	"context"

	"github.com/dona-dllollin/belajar-clean-arch/internal/domain/productModel"
)

type ProductFilter struct {
	Keyword    string
	CategoryID *int64
	Status     string // active/inactive/archived
	Limit      int
	Offset     int
}

type ProductRepoInterface interface {
	// Create product lengkap (beserta image, variant, unit, option)
	Create(ctx context.Context, p *productModel.Product) (int64, error)

	// // Update product lengkap (deep update: images, variants, units)
	Update(ctx context.Context, p *productModel.Product) error

	// // Soft delete / archive product
	Delete(ctx context.Context, id int64) error

	// // Get product lengkap by id
	FindByID(ctx context.Context, id int64) (*productModel.ProductDetail, error)

	// List product dengan filter fleksibel
	FindAll(ctx context.Context, filter ProductFilter) ([]productModel.Product, error)

	// Get Image By ID
	GetImageById(ctx context.Context, id int64) (string, error)

	// // Cek stok varian tertentu
	// GetVariantStock(ctx context.Context, variantID int64) (int, error)

	// // Update stok varian (misal POS)
	// UpdateVariantStock(ctx context.Context, variantID int64, newStock int) error

}

type CategoryInterface interface {
	CreateCategory(ctx context.Context, c *productModel.Category) (int64, error)
	UpdateCategory(ctx context.Context, c *productModel.Category) error
	DeleteCategory(ctx context.Context, id int64) error
	FindAllCategory(ctx context.Context) ([]productModel.Category, error)
	FindCategory(ctx context.Context, id int64) (*productModel.Category, error)
	// FindCategoryByParent(ctx context.Context, id int64) (*productModel.Category, error)
}
