package productrepo

import (
	"context"
	"time"

	"github.com/dona-dllollin/belajar-clean-arch/internal/domain/productModel"
	"github.com/dona-dllollin/belajar-clean-arch/pkgs/logger"
	utils "github.com/dona-dllollin/belajar-clean-arch/utils/errors"
	"github.com/jackc/pgx/v5"
)

type ProductRepository struct {
	db *pgx.Conn
}

func NewProductRepository(db *pgx.Conn) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (conn ProductRepository) Create(ctx context.Context, p *productModel.Product) (int64, error) {

	// insert into table product
	query := `INSERT INTO products (name, description) VALUES ($1, $2) RETURNING id`
	var id int64
	err := conn.db.QueryRow(ctx, query, p.Name, p.Description).Scan(&id)
	if err != nil {
		logger.Error("Error: ", err.Error())
		return 0, utils.MapDbError(err)
	}

	// insert into table categoy_product
	for _, categoryId := range p.CategoryId {
		_, err := conn.db.Exec(ctx, `INSERT INTO category_products (product_id, category_id) VALUES ($1, $2)`, id, categoryId)
		if err != nil {
			logger.Error("Error: ", err.Error())
			return 0, utils.MapDbError(err)
		}
	}

	// insert into table product_images
	for _, image := range p.Images {
		_, err := conn.db.Exec(ctx, `INSERT INTO product_images (product_id, url, sort_order) VALUES ($1, $2, $3)`, id, image.URL, image.SortOrder)
		if err != nil {
			logger.Error("Error: ", err.Error())
			return 0, utils.MapDbError(err)
		}
	}

	return id, nil
}

type CategoryRepository struct {
	db *pgx.Conn
}

func NewCategoryRepsitory(db *pgx.Conn) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

func (conn CategoryRepository) CreateCategory(ctx context.Context, c *productModel.Category) (int64, error) {
	query := `INSERT INTO categories (name, parent_id) VALUES ($1, $2) RETURNING id`

	var id int64
	err := conn.db.QueryRow(ctx, query, c.Name, c.ParentID).Scan(&id)
	if err != nil {
		logger.Error("Error: ", err.Error())
		return 0, utils.MapDbError(err)
	}
	return id, nil
}

func (conn CategoryRepository) FindCategory(ctx context.Context, id int64) (*productModel.Category, error) {
	query := `SELECT id, name, parent_id FROM categories WHERE id = $1`
	var category productModel.Category
	err := conn.db.QueryRow(ctx, query, id).Scan(&category.ID, &category.Name, &category.ParentID)
	if err != nil {
		logger.Error("Error: ", err.Error())
		return nil, utils.MapDbError(err)
	}
	return &category, nil
}

func (conn CategoryRepository) UpdateCategory(ctx context.Context, c *productModel.Category) error {
	query := `UPDATE categories SET name = $2, parent_id = $3, updated_at = $4 WHERE id = $1`

	_, err := conn.db.Exec(ctx, query, c.ID, c.Name, c.ParentID, time.Now())

	if err != nil {
		logger.Error("Error: ", err.Error())
		return utils.MapDbError(err)
	}
	return nil
}

func (conn CategoryRepository) DeleteCategory(ctx context.Context, id int64) error {
	query := `DELETE FROM categories WHERE id = $1`

	_, err := conn.db.Exec(ctx, query, id)
	if err != nil {
		logger.Error("Error: ", err.Error())
		return utils.MapDbError(err)
	}
	return nil
}

func (conn CategoryRepository) FindAllCategory(ctx context.Context) ([]productModel.Category, error) {
	query := `SELECT id, name, parent_id FROM categories`
	rows, err := conn.db.Query(ctx, query)
	if err != nil {
		logger.Error("Error: ", err.Error())
		return nil, utils.MapDbError(err)
	}

	var categories []productModel.Category

	for rows.Next() {
		var category productModel.Category
		if err := rows.Scan(&category.ID, &category.Name, &category.ParentID); err != nil {
			logger.Error("Error: ", err.Error())
			return nil, utils.MapDbError(err)
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		logger.Error("Error: ", err.Error())
		return nil, utils.MapDbError(err)
	}
	return categories, nil
}
