package seed

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SeedDiscounts(ctx context.Context, db *pgxpool.Pool) error {
	query := `
		INSERT INTO discounts (id, name, discount, expired)
		VALUES
			(1, 'Weekend Treats', 0.15, '2025-05-25 16:59:00+07'),
			(2, 'May Day Sale', 0.2, '2025-05-31 16:59:00+07'),
			(3, 'Today Sale', 0.2, '2025-05-20 16:59:00+07'),
			(4, 'Midweek Madness', 0.1, '2025-05-22 23:59:00+07'),
			(5, 'End of Month Bonanza', 0.25, '2025-05-31 23:59:00+07'),
			(6, 'Flash Sale', 0.3, '2025-05-19 23:00:00+07'),
			(7, 'Coffee Lovers Week', 0.12, '2025-05-24 22:00:00+07'),
			(8, 'Buy 2 Get 1 Week', 0.1667, '2025-05-26 23:59:00+07'),
			(9, 'Payday Promo', 0.18, '2025-05-27 23:59:00+07'),
			(10, 'Customer Appreciation', 0.2, '2025-06-01 23:59:00+07')
		ON CONFLICT (name) DO NOTHING;
	`

	_, err := db.Exec(ctx, query)
	if err != nil {
		log.Printf("Failed to seed discounts: %v", err)
		return err
	}

	log.Println("Seeded discounts successfully.")
	return nil
}
