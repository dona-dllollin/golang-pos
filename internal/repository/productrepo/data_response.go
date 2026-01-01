package productrepo

type Image struct {
	ID        int `json:"id"`
	ProductID int64
	URL       string `json:"url"`
	SortOrder int    `json:"sort_order"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Product struct {
	ID          int
	Name        string
	Description string
	Status      string
	Images      []Image
	Categories  []Category
}
