package seed

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SeedPaymentMethods(ctx context.Context, db *pgxpool.Pool) error {
	query := `
		INSERT INTO payment_methods (name)
		VALUES 
			('BRI'),
			('DANA'),
			('BCA'),
			('GoPay'),
			('OVO'),
			('QRIS')
		ON CONFLICT (name) DO NOTHING;
	`

	_, err := db.Exec(ctx, query)
	if err != nil {
		log.Printf("Failed to seed payment_methods: %v", err)
		return err
	}

	log.Println("Seeded payment_methods successfully.")
	return nil
}
