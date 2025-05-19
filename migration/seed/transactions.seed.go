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
			('TRX202505191515154DA093CB',6,197750,15000,23730,236480,'2025-05-19 15:15:15.906678+07',NULL),
			('TRX202505191545591B724006',9,285750,0,34290,320040,'2025-05-19 15:45:59.796077+07',NULL),
			('TRX202505191555201FBAD497',10,230000,15000,27600,272600,'2025-05-19 15:55:19.922704+07',NULL)
		ON CONFLICT (transaction_code) DO NOTHING;
	`

	_, err := db.Exec(ctx, query)
	if err != nil {
		log.Printf("Failed to seed transactions: %v", err)
		return err
	}

	log.Println("Seeded transactions successfully.")
	return nil
}
