package seed

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SeedDiscount(ctx context.Context, db *pgxpool.Pool) error {
	query := `
	INSERT INTO discounts (name, discount, expired)
	VALUES
	('Weekend Treats', 0.15, '2025-05-25 16:59:00+07'),
	('May Day Sale', 0.2, '2025-05-31 16:59:00+07'),
	('Today Sale', 0.2, '2025-05-20 16:59:00+07')
	`

	_, err := db.Exec(ctx, query)
	if err != nil {
		log.Printf("Failed to seed discounts: %v", err)
		return err
	}

	log.Println("Seeded discount successfully.")
	return nil
}
