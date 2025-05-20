package repositories

import (
	"context"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kodakofidev/kodakofi_server/internal/models"
)

type OrderRepoInterface interface {
	CreateOrder(ctx context.Context, data *models.CreateOrderRequest) (*models.CreateOrderResponse, error)
	GetHistoryOrders(ctx context.Context, offset int, status, userId string) ([]models.OrderHistory, error)
	GetTotalSales(ctx context.Context, startDate, endDate time.Time) ([]models.TotalSalesItemReponse, error)
	GetIncomeSales(ctx context.Context, startDate, endDate time.Time, page, limit int, baseURL string) (*models.MetaData, error)
}

type RepoOrder struct {
	DB *pgxpool.Pool
}

func NewOrder(db *pgxpool.Pool) *RepoOrder {
	return &RepoOrder{DB: db}
}

func isDrink(ctx context.Context, tx pgx.Tx, productID string) (bool, error) {
	var categoryID int
	err := tx.QueryRow(ctx, "SELECT category_id FROM products WHERE id = $1", productID).Scan(&categoryID)
	if err != nil {
		return false, err
	}
	return categoryID == 1 || categoryID == 2, nil
}

func (r *RepoOrder) CreateOrder(ctx context.Context, data *models.CreateOrderRequest) (*models.CreateOrderResponse, error) {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// 1. Ambil user_id dari email
	var userID string
	err = tx.QueryRow(ctx, "SELECT id FROM users WHERE email = $1", data.Email).Scan(&userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// 2. Ambil status_id "pending"
	var statusID int
	err = tx.QueryRow(ctx, "SELECT id FROM status WHERE status = 'Pending'").Scan(&statusID)
	if err != nil {
		return nil, fmt.Errorf("failed to get status_id: %w", err)
	}

	// 3. Insert ke orders
	var orderID int
	err = tx.QueryRow(ctx, `
		INSERT INTO orders (user_id, fullname, address, delivery_method_id, payment_method_id, status_id)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		userID, data.Fullname, data.Address, data.DeliveryMethodID, data.PaymentMethodID, statusID).Scan(&orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert order: %w", err)
	}

	total := 0
	var itemsResponse []models.OrderItemResponse

	// Preload Ice Cube data
	var iceCubeID string
	var iceCubePrice float64
	err = tx.QueryRow(ctx, `
		SELECT id, price FROM products
		WHERE category_id = 6 AND name = 'Ice Cube' LIMIT 1`).Scan(&iceCubeID, &iceCubePrice)
	if err != nil {
		return nil, fmt.Errorf("ice cube not found: %w", err)
	}

	iceCubeQtyTotal := 0

	// 4. Handle setiap item
	for _, item := range data.Items {
		isDrinkProduct, err := isDrink(ctx, tx, item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("failed to check product category: %w", err)
		}

		if item.SizeID == 0 {
			if isDrinkProduct {
				return nil, fmt.Errorf("size_id is required for drink products")
			} else {
				item.SizeID = 4
				item.IsIced = false
			}
		}

		// Ambil kategori produk
		// var categoryID int
		// err = tx.QueryRow(ctx, "SELECT category_id FROM products WHERE id = $1", item.ProductID).Scan(&categoryID)
		// if err != nil {
		// 	return nil, fmt.Errorf("failed to get product category: %w", err)
		// }

		// 4b. Ambil harga produk dan nama produk
		var basePrice float64
		var productName string

		err = tx.QueryRow(ctx, "SELECT name, price FROM products WHERE id = $1", item.ProductID).Scan(&productName, &basePrice)
		if err != nil {
			return nil, fmt.Errorf("product not found: %w", err)
		}

		// 4a. Ambil size_id dan added_price berdasarkan item.Size
		// var sizeID int = item.SizeID
		var sizeName string
		var addedPrice float64

		err = tx.QueryRow(ctx, "SELECT size, added_price FROM sizes WHERE id = $1", item.SizeID).
			Scan(&sizeName, &addedPrice)
		if err != nil {
			return nil, fmt.Errorf("invalid size_id: %w", err)
		}

		// 4c. Hitung harga total produk
		subTotal := int(math.Round((basePrice + (basePrice * addedPrice)) * float64(item.Qty)))
		total += subTotal

		// 4c. Update stok size_product
		cmdTag, err := tx.Exec(ctx, `
			UPDATE size_products SET stock = stock - $1
			WHERE product_id = $2 AND size_id = $3 AND stock >= $1`,
			item.Qty, item.ProductID, item.SizeID)
		if err != nil {
			return nil, fmt.Errorf("failed to update stock: %w", err)
		}
		if cmdTag.RowsAffected() == 0 {
			return nil, fmt.Errorf("insufficient stock for product %s with size %d", item.ProductID, item.SizeID)
		}

		// 4d. Insert ke products_orders
		_, err = tx.Exec(ctx, `
			INSERT INTO products_orders (order_id, product_id, base_price, size, is_iced, qty, added_price, sub_total)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			orderID, item.ProductID, int(basePrice), sizeName, item.IsIced, item.Qty, int(basePrice*addedPrice), subTotal)
		if err != nil {
			return nil, fmt.Errorf("failed to insert product_order: %w", err)
		}

		// 4f. Jika IsIced, tambahkan Ice Cube
		if item.IsIced {
			iceCubeQtyTotal += item.Qty
		}

		// Tambahkan ke response
		itemsResponse = append(itemsResponse, models.OrderItemResponse{
			ProductID:   item.ProductID,
			ProductName: productName,
			Qty:         item.Qty,
			Size:        sizeName,
			IsIced:      item.IsIced,
		})
	}

	// Tambahkan Ice Cube ke products_orders jika dibutuhkan
	if iceCubeQtyTotal > 0 {
		iceTotal := int(iceCubePrice) * iceCubeQtyTotal
		total += iceTotal

		// Cek apakah sudah ada Ice Cube di order
		var exists bool
		err = tx.QueryRow(ctx, `
				SELECT EXISTS (
					SELECT 1 FROM products_orders
					WHERE order_id = $1 AND product_id = $2
				)`, orderID, iceCubeID).Scan(&exists)
		if err != nil {
			return nil, fmt.Errorf("failed to check ice cube order: %w", err)
		}

		if exists {
			_, err = tx.Exec(ctx, `
				UPDATE products_orders
				SET qty = qty + $1, sub_total = sub_total + $2
				WHERE order_id = $3 AND product_id = $4`,
				iceCubeQtyTotal, iceTotal, orderID, iceCubeID)
			if err != nil {
				return nil, fmt.Errorf("failed to update ice cube order: %w", err)
			}
		} else {
			_, err = tx.Exec(ctx, `
					INSERT INTO products_orders (order_id, product_id, base_price, size, is_iced, qty, added_price, sub_total)
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
				orderID, iceCubeID, int(iceCubePrice), "Not Drink", false, iceCubeQtyTotal, 0, iceTotal)
			if err != nil {
				return nil, fmt.Errorf("failed to insert ice cube order: %w", err)
			}
		}
	}

	// 5. Hitung tax 12% dan total_amount
	var taxRate float64
	err = tx.QueryRow(ctx, `SELECT tax_value FROM tax ORDER BY created_at DESC LIMIT 1`).Scan(&taxRate)
	if err != nil {
		return nil, fmt.Errorf("failed to get tax: %w", err)
	}

	tax := int(float64(total) * taxRate)

	var deliveryFee int
	err = tx.QueryRow(ctx, `SELECT fee FROM delivery_methods WHERE id = $1`, data.DeliveryMethodID).Scan(&deliveryFee)
	if err != nil {
		return nil, fmt.Errorf("failed to get delivery fee: %w", err)
	}

	totalAmount := total + tax + deliveryFee

	// Generate Transactions Code
	now := time.Now()
	dateTimeStr := now.Format("20060102150405")
	uuidStr := strings.ToUpper(strings.ReplaceAll(uuid.New().String(), "-", "")[:8])
	transactionCode := fmt.Sprintf("TRX%s%s", dateTimeStr, uuidStr)

	_, err = tx.Exec(ctx, `
		INSERT INTO transactions (transaction_code, order_id, total, tax, delivery_fee, total_amount)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		transactionCode, orderID, total, tax, deliveryFee, totalAmount)
	if err != nil {
		return nil, fmt.Errorf("failed to insert transaction: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &models.CreateOrderResponse{
		Email:          data.Email,
		Fullname:       data.Fullname,
		Address:        data.Address,
		DeliveryMethod: data.DeliveryMethodID,
		PaymentMethod:  data.PaymentMethodID,
		Items:          itemsResponse,
		DeliveryFee:    deliveryFee,
		Total:          total,
		Tax:            tax,
		TotalAmount:    totalAmount,
	}, nil

}

// repo get history orders
func (r *RepoOrder) GetHistoryOrders(ctx context.Context, offset int, status, userId string) ([]models.OrderHistory, error) {

	// query := "select t.transaction_code, o.created_at, t.total_amount, o.id, s.status from orders o join transactions t on o.id = t.order_id join status s on s.id = o.status_id where o.user_id = $1 "
	query := `select DISTINCT ON (o.id) t.transaction_code, o.created_at, t.total_amount, o.id AS order_id, s.status, pi2."path" FROM orders o LEFT JOIN transactions t ON o.id = t.order_id left JOIN status s ON s.id = o.status_id LEFT JOIN products_orders po ON po.order_id = o.id left JOIN product_images pi2 ON pi2.product_id = po.product_id WHERE o.user_id = $1`

	value := []interface{}{userId}
	valueIndex := 2

	if status != "" {
		query += fmt.Sprintf(" and s.status = $%d", valueIndex)
		value = append(value, status)
		valueIndex++
	}

	if offset != -1 {
		query += fmt.Sprintf(" limit 4 offset $%d", valueIndex)
		value = append(value, offset)
	}

	rows, err := r.DB.Query(ctx, query, value...)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	defer rows.Close()
	var result []models.OrderHistory

	for rows.Next() {
		var history models.OrderHistory
		if err := rows.Scan(&history.TransactionCode, &history.Date, &history.GrandTotal, &history.OrderId, &history.Status, &history.Path); err != nil {
			log.Println(err.Error())
			return nil, err
		}
		result = append(result, history)
	}
	return result, nil
}

func (r *RepoOrder) GetTotalSales(ctx context.Context, startDate, endDate time.Time) ([]models.TotalSalesItemReponse, error) {
	query := `
		SELECT
			DATE(o.created_at) as date,
			SUM(po.qty) as total_items
		FROM orders o
		JOIN products_orders po ON o.id = po.order_id
		WHERE o.created_at BETWEEN $1 AND $2
		GROUP BY date
		ORDER BY date ASC
	`

	rows, err := r.DB.Query(ctx, query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.TotalSalesItemReponse
	var totalItems int

	for rows.Next() {
		var item models.TotalSalesItemReponse
		err := rows.Scan(&item.Date, &item.Item)
		if err != nil {
			return nil, err
		}
		totalItems += item.Item
		result = append(result, item)
	}

	// Tambahkan total ke semua elemen agar dapat digunakan untuk summary di FE
	for i := range result {
		result[i].TotalItemOrder = totalItems
	}

	return result, nil
}

func (r *RepoOrder) GetIncomeSales(ctx context.Context, startDate, endDate time.Time, page, limit int, baseURL string) (*models.MetaData, error) {
	offset := (page - 1) * limit

	query := `
		SELECT
			DATE(o.created_at) as order_date,
			p.name as product_name,
			SUM(po.qty) as total_item_order,
			SUM(po.sub_total) as income
		FROM products_orders po
		JOIN orders o ON o.id = po.order_id
		JOIN products p ON p.id = po.product_id
		JOIN transactions t ON o.id = t.order_id
		WHERE o.created_at BETWEEN $1 AND $2
		GROUP BY order_date, p.name
		ORDER BY income DESC
		OFFSET $3 LIMIT $4
	`

	rows, err := r.DB.Query(ctx, query, startDate, endDate, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.TotalIncomeItemReponse
	for rows.Next() {
		var item models.TotalIncomeItemReponse
		if err := rows.Scan(&item.Date, &item.ProductName, &item.TotalItemOrder, &item.Income); err != nil {
			return nil, err
		}
		results = append(results, item)
	}

	// Hitung total data
	var totalData int
	countQuery := `
		SELECT COUNT(*) FROM (
			SELECT DISTINCT DATE(o.created_at), p.name
			FROM products_orders po
			JOIN orders o ON o.id = po.order_id
			JOIN products p ON p.id = po.product_id
			WHERE o.created_at BETWEEN $1 AND $2
		) AS subquery
	`
	err = r.DB.QueryRow(ctx, countQuery, startDate, endDate).Scan(&totalData)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalData) / float64(limit)))

	meta := models.MetaData{
		Data:       results,
		TotalData:  totalData,
		Page:       page,
		TotalPages: totalPages,
	}

	if page < totalPages {
		meta.NextLink = fmt.Sprintf("%s?page=%d&limit=%d&start_date=%s&end_date=%s", baseURL, page+1, limit, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	}
	if page > 1 {
		meta.PrevLink = fmt.Sprintf("%s?page=%d&limit=%d&start_date=%s&end_date=%s", baseURL, page-1, limit, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	}

	return &meta, nil
}
