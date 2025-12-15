package dto

type CreateProductRequest struct {
	Name        string `validate:"required"`
	Description string
	CategoryId  []*int64
	// ImageFiles  []*multipart.FileHeader // << ini untuk upload
	// Variant, Unit, dll bebas
}
