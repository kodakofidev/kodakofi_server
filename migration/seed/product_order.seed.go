package seed

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SeedProductOrder(ctx context.Context, db *pgxpool.Pool) error {
	query := `
	INSERT INTO products_orders (product_id, order_id, qty)
	VALUES
	('c4e99f87-efd0-49a3-bd51-e1ba499f5b0b', 1, 1),
	('47b80d70-0512-4c11-9e3f-364563b40e4a', 1, 1)
	`

	_, err := db.Exec(ctx, query)
	if err != nil {
		log.Printf("Failed to seed product order: %v", err)
		return err
	}

	log.Println("Seeded product order successfully.")
	return nil
}
