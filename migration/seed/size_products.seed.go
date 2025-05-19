package seed

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SeedSizeProducts1(ctx context.Context, db *pgxpool.Pool) error {
	query := `
		WITH size_options AS (
			SELECT unnest(array[1, 2, 3]) AS size_id
		)
		INSERT INTO size_products (product_id, size_id, stock)
		SELECT 
			p.id,
			s.size_id,
			CASE s.size_id
				WHEN 1 THEN 100
				WHEN 2 THEN 60
				WHEN 3 THEN 30
			END AS stock
		FROM products p
		CROSS JOIN size_options s
		WHERE p.category_id IN (1, 2)
		ON CONFLICT (product_id, size_id) DO NOTHING;
	`

	_, err := db.Exec(ctx, query)
	if err != nil {
		log.Printf("Failed to seed size_products 1: %v", err)
		return err
	}

	log.Println("Seeded size_products 1 successfully.")
	return nil
}

func SeedSizeProducts2(ctx context.Context, db *pgxpool.Pool) error {
	query := `
		INSERT INTO size_products (product_id, size_id, stock)
		SELECT p.id, 4, 
			CASE 
				WHEN p.category_id IN (3,4,5) THEN 50
				ELSE 10  -- kategori 6 (topping dll) stok 10
			END as stock
		FROM products p
		WHERE p.category_id IN (3,4,5,6)
		ON CONFLICT (product_id, size_id) DO NOTHING;
	`

	_, err := db.Exec(ctx, query)
	if err != nil {
		log.Printf("Failed to seed size_products 2: %v", err)
		return err
	}

	log.Println("Seeded size_products 2 successfully.")
	return nil
}

func SeedSizeProducts(ctx context.Context, db *pgxpool.Pool) error {
	if err := SeedSizeProducts1(ctx, db); err != nil {
		return err
	}
	if err := SeedSizeProducts2(ctx, db); err != nil {
		return err
	}
	return nil
}
