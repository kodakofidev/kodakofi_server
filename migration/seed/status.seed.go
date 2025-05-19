package seed

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SeedStatus(ctx context.Context, db *pgxpool.Pool) error {
	query := `
		INSERT INTO status (status)
		VALUES 
			('Pending'),
			('Processing'),
			('Completed'),
			('Cancelled')
		ON CONFLICT (status) DO NOTHING;
	`

	_, err := db.Exec(ctx, query)
	if err != nil {
		log.Printf("Failed to seed status: %v", err)
		return err
	}

	log.Println("Seeded status successfully.")
	return nil
}
