package productrepo

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/dona-dllollin/belajar-clean-arch/internal/domain/productModel"
	"github.com/dona-dllollin/belajar-clean-arch/pkgs/logger"
	utils "github.com/dona-dllollin/belajar-clean-arch/utils/errors"
	"github.com/jackc/pgx/v5"
)

// ===========================================
// Product Repository
// ===========================================

type ProductRepository struct {
	db *pgx.Conn
}

func NewProductRepository(db *pgx.Conn) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

// ********** Implementation Create Product**********
func (conn ProductRepository) Create(ctx context.Context, p *productModel.Product) (int64, error) {

	tx, err := conn.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	// insert into table product
	var productID int64
	err = tx.QueryRow(ctx,
		`INSERT INTO products (name, description)
         VALUES ($1, $2) RETURNING id`,
		p.Name, p.Description,
	).Scan(&productID)
	if err != nil {
		return 0, utils.MapDbError(err)
	}

	// insert into table categoy_product
	if len(p.CategoryId) > 0 {
		batch := &pgx.Batch{}
		for _, cid := range p.CategoryId {
			batch.Queue(
				"INSERT INTO category_products (product_id, category_id) VALUES ($1, $2)",
				productID, cid,
			)
		}

		br := tx.SendBatch(ctx, batch)
		for range p.CategoryId {
			if _, err := br.Exec(); err != nil {
				br.Close()
				return 0, err
			}
		}
		br.Close()
	}

	// insert into table product_images
	for _, img := range p.Images {
		_, err = tx.Exec(ctx,
			"INSERT INTO product_images (product_id, url, sort_order) VALUES ($1, $2, $3)",
			productID, img.URL, img.SortOrder,
		)
		if err != nil {
			return 0, err
		}
	}
	// insert variants
	for _, v := range p.Variants {
		if err := conn.AddVariant(ctx, v, tx, productID); err != nil {
			return 0, err
		}
	}

	return productID, tx.Commit(ctx)
}

// ********** Implementation Add Variant Product**********
func (conn ProductRepository) AddVariant(ctx context.Context, variant productModel.Variant, tx pgx.Tx, id int64) error {

	var variant_id int64
	err := tx.QueryRow(ctx,
		"INSERT INTO variants (product_id, sku, base_unit, stock, cost_price) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		id,
		variant.SKU,
		variant.BaseUnit,
		variant.Stock,
		variant.CostPrice).
		Scan(&variant_id)
	if err != nil {
		logger.Error("Error: ", err.Error())
		return err
	}

	for _, vOption := range variant.Options {
		_, err = tx.Exec(ctx,
			"INSERT INTO variant_options (variant_id, name, value) VALUES ($1, $2, $3)",
			variant_id,
			vOption.Name,
			vOption.Value)
		if err != nil {
			logger.Error("Error: ", err.Error())
			return err
		}

	}

	for _, vUnit := range variant.Units {
		_, err = tx.Exec(ctx,
			"INSERT INTO variant_units (variant_id, name, barcode, conversion_rate, price) VALUES ($1, $2, $3, $4, $5)",
			variant_id,
			vUnit.Name,
			vUnit.Barcode,
			vUnit.ConversionRate,
			vUnit.Price)
		if err != nil {
			logger.Error("Error: ", err.Error())
			return err
		}

	}

	return nil

}

// ********** Implementation FindAll Product**********
func (conn ProductRepository) FindAll(ctx context.Context, filter ProductFilter) ([]productModel.Product, error) {
	query := `SELECT 
			p.id,
			p.name, 
			p.description, 
			p.status,
			COALESCE(
			JSONB_AGG(
			DISTINCT JSONB_BUILD_OBJECT(
			'id', pi.id,
			'productId', pi.product_id,
			'url', pi.url, 
			'sortOrder', pi.sort_order
			)
			) FILTER (WHERE pi.id IS NOT NULL),
			 '[]'::jsonb
			) AS images,

			COALESCE(
			JSONB_AGG(DISTINCT pc.category_id) 
			FILTER (WHERE pc.category_id IS NOT NULL),
			 '[]'::jsonb
			) AS categories
			FROM products p 
			LEFT JOIN product_images pi
				ON pi.product_id = p.id
			LEFT JOIN category_products pc
				ON pc.product_id = p.id`

	var args []interface{}
	var conditions []string

	if filter.CategoryID != nil {
		// query += " JOIN category_products cp ON cp.product_id = p.id"
		conditions = append(conditions, fmt.Sprintf("pc.category_id = $%d", len(args)+1))
		args = append(args, *filter.CategoryID)
	}

	if filter.Keyword != "" {
		conditions = append(conditions, fmt.Sprintf("p.name ILIKE $%d", len(args)+1))
		args = append(args, "%"+filter.Keyword+"%")
	}

	if filter.Status != "" {
		conditions = append(conditions, fmt.Sprintf("p.status = $%d", len(args)+1))
		args = append(args, filter.Status)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " GROUP BY p.id, p.name, p.description ORDER BY p.id DESC"

	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", len(args)+1)
		args = append(args, filter.Limit)
	}

	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", len(args)+1)
		args = append(args, filter.Offset)
	}

	rows, err := conn.db.Query(ctx, query, args...)
	if err != nil {
		logger.Error("Error: ", err.Error())
		return nil, utils.MapDbError(err)
	}
	defer rows.Close()

	var (
		imagesJSON     []byte
		categoriesJSON []byte
	)
	// var products []Product
	var productResponse []productModel.Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Status, &imagesJSON, &categoriesJSON); err != nil {
			logger.Error("Error: ", err.Error())
			return nil, utils.MapDbError(err)
		}
		// products = append(products, p)

		var (
			images     []productModel.ProductImage
			categories []*int64
		)
		json.Unmarshal(imagesJSON, &images)
		err := json.Unmarshal(categoriesJSON, &categories)
		if err != nil {
			logger.Error(err.Error())
		}

		productResponse = append(productResponse, productModel.Product{
			ID:          int64(p.ID),
			Name:        p.Name,
			Description: p.Description,
			Status:      p.Status,
			Images:      images,
			CategoryId:  categories,
		})
	}

	if err := rows.Err(); err != nil {
		logger.Error("Error: ", err.Error())
		return nil, utils.MapDbError(err)
	}

	return productResponse, nil
}

