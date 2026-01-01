package productcase

import (
	"context"

	"github.com/dona-dllollin/belajar-clean-arch/internal/domain/productModel"
	Repository "github.com/dona-dllollin/belajar-clean-arch/internal/repository/productrepo"
	"github.com/dona-dllollin/belajar-clean-arch/pkgs/logger"
)

type ProductService interface {
	// ------ PRODUCT ------
	CreateProduct(ctx context.Context, p *productModel.Product) (*int64, error)
	UpdateProduct(ctx context.Context, p *productModel.Product) error
	DeleteProduct(ctx context.Context, id int64) error
	GetProductByID(ctx context.Context, id int64) (*productModel.ProductDetail, error)
	ListProducts(ctx context.Context, filter ProductFilter) ([]productModel.Product, error)
	GetProductImage(ctx context.Context, id int64) (string, error)

	// ------ CATEGORY ------
	CreateCategory(ctx context.Context, c *productModel.Category) (*int64, error)
	UpdateCategory(ctx context.Context, c *productModel.Category) error
	DeleteCategory(ctx context.Context, id int64) error
	ListCategories(ctx context.Context) ([]productModel.Category, error)
	GetCategory(ctx context.Context, id int64) (*productModel.Category, error)
	// GetCategoryByParentId(ctx context.Context, id int64) ([]productModel.Category, error)

	// // ------ VARIANT ------
	// AddVariant(ctx context.Context, v *productModel.Variant) (int64, error)
	// UpdateVariant(ctx context.Context, v *productModel.Variant) error
	// DeleteVariant(ctx context.Context, id int64) error
	// GetVariantByID(ctx context.Context, id int64) (*productModel.Variant, error)

	// // ------ VARIANT OPTION ------
	// AddVariantOption(ctx context.Context, variantID int64, opt *productModel.VariantOption) error
	// UpdateVariantOption(ctx context.Context, variantID int64, opt *productModel.VariantOption) error
	// DeleteVariantOption(ctx context.Context, variantID int64, optionName string) error

	// // ------ VARIANT UNIT ------
	// AddVariantUnit(ctx context.Context, unit *productModel.VariantUnit) (int64, error)
	// UpdateVariantUnit(ctx context.Context, unit *productModel.VariantUnit) error
	// DeleteVariantUnit(ctx context.Context, unitID int64) error
}

type ProductFilter struct {
	Search     string
	CategoryID *int64
	Status     *string // active/inactive/archived
	Limit      int
	Offset     int
}

type ProductUseCase struct {
	productRepo  Repository.ProductRepoInterface
	categoryRepo Repository.CategoryInterface
}

func NewProductService(productRepo Repository.ProductRepoInterface, categoryRepo Repository.CategoryInterface) *ProductUseCase {
	return &ProductUseCase{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
	}
}

// ----------------------------------------------------------------------
// PRODUCT BASIC
// ----------------------------------------------------------------------

func (s *ProductUseCase) CreateProduct(ctx context.Context, p *productModel.Product) (*int64, error) {

	id, err := s.productRepo.Create(ctx, p)
	if err != nil {
		logger.Errorf("Create fail, error: %s", err)
		return nil, err
	}

	return &id, nil

}

func (s *ProductUseCase) UpdateProduct(ctx context.Context, p *productModel.Product) error {
	return s.productRepo.Update(ctx, p)
}

func (s *ProductUseCase) DeleteProduct(ctx context.Context, id int64) error {
	return s.productRepo.Delete(ctx, id)
}

func (s *ProductUseCase) GetProductByID(ctx context.Context, id int64) (*productModel.ProductDetail, error) {
	return s.productRepo.FindByID(ctx, id)
}

func (s *ProductUseCase) ListProducts(ctx context.Context, filter ProductFilter) ([]productModel.Product, error) {
	// Map struct usecase filter -> repo filter
	repoFilter := Repository.ProductFilter{
		Keyword:    filter.Search,
		CategoryID: filter.CategoryID,
		Limit:      filter.Limit,
		Offset:     filter.Offset,
	}

	if filter.Status != nil {
		repoFilter.Status = *filter.Status
	}

	products, err := s.productRepo.FindAll(ctx, repoFilter)
	if err != nil {
		logger.Errorf("ListProducts fail, error: %s", err)
		return nil, err
	}

	return products, nil
}

