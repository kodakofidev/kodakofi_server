package repositories

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kodakofidev/kodakofi_server/internal/models"
)

type OrderRepoInterface interface {
	CreateOrder()
	GetHistoryOrders(ctx context.Context, offset int, status, userId string) (models.OrderHistories, error)
}

type RepoOrder struct {
	DB *pgxpool.Pool
}

func NewOrder(db *pgxpool.Pool) *RepoOrder {
	return &RepoOrder{DB: db}
}

func (r *RepoOrder) CreateOrder() {

}

// repo get history orders
func (r *RepoOrder) GetHistoryOrders(ctx context.Context, offset int, status, userId string) (models.OrderHistories, error) {

	query := "select t.transaction_code, o.created_at, t.grand_total, o.id, s.status from orders o join transactions t on o.id = t.order_id join status s on s.id = o.status_id where o.user_id = $1 "


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
	var result models.OrderHistories

	for rows.Next() {
		var history models.OrderHistory
		if err := rows.Scan(&history.TransactionCode, &history.Date, &history.GrandTotal, &history.OrderId, &history.Status); err != nil {
			log.Println(err.Error())
			return nil, err
		}
		result = append(result, history)
	}
	return result, nil
}