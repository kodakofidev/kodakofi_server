package seed

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SeedTax(ctx context.Context, db *pgxpool.Pool) error {
	query := `
		INSERT INTO tax (tax_value)
		VALUES
			(0.12)
		ON CONFLICT (id) DO NOTHING;
	`

	_, err := db.Exec(ctx, query)
	if err != nil {
		log.Printf("Failed to seed tax: %v", err)
		return err
	}

	log.Println("Seeded tax successfully.")
	return nil
}