func (s *ProductUseCase) GetProductImage(ctx context.Context, id int64) (string, error) {
	return s.productRepo.GetImageById(ctx, id)
}

// ----------------------------------------------------------------------
// Category Product
// ----------------------------------------------------------------------

func (s *ProductUseCase) CreateCategory(ctx context.Context, c *productModel.Category) (*int64, error) {
	id, err := s.categoryRepo.CreateCategory(ctx, c)
	if err != nil {
		logger.Errorf("Create fail, error: %s", err)
		return nil, err
	}

	return &id, nil
}

func (s *ProductUseCase) UpdateCategory(ctx context.Context, c *productModel.Category) error {
	err := s.categoryRepo.UpdateCategory(ctx, c)
	if err != nil {
		logger.Errorf("Update fail, error: %s", err)
		return err
	}

	return nil
}

func (s *ProductUseCase) DeleteCategory(ctx context.Context, id int64) error {
	err := s.categoryRepo.DeleteCategory(ctx, id)
	if err != nil {
		logger.Errorf("Delete fail, error: %s", err)
		return err
	}

	return nil
}

func (s *ProductUseCase) ListCategories(ctx context.Context) ([]productModel.Category, error) {
	categories, err := s.categoryRepo.FindAllCategory(ctx)
	if err != nil {
		logger.Errorf("failed to Get List Categories, error: %s", err)
		return nil, err
	}

	return categories, nil
}

func (s *ProductUseCase) GetCategory(ctx context.Context, id int64) (*productModel.Category, error) {
	categories, err := s.categoryRepo.FindCategory(ctx, id)
	if err != nil {
		logger.Errorf("fail Get Category, error: %s", err)
		return nil, err
	}

	return categories, nil
}

// // ----------------------------------------------------------------------
// // PARTIAL EDIT: ADD VARIANT
// // ----------------------------------------------------------------------

// func (s *ProductUseCase) AddVariant(ctx context.Context, productID int64, v productModel.Variant) error {
// 	// 1. Load produk lengkap
// 	p, err := s.productRepo.FindByID(ctx, productID)
// 	if err != nil {
// 		return err
// 	}

// 	// 2. Tambah variant
// 	p.Variants = append(p.Variants, v)

// 	// 3. Simpan full product kembali
// 	return s.productRepo.Update(ctx, p)
// }

// // ----------------------------------------------------------------------
// // PARTIAL EDIT: ADD VARIANT OPTION
// // ----------------------------------------------------------------------

// func (s *ProductUseCase) AddVariantOption(ctx context.Context, variantID int64, opt productModel.VariantOption) error {
// 	// 1. Cari product yang memiliki variant itu
// 	// (biasanya lebih efisien FindByVariantID, tapi kita pakai FindByID full)
// 	p, err := s.findProductByVariantID(ctx, variantID)
// 	if err != nil {
// 		return err
// 	}

// 	// 2. Insert option pada variant yang tepat
// 	for i := range p.Variants {
// 		if p.Variants[i].ID == variantID {
// 			p.Variants[i].Options = append(p.Variants[i].Options, opt)
// 			break
// 		}
// 	}

// 	// 3. Save sebagai full update
// 	return s.productRepo.Update(ctx, p)
// }

// // ----------------------------------------------------------------------
// // PARTIAL EDIT: ADD VARIANT UNIT
// // ----------------------------------------------------------------------

// func (s *ProductUseCase) AddVariantUnit(ctx context.Context, variantID int64, u productModel.VariantUnit) error {
// 	p, err := s.findProductByVariantID(ctx, variantID)
// 	if err != nil {
// 		return err
// 	}

// 	for i := range p.Variants {
// 		if p.Variants[i].ID == variantID {
// 			p.Variants[i].Units = append(p.Variants[i].Units, u)
// 			break
// 		}
// 	}

// 	return s.productRepo.Update(ctx, p)
// }

// // ----------------------------------------------------------------------
// // Helper untuk mencari product berdasarkan variantID
// // ----------------------------------------------------------------------

// func (s *ProductUseCase) findProductByVariantID(ctx context.Context, variantID int64) (*productModel.Product, error) {
// 	// Cara paling aman: load product by variantID via repo, tapi
// 	// jika tidak ada, maka harus scan setelah FindAll/FindByID (kurang efisien)
// 	// Untuk contoh, kita pakai FindByID khusus (misal kamu buat di repo)

// 	return s.productRepo.FindByVariantID(ctx, variantID)
// }
