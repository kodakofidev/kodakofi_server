package seed

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SeedCategories(ctx context.Context, db *pgxpool.Pool) error {
	query := `
		INSERT INTO categories (name)
		VALUES 
			('Coffee'),
			('Non-Coffee'),
			('Food'),
			('Dessert'),
			('Snack'),
			('Topping')
		ON CONFLICT (name) DO NOTHING;
	`

	_, err := db.Exec(ctx, query)
	if err != nil {
		log.Printf("Failed to seed categories: %v", err)
		return err
	}

	log.Println("Seeded categories successfully.")
	return nil
}
