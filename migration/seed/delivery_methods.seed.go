package seed

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SeedDeliveryMethods(ctx context.Context, db *pgxpool.Pool) error {
	query := `
		INSERT INTO delivery_methods (name, fee)
		VALUES 
			('Dine In', 0),
			('Door Delivery', 15000),
			('Pick Up', 0)
		ON CONFLICT (name) DO NOTHING;
	`

	_, err := db.Exec(ctx, query)
	if err != nil {
		log.Printf("Failed to seed delivery_methods: %v", err)
		return err
	}

	log.Println("Seeded delivery_methods successfully.")
	return nil
}
