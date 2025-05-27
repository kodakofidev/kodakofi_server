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
	GetDataSales(ctx context.Context, startDate, endDate time.Time) (*models.ProductSalesDataRes, error)
	UpdateStatusOrder(ctx context.Context, orderID, statusID int) (*models.UpdateOrderStatusRes, error)
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

func (r *RepoOrder) GetDataSales(ctx context.Context, startDate, endDate time.Time) (*models.ProductSalesDataRes, error) {
	var res models.ProductSalesDataRes

	// 1. Ambil status order
	statusQuery := `
		SELECT s.status, COUNT(o.id)
		FROM orders o
		JOIN status s ON o.status_id = s.id
		WHERE o.created_at BETWEEN $1 AND $2
		GROUP BY s.status
	`
	rows, err := r.DB.Query(ctx, statusQuery, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	statusData := models.SalesDataStatus{}
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		switch status {
		case "Pending":
			statusData.Pending = count
		case "Processing":
			statusData.Processing = count
		case "Completed":
			statusData.Completed = count
		}
	}
	res.Status = statusData

	// 2. Ambil total item terjual per hari (FIXED)
	dailySalesQuery := `
		SELECT DATE(o.created_at) AS order_date, SUM(po.qty) AS total_qty
		FROM orders o
		JOIN products_orders po ON o.id = po.order_id
		WHERE o.created_at BETWEEN $1 AND $2
		GROUP BY order_date
		ORDER BY order_date
	`
	dailyRows, err := r.DB.Query(ctx, dailySalesQuery, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer dailyRows.Close()

	dailyItems := []models.DailySoldItems{}
	for dailyRows.Next() {
		var date time.Time
		var productsSold int
		if err := dailyRows.Scan(&date, &productsSold); err != nil {
			return nil, err
		}
		dailyItems = append(dailyItems, models.DailySoldItems{
			Date:         date.Format("2006-01-02"), // konversi ke string
			ProductsSold: productsSold,
		})
	}
	res.DailySoldItems = dailyItems

	// 3. Total seluruh item terjual
	err = r.DB.QueryRow(ctx, `
		SELECT COALESCE(SUM(po.qty), 0)
		FROM orders o
		JOIN products_orders po ON o.id = po.order_id
		WHERE o.created_at BETWEEN $1 AND $2
	`, startDate, endDate).Scan(&res.TotalSoldItems)
	if err != nil {
		return nil, err
	}

	// 4. Data income per produk
	incomeQuery := `
		SELECT p.name, SUM(po.qty) AS total_qty, SUM(po.sub_total) AS income
		FROM products_orders po
		JOIN products p ON po.product_id = p.id
		JOIN orders o ON o.id = po.order_id
		WHERE o.created_at BETWEEN $1 AND $2
		GROUP BY p.name
		ORDER BY total_qty DESC
	`
	incomeRows, err := r.DB.Query(ctx, incomeQuery, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer incomeRows.Close()

	incomeData := []models.TotalIncomeItemResponse{}
	for incomeRows.Next() {
		var item models.TotalIncomeItemResponse
		var incomeInt int
		if err := incomeRows.Scan(&item.ProductName, &item.TotalItemSold, &incomeInt); err != nil {
			return nil, err
		}
		item.Income = fmt.Sprintf("%d", incomeInt)
		incomeData = append(incomeData, item)
	}
	res.IncomeDataPerItem = incomeData
	res.TotalData = len(incomeData)

	return &res, nil
}

func (r *RepoOrder) UpdateStatusOrder(ctx context.Context, orderID, statusID int) (*models.UpdateOrderStatusRes, error) {
	query := `
		UPDATE orders
		SET
			status_id = $1,
			updated_at = now()
		WHERE id = $2
		RETURNING id,
		(SELECT status FROM status
			WHERE id = $1),
		(SELECT transaction_code FROM transactions
			WHERE order_id = $2),
		updated_at
	`

	var response models.UpdateOrderStatusRes
	err := r.DB.QueryRow(ctx, query, statusID, orderID).Scan(&response.OrderID, &response.Status, &response.TransactionCode, &response.UpdateAt)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
