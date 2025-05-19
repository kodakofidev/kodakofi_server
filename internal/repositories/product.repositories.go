package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kodakofidev/kodakofi_server/internal/models"
)

type ProductRepoInterface interface {
	GetAllProducts(c context.Context, params *models.ProductQueryParams) (*models.PaginatedResponse, error)
	GetDetailProduct(c context.Context, id string) (*models.Product, error)
	GetRecommendation(c context.Context, limit int) (models.Products, error)
	AddProduct(c context.Context, newProduct *models.ProductRequest, images []string) error
}

type RepoProduct struct {
	DB *pgxpool.Pool
}

func NewProduct(db *pgxpool.Pool) *RepoProduct {
	return &RepoProduct{DB: db}
}

func (r *RepoProduct) GetAllProducts(c context.Context, params *models.ProductQueryParams) (*models.PaginatedResponse, error) {
	pageSize := 6
	offset := (params.Page - 1) * pageSize

	// Build WHERE clauses
	var whereClauses []string
	var args []interface{}
	argIndex := 1 // Parameter index dimulai dari 1 untuk CTE

	// Filter harga
	if params.Min > 0 || params.Max > 0 {
		minPrice := params.Min
		maxPrice := params.Max
		if maxPrice <= 0 {
			maxPrice = 1000000
		}
		whereClauses = append(whereClauses, fmt.Sprintf("p.price BETWEEN $%d AND $%d", argIndex, argIndex+1))
		args = append(args, minPrice, maxPrice)
		argIndex += 2
	}

	// Filter pencarian nama
	if params.Search != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("p.name ILIKE '%%' || $%d || '%%'", argIndex))
		args = append(args, params.Search)
		argIndex++
	}

	// Filter kategori
	if params.Category != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("c.name = $%d", argIndex))
		args = append(args, params.Category)
		argIndex++
	}

	// Filter diskon
	if params.Discount != "" {
		whereClauses = append(whereClauses, "d.discount IS NOT NULL")
	}

	// Filter newest
	if params.Options == "newest" {
		whereClauses = append(whereClauses, "p.created_at >= NOW()")
	}

	// Filter newest
	whereClauses = append(whereClauses, "p.is_deleted = false ")

	cteQuery := `
    WITH filtered_products AS (
        SELECT p.id
        FROM products p
        LEFT JOIN product_discounts pd ON pd.product_id = p.id 
        LEFT JOIN discounts d ON d.id = pd.discount_id 
        LEFT JOIN products_orders po ON po.product_id = p.id 
        LEFT JOIN orders o ON o.id = po.order_id 
        LEFT JOIN product_images pi ON pi.product_id = p.id 
        LEFT JOIN ratings r ON r.product_id = p.id 
        JOIN categories c ON c.id = p.category_id`

	// Tambahkan WHERE clause jika ada
	if len(whereClauses) > 0 {
		cteQuery += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	cteQuery += ` GROUP BY p.id, d.id, c.name
    ),
    product_stats AS (
        SELECT COUNT(*) as total_filtered FROM filtered_products
    )`

	// Build main query
	mainQuery := `
    SELECT 
        p.id, p.name, p.category_id, p.price, p.description,
        d.name AS discount_name, d.discount, 
        COALESCE(SUM(po.qty), 0) AS total_order, 
        COALESCE(json_agg(pi.path) FILTER (WHERE pi.path IS NOT NULL), '[]'::json) AS images, 
        COUNT(r.*) AS total_ratings,
        c.name AS category_name,
        ps.total_filtered
    FROM filtered_products fp
    JOIN products p ON p.id = fp.id
    LEFT JOIN product_discounts pd ON pd.product_id = p.id
    LEFT JOIN discounts d ON d.id = pd.discount_id
    LEFT JOIN products_orders po ON po.product_id = p.id
    LEFT JOIN product_images pi ON pi.product_id = p.id
    LEFT JOIN ratings r ON r.product_id = p.id
    JOIN categories c ON c.id = p.category_id
    CROSS JOIN product_stats ps
    GROUP BY p.id, d.id, c.name, ps.total_filtered`
	// Tambahkan ORDER BY
	switch params.Options {
	case "oldest":
		mainQuery += " ORDER BY p.created_at ASC"
	case "asc":
		mainQuery += " ORDER BY p.name ASC"
	case "desc":
		mainQuery += " ORDER BY p.name DESC"
	case "cheapest":
		mainQuery += " ORDER BY p.price ASC"
	case "favorite":
		mainQuery += " ORDER BY po.qty DESC"
	default:
		mainQuery += " ORDER BY p.created_at DESC"
	}

	// Tambahkan LIMIT dan OFFSET
	mainQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, pageSize, offset)

	// Gabungkan query lengkap
	fullQuery := cteQuery + mainQuery

	// Eksekusi query
	rows, err := r.DB.Query(c, fullQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	var products models.Products
	var totalFiltered int

	for rows.Next() {
		var product models.Product
		var imagesJSON []byte
		var discountName sql.NullString
		var discount sql.NullFloat64

		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.CategoryID,
			&product.Price,
			&product.Description,
			&discountName,
			&discount,
			&product.TotalOrder,
			&imagesJSON,
			&product.TotalRatings,
			&product.CategoryName,
			&totalFiltered,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		// Handle NULL values
		if discountName.Valid {
			product.DiscountName = &discountName.String
		}
		if discount.Valid {
			product.Discount = &discount.Float64
		}

		// Parse JSON images
		if err := json.Unmarshal(imagesJSON, &product.Images); err != nil {
			return nil, fmt.Errorf("error parsing images JSON: %w", err)
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after row iteration: %w", err)
	}

	// Hitung total pages
	totalPages := totalFiltered / pageSize
	if totalFiltered%pageSize > 0 {
		totalPages++
	}

	// Build pagination links
	basePath := "/api/product"
	links := map[string]string{
		"prev": "",
		"next": "",
	}

	if params.Page > 1 {
		links["prev"] = fmt.Sprintf("%s?page=%d", basePath, params.Page-1)
	}

	if params.Page < totalPages {
		links["next"] = fmt.Sprintf("%s?page=%d", basePath, params.Page+1)
	}

	response := &models.PaginatedResponse{
		Data: products,
		Pagination: models.Pagination{
			Page:       params.Page,
			PageSize:   pageSize,
			TotalItems: totalFiltered,
			TotalPages: totalPages,
			Links:      links,
		},
	}

	return response, nil
}

func (r *RepoProduct) GetDetailProduct(c context.Context, id string) (*models.Product, error) {
	query := `
		SELECT p.id, p.name, p.category_id, p.price, p.description,
			d.name AS discount_name, d.discount, 
			COALESCE((
			  SELECT SUM(po.qty)
			  FROM products_orders po
			  JOIN orders o ON o.id = po.order_id
			  WHERE po.product_id = p.id
			), 0) AS total_order,
  			(
  			  SELECT json_agg(pi.path)
  			  FROM product_images pi
  			  WHERE pi.product_id = p.id
  			) AS images,
  			(
  			  SELECT COUNT(*)
  			  FROM ratings r
  			  WHERE r.product_id = p.id AND r.rating = TRUE
  			) AS total_ratings
		FROM products p
		LEFT JOIN product_discounts pd ON pd.product_id = p.id
		LEFT JOIN discounts d ON d.id = pd.discount_id
		JOIN categories c ON c.id = p.category_id
		WHERE p.id = $1`

	var detail models.Product
	err := r.DB.QueryRow(c, query, id).Scan(
		&detail.ID,
		&detail.Name,
		&detail.CategoryID,
		&detail.Price,
		&detail.Description,
		&detail.DiscountName,
		&detail.Discount,
		&detail.TotalOrder,
		&detail.Images,
		&detail.TotalRatings,
	)

	if err != nil {
		return &models.Product{}, err
	}
	return &detail, nil
}

func (r *RepoProduct) GetRecommendation(c context.Context, limit int) (models.Products, error) {
	query := `
		SELECT p.id, p.name, p.price, d.name AS discount_name, d.discount,
			(
			  SELECT json_agg(pi.path)
			  FROM product_images pi
			  WHERE pi.product_id = p.id
			) AS images,
			(
			  SELECT COUNT(*)
			  FROM ratings r
			  WHERE r.product_id = p.id AND r.rating = TRUE
			) AS total_ratings
		FROM products p
		LEFT JOIN product_discounts pd ON pd.product_id = p.id
		LEFT JOIN discounts d ON d.id = pd.discount_id
		JOIN categories c ON c.id = p.category_id
		ORDER BY total_ratings DESC
		LIMIT $1
	`

	rows, err := r.DB.Query(c, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Price,
			&p.DiscountName,
			&p.Discount,
			&p.Images,
			&p.TotalRatings,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}
func (r *RepoProduct) AddProduct(c context.Context, newProduct *models.ProductRequest, listImage []string) error {
	tx, err := r.DB.Begin(c)
	if err != nil {
		log.Println("[DEBUG 1]", err)
		return err
	}

	func() {
		if err != nil {
			tx.Rollback(c)
		}
	}()
	query := `insert into products (name, category_id, price,description) values ($1, $2,$3,$4) returning id`
	var productID models.ProductRequest
	values := []any{newProduct.Name, newProduct.CategoryID, newProduct.Price, newProduct.Description}
	log.Println("[DEBUG LENGTH]", len(newProduct.Size))
	if len(newProduct.Size) == 0 {
		newProduct.Size = append(newProduct.Size, 4)
	}
	err = tx.QueryRow(c, query, values...).Scan(&productID.Id)
	log.Println("[DEBUG ID]", productID.Id)
	if err != nil {
		log.Println("[ERR ]", err)
	}
	log.Println("[ID]", productID)
	querySize := `insert into size_products (product_id, stock,size_id) values`
	valuesSize := []any{productID.Id, newProduct.Stock}
	log.Println("DEBUG SIZE", newProduct.Size)
	for i, size := range newProduct.Size {
		if i > 0 {
			querySize += ","
		}
		querySize += fmt.Sprintf("($1, $2,$%d)", i+3)
		valuesSize = append(valuesSize, size)
	}
	log.Println("[valuessize]", valuesSize)
	cmd, err := tx.Exec(c, querySize, valuesSize...)
	if err != nil {
		log.Println("[DEBUG 2]", err)
		return errors.New("add product failed")
	}
	row := cmd.RowsAffected()
	if row == 0 {
		return errors.New("add product failed ")
	}

	queryImage := `insert into product_images (product_id, path) values`
	valuesImage := []any{productID.Id}
	for i, image := range listImage {
		if i > 0 {
			queryImage += ","
		}
		queryImage += fmt.Sprintf("($1,$%d)", i+2)
		valuesImage = append(valuesImage, image)
	}
	log.Println("[DEBUG 2]", queryImage)

	cmd, err = tx.Exec(c, queryImage, valuesImage...)

	if err != nil {
		log.Println("[DEBUG 3]", err)

		return errors.New("add path image failed")
	}
	row = cmd.RowsAffected()
	if row == 0 {
		return errors.New("add path image failed")
	}
	if err := tx.Commit(c); err != nil {
		return err
	}

	return nil
}
