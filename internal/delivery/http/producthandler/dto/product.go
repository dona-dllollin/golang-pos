package dto

type CreateProductRequest struct {
	Name        string `validate:"required"`
	Description string
	CategoryId  []*int64
	// ImageFiles  []*multipart.FileHeader // << ini untuk upload
	// Variant, Unit, dll bebas
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
