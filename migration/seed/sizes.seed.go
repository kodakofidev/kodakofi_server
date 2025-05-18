package seed

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SeedSizes(ctx context.Context, db *pgxpool.Pool) error {
	query := `
		INSERT INTO sizes (size, added_price)
		VALUES 
			('Reguler', 0),
			('Medium', 0.25),
			('Large', 0.5),
			('Not Drink', 0)
		ON CONFLICT (size) DO NOTHING;
	`

	_, err := db.Exec(ctx, query)
	if err != nil {
		log.Printf("Failed to seed sizes: %v", err)
		return err
	}

	log.Println("Seeded sizes successfully.")
	return nil
}
