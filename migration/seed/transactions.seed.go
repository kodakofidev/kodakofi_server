package seed

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SeedTrasnsactions(ctx context.Context, db *pgxpool.Pool) error {
	query := `
		INSERT INTO public.transactions (transaction_code,order_id,total,delivery_fee,tax,total_amount,created_at,updated_at)
		VALUES
			('TRX20250518205722E2C80723',1,126250,0,15150,141400,'2025-05-18 20:57:22.586127+07',NULL),
			('TRX2025051820575353C0C8B7',2,84250,0,10110,94360,'2025-05-18 20:57:53.860762+07',NULL),
			('TRX20250518210324B819A924',3,117250,0,14070,131320,'2025-05-18 21:03:24.705226+07',NULL)
		ON CONFLICT (id) DO NOTHING;
	`

	_, err := db.Exec(ctx, query)
	if err != nil {
		log.Printf("Failed to seed transactions: %v", err)
		return err
	}

	log.Println("Seeded transactions successfully.")
	return nil
}
