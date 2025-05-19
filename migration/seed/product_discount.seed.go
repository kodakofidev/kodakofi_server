package seed

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SeedProductDiscount(ctx context.Context, db *pgxpool.Pool) error {
	query := `
	INSERT INTO product_discounts (product_id, discount_id)
	VALUES
	('287ec09d-928c-4562-9f29-86ad95dce6f6'::uuid, 1),
	('aead3bb8-eaa3-4408-a192-cc36b227f464'::uuid, 1),
	('e6622f29-2a8b-4110-9a82-8616eed29570'::uuid, 1),
	('40425c10-f932-4b44-97a6-681b56a5ddfa'::uuid, 2),
	('5740f5fe-5178-4933-9e7b-63cab72fd79a'::uuid, 2),
	('31b9935a-7bd3-4ae0-898d-386e4cffb82e'::uuid, 2),
	('0076aee4-9db2-4941-a69d-e07ff562dc3b'::uuid, 3),
	('4e841656-596c-434d-b5bb-1f27f5d7418c'::uuid, 3),
	('9c4817f0-455e-415f-9c79-8a7e6f4fc1ab'::uuid, 3)
	`

	_, err := db.Exec(ctx, query)
	if err != nil {
		log.Printf("Failed to seed product discounts: %v", err)
		return err
	}

	log.Println("Seeded product discount successfully.")
	return nil
}
