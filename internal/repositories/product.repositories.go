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
	GetListImageProduct(c context.Context, id string) ([]string, error)
	DeleteImage(c context.Context, product_id string) error
	UpdateProduct(ctx context.Context, productID string, updateData *models.ProductRequest, newImages []string, shouldUpdateImages bool, currentImages []string) error
	ToggleLike(c context.Context, userID, productID string) (bool, error)
}

type RepoProduct struct {
	DB *pgxpool.Pool
}

func NewProduct(db *pgxpool.Pool) *RepoProduct {
	return &RepoProduct{DB: db}
}

func (r *RepoProduct) GetAllProducts(c context.Context, params *models.ProductQueryParams) (*models.PaginatedResponse, error) {
	pageSize := 6
	if params.Page < 1 {
		params.Page = 1
	}
	offset := (params.Page - 1) * pageSize

	var whereClauses []string
	var args []interface{}
	argIndex := 1

	if params.Min >= 0 || params.Max >= 0 {
		minPrice := params.Min
		if minPrice < 0 {
			minPrice = 0
		}
		maxPrice := params.Max
		if maxPrice <= 0 {
			maxPrice = 1000000
		}
		whereClauses = append(whereClauses, fmt.Sprintf("p.price BETWEEN $%d AND $%d", argIndex, argIndex+1))
		args = append(args, minPrice, maxPrice)
		argIndex += 2
	}

	if params.Search != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("p.name ILIKE '%%' || $%d || '%%'", argIndex))
		args = append(args, params.Search)
		argIndex++
	}

	if params.Category != nil {
		var categoryConditions []string
		for _, category := range params.Category {
			categoryConditions = append(categoryConditions, fmt.Sprintf("c.name = $%d", argIndex))
			args = append(args, category)
			argIndex++
		}
		whereClauses = append(whereClauses, "("+strings.Join(categoryConditions, " OR ")+")")
	}

	if params.Discount != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("d.id = $%d", argIndex))
		args = append(args, params.Discount)
		argIndex++
	}

	whereClauses = append(whereClauses, "p.is_deleted = false")

	cteQuery := `
	WITH filtered_products AS (
		SELECT DISTINCT p.id
		FROM products p
		LEFT JOIN product_discounts pd ON pd.product_id = p.id 
		LEFT JOIN discounts d ON d.id = pd.discount_id 
		JOIN categories c ON c.id = p.category_id`

	if len(whereClauses) > 0 {
		cteQuery += " WHERE " + strings.Join(whereClauses, " AND ")
	}
	cteQuery += `
	),
	product_stats AS (
		SELECT COUNT(*) as total_filtered FROM filtered_products
	)`

	mainQuery := `
	SELECT 
		p.id, p.name, p.category_id, p.price, p.description,
		d.name AS discount_name, d.discount, 
		COALESCE(SUM(po.qty), 0) AS total_order, 
		COALESCE(json_agg(DISTINCT jsonb_build_object('id', s.id, 'size', s.size, 'stock', sp.stock)) FILTER (WHERE s.id IS NOT NULL), '[]') AS size,
		COALESCE(json_agg(DISTINCT pi.path) FILTER (WHERE pi.path IS NOT NULL), '[]') AS images, 
		COUNT(DISTINCT r.*) AS total_ratings,
		c.name AS category_name,
		ps.total_filtered
	FROM filtered_products fp
	JOIN products p ON p.id = fp.id
	LEFT JOIN product_discounts pd ON pd.product_id = p.id
	LEFT JOIN size_products sp ON sp.product_id = p.id
	LEFT JOIN sizes s ON s.id = sp.size_id 
	LEFT JOIN discounts d ON d.id = pd.discount_id
	LEFT JOIN products_orders po ON po.product_id = p.id
	LEFT JOIN product_images pi ON pi.product_id = p.id
	LEFT JOIN ratings r ON r.product_id = p.id
	JOIN categories c ON c.id = p.category_id
	CROSS JOIN product_stats ps
	GROUP BY p.id, d.id, c.name, ps.total_filtered`

	switch params.Options {
	case "oldest":
		mainQuery += " ORDER BY p.created_at ASC"
	case "asc":
		mainQuery += " ORDER BY p.name ASC"
	case "desc":
		mainQuery += " ORDER BY p.name DESC"
	case "cheapest":
		mainQuery += " ORDER BY p.price ASC"
	case "bestseller":
		mainQuery += " ORDER BY total_order DESC"
	case "ratings":
		mainQuery += " ORDER BY total_ratings DESC"
	default:
		mainQuery += " ORDER BY p.created_at DESC"
	}
	log.Println()
	mainQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, pageSize, offset)

	fullQuery := cteQuery + mainQuery

	log.Println("[debug query]", fullQuery)

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
		var sizesJSON []byte
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
			&sizesJSON,
			&imagesJSON,
			&product.TotalRatings,
			&product.CategoryName,
			&totalFiltered,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		if discountName.Valid {
			product.DiscountName = &discountName.String
		}
		if discount.Valid {
			product.Discount = &discount.Float64
		}
		if err := json.Unmarshal(imagesJSON, &product.Images); err != nil {
			return nil, fmt.Errorf("error parsing images JSON: %w", err)
		}
		if err := json.Unmarshal(sizesJSON, &product.Sizes); err != nil {
			return nil, fmt.Errorf("error parsing sizes JSON: %w", err)
		}
		products = append(products, product)
	}

	totalPages := totalFiltered / pageSize
	if totalFiltered%pageSize > 0 {
		totalPages++
	}

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
  			  WHERE r.product_id = p.id
  			) AS total_ratings,
  			(
  			  SELECT json_agg(json_build_object(
  			    'id', s.id,
  			    'size', s.size,
  			    'stock', sp.stock
  			  ))
  			  FROM size_products sp
  			  JOIN sizes s ON s.id = sp.size_id
  			  WHERE sp.product_id = p.id
  			) AS sizes
		FROM products p
		LEFT JOIN product_discounts pd ON pd.product_id = p.id
		LEFT JOIN discounts d ON d.id = pd.discount_id
		JOIN categories c ON c.id = p.category_id
		WHERE p.id = $1`

	var detail models.Product
	var sizesJson []byte
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
		&sizesJson,
	)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(sizesJson, &detail.Sizes); err != nil {
		return nil, err
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
			  WHERE r.product_id = p.id AND r.rating > 0
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

func (r *RepoProduct) DeleteImage(c context.Context, product_id string) error {
	query := ` delete from product_images where product_id = $1`
	_, err := r.DB.Exec(c, query, product_id)
	if err != nil {
		return err
	}
	return nil
}

func (r *RepoProduct) GetListImageProduct(c context.Context, id string) ([]string, error) {
	log.Println("[DEBUG ID]", id)
	query := ` select json_agg(path) from product_images where product_id =$1`
	var listImage []string
	err := r.DB.QueryRow(c, query, id).Scan(&listImage)
	if err != nil {
		return nil, err
	}
	return listImage, nil
}

func (r *RepoProduct) UpdateProduct(ctx context.Context, productID string, updateData *models.ProductRequest, newImages []string, shouldUpdateImages bool, currentImages []string,
) error {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Build dynamic update query untuk product info
	var queryParts []string
	var params []interface{}
	paramCount := 1

	if updateData.Name != nil {
		queryParts = append(queryParts, fmt.Sprintf("name = $%d", paramCount))
		params = append(params, updateData.Name)
		paramCount++
	}
	if updateData.CategoryID != nil {
		queryParts = append(queryParts, fmt.Sprintf("category_id = $%d", paramCount))
		params = append(params, *updateData.CategoryID)
		paramCount++
	}
	if updateData.Description != nil {
		queryParts = append(queryParts, fmt.Sprintf("description = $%d", paramCount))
		params = append(params, *&updateData.Description)
		paramCount++
	}
	if updateData.Price != nil {
		queryParts = append(queryParts, fmt.Sprintf("price = $%d", paramCount))
		params = append(params, *updateData.Price)
		paramCount++
	}

	if len(queryParts) > 0 {
		query := fmt.Sprintf(`UPDATE products SET %s, updated_at = NOW() WHERE id = $%d`, strings.Join(queryParts, ", "), paramCount)

		params = append(params, productID)

		_, err := tx.Exec(ctx, query, params...)
		if err != nil {
			return fmt.Errorf("failed to update product: %w", err)
		}
	}

	// Handle size updates
	if updateData.Size != nil {
		_, err = tx.Exec(ctx, "DELETE FROM size_products WHERE product_id = $1", productID)
		if err != nil {
			return fmt.Errorf("failed to delete sizes: %w", err)
		}

		if len(updateData.Size) > 0 {
			querySize := `insert into size_products (product_id, stock,size_id) values`
			valuesSize := []any{productID, updateData.Stock}
			for i, size := range updateData.Size {
				if i > 0 {
					querySize += ","
				}
				querySize += fmt.Sprintf("($1, $2,$%d)", i+3)
				valuesSize = append(valuesSize, size)
			}
			log.Println("[valuessize]", valuesSize)
			cmd, err := tx.Exec(ctx, querySize, valuesSize...)
			if err != nil {
				log.Println("[DEBUG 2]", err)
				return errors.New("add product failed")
			}
			row := cmd.RowsAffected()
			if row == 0 {
				return errors.New("add product failed ")
			}

		}
	}
	_, err = tx.Exec(ctx, "DELETE FROM product_images WHERE product_id = $1", productID)

	queryImage := `insert into product_images (product_id, path) values`
	valuesImage := []any{productID}
	for i, image := range newImages {
		if i > 0 {
			queryImage += ","
		}
		queryImage += fmt.Sprintf("($1,$%d)", i+2)
		valuesImage = append(valuesImage, image)
	}
	log.Println("[DEBUG 2]", queryImage)

	cmd, err := tx.Exec(ctx, queryImage, valuesImage...)

	if err != nil {
		log.Println("[DEBUG 3]", err)

		return errors.New("add path image failed")
	}
	row := cmd.RowsAffected()
	if row == 0 {
		return errors.New("add path image failed")
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *RepoProduct) GetLikeStatus(c context.Context, userID, productID string,
) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM product_likes  WHERE user_id = $1 AND product_id = $2)`

	err := r.DB.QueryRow(c, query, userID, productID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// ToggleLike - switch like status
func (r *RepoProduct) ToggleLike(c context.Context, userID, productID string) (bool, error) {
	// Start transaction
	log.Println("[DEBUG USERID]", userID)
	log.Println("[DEBUG PRODUCTID]", productID)
	tx, err := r.DB.Begin(c)
	if err != nil {
		return false, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(c)
		}
	}()

	// Check current status
	var exists bool
	err = tx.QueryRow(c, `SELECT EXISTS(SELECT 1 FROM ratings WHERE user_id = $1::uuid AND product_id = $2::uuid)`, userID, productID).Scan(&exists)

	if err != nil {
		log.Println("[DEBUG LIKE1]", err)
		return false, err
	}

	// Toggle based on current status
	if exists {
		_, err = tx.Exec(c, `DELETE FROM ratings WHERE user_id = $1::uuid AND product_id = $2::uuid`, userID, productID)
		if err != nil {
			log.Println("[DEBUG LIKE2]", err)

			return false, err
		}
	} else {
		_, err = tx.Exec(c, `INSERT INTO ratings (user_id, product_id) VALUES ($1::uuid, $2::uuid)`, userID, productID)
		if err != nil {
			log.Println("[DEBUG LIKE3]", err)

			return false, err
		}
	}

	// Commit transaction
	if err = tx.Commit(c); err != nil {
		return false, err
	}

	return !exists, nil
}
