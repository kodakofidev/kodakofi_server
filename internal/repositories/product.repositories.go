package repositories

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kodakofidev/kodakofi_server/internal/models"
)

type ProductRepoInterface interface {
	GetAllProducts(c context.Context, params *models.ProductQueryParams) (*models.PaginatedResponse, error)
	GetDetailProduct()
}

type RepoProduct struct {
	DB *pgxpool.Pool
}

func NewProduct(db *pgxpool.Pool) *RepoProduct {
	return &RepoProduct{DB: db}
}

func (r *RepoProduct) GetAllProducts(c context.Context, params *models.ProductQueryParams) (*models.PaginatedResponse, error) {
	const pageSize = 6
	offset := (params.Page - 1) * pageSize
	query := ` select p.id, p.name, p.category_id,p.price, p.description,d.name AS discount_name,d.discount,SUM(po.qty) AS total_order,json_agg(pi.path) AS images,COUNT(r.*) AS total_ratings,c.name AS category_name FROM products p LEFT JOIN product_discounts pd ON pd.product_id = p.id LEFT JOIN discounts d ON d.id = pd.discount_id LEFT JOIN products_orders po ON po.product_id = p.id LEFT JOIN orders o ON o.id = po.order_id LEFT JOIN product_images pi ON pi.product_id = p.id LEFT JOIN ratings r ON r.product_id = p.id JOIN categories c ON c.id = p.category_id`
	var whereClauses []string
	args := []any{pageSize, offset, params.Min, params.Max}
	argIndex := 3
	if params.Min > 0 || params.Max > 0 {
		whereClauses = append(whereClauses, fmt.Sprintf("p.price BETWEEN $%d AND $%d", argIndex, argIndex+1))
		argIndex += 2
	}
	if params.Name != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("p.name ILIKE $%d", argIndex))
		args = append(args, "%"+params.Name+"%")
		argIndex++
	}
	if params.Category != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("c.name = $%d", argIndex))
		args = append(args, params.Category)
		argIndex++
	}
	if params.Discount != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("d.discount is not null"))
	}

	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}
	query += ` GROUP BY p.id, d.id, c.name ORDER BY  `

	if params.Options == "newest" {
		whereClauses = append(whereClauses, "p.created_at >= NOW()")
	}

	switch params.Options {
	case "oldest":
		query += "p.price DESC"
	case "asc":
		query += "total_order DESC"
	case "desc":
		query += "average_rating DESC"
	case "cheapest":
		query += "average_rating DESC"
	case "favorite":
		query += "average_rating DESC"
	default:
		query += "p.price ASC"
	}
	query += `LIMIT $1 OFFSET $2`
	rows, err := r.DB.Query(c, query, args...)
	if err != nil {
		return nil, err
	}
	var products models.Products
	for rows.Next() {
		var product models.Product
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.CategoryID,
			&product.Price,
			&product.Description,
			&product.DiscountName,
			&product.Discount,
			&product.TotalOrder,
			&product.Images,
			&product.TotalRatings,
			&product.CategoryName,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	countQuery := `SELECT COUNT(DISTINCT p.id)FROM products pLEFT JOIN product_discounts pd ON pd.product_id = p.idLEFT JOIN discounts d ON d.id = pd.discount_idJOIN categories c ON c.id = p.category_id`

	var totalItems int
	err = r.DB.QueryRow(c, countQuery, args[:len(args)-2]...).Scan(&totalItems)
	if err != nil {
		return nil, errors.New("failed to count products")
	}

	totalPages := totalItems / pageSize
	if totalItems%pageSize > 0 {
		totalPages++
	}
	response := &models.PaginatedResponse{
		Data: products,
		Pagination: models.Pagination{
			Page:       params.Page,
			PageSize:   pageSize,
			TotalItems: totalItems,
			TotalPages: totalPages,
		},
	}
	return response, nil
}

func (r *RepoProduct) GetDetailProduct() {

}