// ********** Implementation FindByID Product**********
func (conn ProductRepository) FindByID(ctx context.Context, id int64) (*productModel.ProductDetail, error) {
	var p productModel.ProductDetail
	query := `SELECT id, name, description, status FROM products WHERE id = $1`
	err := conn.db.QueryRow(ctx, query, id).Scan(&p.ID, &p.Name, &p.Description, &p.Status)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, utils.ErrNotFound
		}
		return nil, utils.MapDbError(err)
	}

	// Get Categories
	catQuery := `SELECT category_id, name FROM category_products WHERE product_id = $1`
	rows, err := conn.db.Query(ctx, catQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var category productModel.Category
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			return nil, err
		}
		p.Categories = append(p.Categories, category)
	}

	// Get Images
	imgQuery := `SELECT id, url, sort_order FROM product_images WHERE product_id = $1 ORDER BY sort_order ASC`
	imgRows, err := conn.db.Query(ctx, imgQuery, id)
	if err != nil {
		return nil, err
	}
	defer imgRows.Close()
	for imgRows.Next() {
		var img productModel.ProductImage
		img.ProductID = p.ID
		if err := imgRows.Scan(&img.ID, &img.URL, &img.SortOrder); err != nil {
			return nil, err
		}
		p.Images = append(p.Images, img)
	}

	// Get Variants
	varQuery := `SELECT id, sku, base_unit, stock, cost_price FROM variants WHERE product_id = $1`
	varRows, err := conn.db.Query(ctx, varQuery, id)
	if err != nil {
		return nil, err
	}
	defer varRows.Close()

	// Temporarily store variants to fetch their options/units
	for varRows.Next() {
		var v productModel.Variant
		v.ProductID = p.ID
		if err := varRows.Scan(&v.ID, &v.SKU, &v.BaseUnit, &v.Stock, &v.CostPrice); err != nil {
			return nil, err
		}
		p.Variants = append(p.Variants, v)
	}
	varRows.Close()

	// Fill Variant Options and Units
	for i := range p.Variants {
		// Options
		optQuery := `SELECT name, value FROM variant_options WHERE variant_id = $1`
		optRows, err := conn.db.Query(ctx, optQuery, p.Variants[i].ID)
		if err != nil {
			return nil, err
		}
		for optRows.Next() {
			var opt productModel.VariantOption
			if err := optRows.Scan(&opt.Name, &opt.Value); err != nil {
				optRows.Close()
				return nil, err
			}
			p.Variants[i].Options = append(p.Variants[i].Options, opt)
		}
		optRows.Close()

		// Units
		unitQuery := `SELECT id, name, barcode, conversion_rate, price FROM variant_units WHERE variant_id = $1`
		unitRows, err := conn.db.Query(ctx, unitQuery, p.Variants[i].ID)
		if err != nil {
			return nil, err
		}
		for unitRows.Next() {
			var u productModel.VariantUnit
			u.VariantID = p.Variants[i].ID
			if err := unitRows.Scan(&u.ID, &u.Name, &u.Barcode, &u.ConversionRate, &u.Price); err != nil {
				unitRows.Close()
				return nil, err
			}
			p.Variants[i].Units = append(p.Variants[i].Units, u)
		}
		unitRows.Close()
	}

	return &p, nil
}

