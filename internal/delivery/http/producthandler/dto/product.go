package dto

import "github.com/dona-dllollin/belajar-clean-arch/internal/domain/productModel"

type CreateProductRequest struct {
	Name        string `validate:"required"`
	Description string
	CategoryId  []*int64
	// ImageFiles  []*multipart.FileHeader // << ini untuk upload
	Variants []VariantRequest `json:"variants"`
}

type VariantRequest struct {
	SKU       string                 `json:"sku"`
	Options   []VariantOptionRequest `json:"options"`
	BaseUnit  string                 `json:"base_unit"`
	Stock     int                    `json:"stock"`
	CostPrice int64                  `json:"cost_price"`
	Units     []VariantUnitRequest   `json:"units"`
}

type VariantOptionRequest struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type VariantUnitRequest struct {
	ID             int64   // Just to be compatible to productModel
	VariantID      int64   // Same
	Name           string  `json:"name"`
	SKU            *string `json:"sku"`
	Barcode        *string `json:"barcode"`
	ConversionRate int     `json:"conversion_rate"`
	Price          int64   `json:"price"`
}

type CreateCategory struct {
	Name     string `json:"name" validate:"required,unique"`
	ParentID *int64 `json:"parent_id,omitempty"`
}

type UpdateCategory struct {
	ID       int64  `json:"id" validate:"required"`
	Name     string `json:"name" validate:"required,unique"`
	ParentID *int64 `json:"parent_id,omitempty"`
}

type ListCategories struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	ParentID *int64 `json:"parent_id"`
}

func MapOptions(options []VariantOptionRequest) []productModel.VariantOption {
	productOptions := []productModel.VariantOption{}
	for _, option := range options {
		productOptions = append(productOptions, productModel.VariantOption(option))
	}

	return productOptions

}

func MapUnits(units []VariantUnitRequest) []productModel.VariantUnit {
	productUnits := []productModel.VariantUnit{}
	for _, unit := range units {
		productUnits = append(productUnits, productModel.VariantUnit(unit))
	}

	return productUnits
}
