package productrepo

import (
	"context"

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

	return id, err
}
