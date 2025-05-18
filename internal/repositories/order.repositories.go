package repositories

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kodakofidev/kodakofi_server/internal/models"
)

type OrderRepoInterface interface {
	CreateOrder(ctx context.Context, data *models.CreateOrderRequest) (*models.CreateOrderResponse, error)
}

type RepoOrder struct {
	DB *pgxpool.Pool
}

func NewOrder(db *pgxpool.Pool) *RepoOrder {
	return &RepoOrder{DB: db}
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

	// 4. Handle setiap item
	for _, item := range data.Items {
		// 4b. Ambil harga produk dan nama produk
		var basePrice float64
		var productName string

		err = tx.QueryRow(ctx, "SELECT name, price FROM products WHERE id = $1", item.ProductID).Scan(&productName, &basePrice)
		if err != nil {
			return nil, fmt.Errorf("product not found: %w", err)
		}

		// 4a. Ambil size_id dan added_price berdasarkan item.Size
		var addedPrice float64

		err = tx.QueryRow(ctx, "SELECT added_price FROM sizes WHERE id = $1", item.SizeID).Scan(&addedPrice)
		if err != nil {
			return nil, fmt.Errorf("size_id not found: %w", err)
		}

		// 4c. Hitung harga total produk
		subTotal := int((basePrice + (basePrice * addedPrice)) * float64(item.Qty))
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
			INSERT INTO products_orders (order_id, product_id, base_price, qty, added_price, sub_total)
			VALUES ($1, $2, $3, $4, $5, $6)`,
			orderID, item.ProductID, int(basePrice), item.Qty, int(basePrice*addedPrice), subTotal)
		if err != nil {
			return nil, fmt.Errorf("failed to insert product_order: %w", err)
		}

		// 4f. Jika iced, tambahkan Ice Cube
		if item.IsIced {
			iceCubeQty := 0
			var iceCubeID string
			var iceCubePrice float64
			if iceCubeID == "" {
				err = tx.QueryRow(ctx, `
				SELECT id, price FROM products
				WHERE category_id = 6 AND name = 'Ice Cube' LIMIT 1`).Scan(&iceCubeID, &iceCubePrice)
				if err != nil {
					return nil, fmt.Errorf("ice cube not found: %w", err)
				}
			}
			iceCubeQty += item.Qty

			if iceCubeQty > 0 {
				// Cek apakah sudah ada Ice Cube di products_orders
				var exists bool
				err = tx.QueryRow(ctx, `
					SELECT EXISTS (
						SELECT 1 FROM products_orders
						WHERE order_id = $1 AND product_id = $2
					)`, orderID, iceCubeID).Scan(&exists)
				if err != nil {
					return nil, fmt.Errorf("failed to check existing ice cube order: %w", err)
				}

				iceTotal := int(iceCubePrice) * item.Qty
				total += iceTotal

				// Insert Ice Cube ke products_orders
				if exists {
					// Update qty dan sub_total
					_, err = tx.Exec(ctx, `
					UPDATE products_orders
					SET qty = qty + $1, sub_total = sub_total + $2
					WHERE order_id = $3 AND product_id = $4`,
						iceCubeQty, iceTotal, orderID, iceCubeID)
					if err != nil {
						return nil, fmt.Errorf("failed to update ice cube order: %w", err)
					}
				} else {
					// Insert baru
					_, err = tx.Exec(ctx, `
					INSERT INTO products_orders (order_id, product_id, base_price, qty, added_price, sub_total)
					VALUES ($1, $2, $3, $4, $5, $6)`,
						orderID, iceCubeID, int(iceCubePrice), iceCubeQty, 0, iceTotal)
					if err != nil {
						return nil, fmt.Errorf("failed to insert ice cube: %w", err)
					}
				}
			}
		}

		itemsResponse = append(itemsResponse, models.OrderItemResponse{
			ProductID:   item.ProductID,
			ProductName: productName,
			Qty:         item.Qty,
			Size:        item.SizeID,
			IsIced:      item.IsIced,
		})
	}

	// 5. Hitung tax 12% dan total_amount
	var taxRate float64
	err = tx.QueryRow(ctx, `SELECT tax_value FROM tax ORDER BY created_at DESC LIMIT 1`).Scan(&taxRate)
	if err != nil {
		return nil, fmt.Errorf("failed to get tax: %w", err)
	}

	tax := int(float64(total) * taxRate)
	totalAmount := total + tax

	var deliveryFee int
	err = tx.QueryRow(ctx, `SELECT fee FROM delivery_methods WHERE id = $1`, data.DeliveryMethodID).Scan(&deliveryFee)
	if err != nil {
		return nil, fmt.Errorf("failed to get delivery fee: %w", err)
	}
	totalAmount += deliveryFee

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
		Tax:            tax,
		Total:          total,
		TotalAmount:    totalAmount,
	}, nil
}
