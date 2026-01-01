package productModel

//Product
type Product struct {
	ID          int64
	Name        string
	Description string
	CategoryId  []*int64
	Status      string // "active", "inactive", "archived"
	Images      []ProductImage
	Variants    []Variant
}

//Product
type ProductDetail struct {
	ID          int64
	Name        string
	Description string
	Categories  []Category
	Status      string // "active", "inactive", "archived"
	Images      []ProductImage
	Variants    []Variant
}

//Product Category
type Category struct {
	ID       int64
	Name     string
	ParentID *int64
}

//Product Image
type ProductImage struct {
	ID        int64
	ProductID int64
	URL       string
	SortOrder int
}

// Variant (kombinasi warna, ukuran, dll)
type Variant struct {
	ID        int64
	ProductID int64
	SKU       string
	Options   []VariantOption
	BaseUnit  string
	Stock     int // Stok dalam satan dasar
	CostPrice int64
	Units     []VariantUnit
}

// Varian Option (warna = merah, ukuran = M)
type VariantOption struct {
	Name  string
	Value string
}

// Variant Unit (multi-satuan POS)
type VariantUnit struct {
	ID             int64
	VariantID      int64
	Name           string //pcs, pack, dus, etc
	SKU            *string
	Barcode        *string
	ConversionRate int   // pack = 5 pcs -> 5
	Price          int64 //harga per unit
}