// ********** Implementation Update Product**********
func (conn ProductRepository) Update(ctx context.Context, p *productModel.Product) error {

	tx, err := conn.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Update base product
	_, err = tx.Exec(ctx, `UPDATE products SET name=$1, description=$2, updated_at=NOW() WHERE id=$3`,
		p.Name, p.Description, p.ID)
	if err != nil {
		logger.Error("Error: ", err.Error())
		return utils.MapDbError(err)
	}

	// Update Categories: Delete all and re-insert
	_, err = tx.Exec(ctx,
		`INSERT INTO category_products (product_id, category_id) 
		 SELECT $1, unnest($2::bigint[]) EXCEPT SELECT product_id, category_id FROM category_products
		 WHERE product_id = $1`, p.ID, p.CategoryId)
	if err != nil {
		logger.Error("Error: ", err.Error())
		return utils.MapDbError(err)
	}
	_, err = tx.Exec(ctx,
		`DELETE FROM category_products 
		 WHERE product_id = $1 AND category_id NOT IN (SELECT unnest($2::bigint[]))`,
		p.ID, p.CategoryId)
	if err != nil {
		logger.Error("Error: ", err.Error())
		return utils.MapDbError(err)
	}

	// Update Images
	if len(p.Images) > 0 {
		err = conn.UpdateImage(ctx, p.ID, p.Images, tx)
		if err != nil {
			logger.Error("Error: ", err.Error())
			return utils.MapDbError(err)
		}
	}

	return tx.Commit(ctx)
}

// ********** Implementation Update Image**********
func (conn ProductRepository) UpdateImage(ctx context.Context, productId int64, images []productModel.ProductImage, tx pgx.Tx) error {

	oldImages, err := conn.GetImageByProductId(ctx, productId)
	if err != nil {
		logger.Error("Error: ", err.Error())
		return utils.MapDbError(err)
	}

	oldMap := make(map[int64]productModel.ProductImage)
	for _, img := range oldImages {
		oldMap[img.ID] = img
	}
	newMap := make(map[int64]productModel.ProductImage)
	for _, img := range images {
		newMap[img.ID] = img
	}

	for _, img := range oldMap {
		if _, ok := newMap[img.ID]; !ok {
			_, err = tx.Exec(ctx, `DELETE FROM product_images WHERE id=$1`, img.ID)
			if err != nil {
				logger.Error("Error: ", err.Error())
				return utils.MapDbError(err)
			}
		}
	}
	for _, img := range newMap {
		if _, ok := oldMap[img.ID]; !ok {
			_, err = tx.Exec(ctx, `INSERT INTO product_images (product_id, url, sort_order) VALUES ($1, $2, $3)`, productId, img.URL, img.SortOrder)
			if err != nil {
				logger.Error("Error: ", err.Error())
				return utils.MapDbError(err)
			}
		} else {
			_, err = tx.Exec(ctx, `UPDATE product_images SET url=$1, sort_order=$2 WHERE id=$3`, img.URL, img.SortOrder, img.ID)
			if err != nil {
				logger.Error("Error: ", err.Error())
				return utils.MapDbError(err)
			}
		}
	}
	return nil
}

// ********** Implementation Get Image By Product ID**********
func (conn ProductRepository) GetImageByProductId(ctx context.Context, productId int64) ([]productModel.ProductImage, error) {
	query := `SELECT id, url, sort_order FROM product_images WHERE product_id = $1`
	rows, err := conn.db.Query(ctx, query, productId)
	if err != nil {
		return nil, utils.MapDbError(err)
	}
	defer rows.Close()
	var images []productModel.ProductImage
	for rows.Next() {
		var img productModel.ProductImage
		if err := rows.Scan(&img.ID, &img.URL, &img.SortOrder); err != nil {
			return nil, utils.MapDbError(err)
		}
		images = append(images, img)
	}
	return images, nil
}

// ***** Implementation Get Image By ID ******
func (conn ProductRepository) GetImageById(ctx context.Context, id int64) (string, error) {
	query := `SELECT url FROM product_images WHERE id = $1`
	var url string
	err := conn.db.QueryRow(ctx, query, id).Scan(&url)
	if err != nil {
		return url, utils.MapDbError(err)
	}
	return url, nil
}

// ********** Implementation Delete Product**********
func (conn ProductRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE products SET status = 'archived', updated_at = NOW() WHERE id = $1`
	_, err := conn.db.Exec(ctx, query, id)
	if err != nil {
		return utils.MapDbError(err)
	}
	return nil
}

// ===========================================
// Category Repository
// ===========================================
type CategoryRepository struct {
	db *pgx.Conn
}

func NewCategoryRepsitory(db *pgx.Conn) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

// ********** Implementation Create Category**********
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

// ********** Implementation Get Category By Id**********
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

// ********** Implementation Update Category**********
func (conn CategoryRepository) UpdateCategory(ctx context.Context, c *productModel.Category) error {
	query := `UPDATE categories SET name = $2, parent_id = $3, updated_at = $4 WHERE id = $1`

	_, err := conn.db.Exec(ctx, query, c.ID, c.Name, c.ParentID, time.Now())

	if err != nil {
		logger.Error("Error: ", err.Error())
		return utils.MapDbError(err)
	}
	return nil
}

// ********** Implementation Delete Category**********
func (conn CategoryRepository) DeleteCategory(ctx context.Context, id int64) error {
	query := `DELETE FROM categories WHERE id = $1`

	_, err := conn.db.Exec(ctx, query, id)
	if err != nil {
		logger.Error("Error: ", err.Error())
		return utils.MapDbError(err)
	}
	return nil
}

// ********** Implementation Get list Category**********
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
