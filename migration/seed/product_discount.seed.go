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
	('ec4e8c26-b66f-4964-ab5a-4dd8f7ec0e05', 1),
	('9a283cef-5a9f-437f-ac24-cf489094c7aa', 1),
	('839f121f-8351-4e2c-8e10-d8ac81cc3e45', 1),
	('f30b51f9-9877-4e7f-b81a-3406c20b9629', 2),
	('ca0905eb-2088-4a11-af8f-7b76601cc206', 2),
	('20b2c260-292f-4508-bff4-d377acc2fb8c', 2),
	('c4e99f87-efd0-49a3-bd51-e1ba499f5b0b', 3),
	('15513685-a713-4830-a0bc-c1a99902ce35', 3),
	('3252ddbd-1c69-49b1-86fc-f52312ac8c8c', 3)
	`

	_, err := db.Exec(ctx, query)
	if err != nil {
		log.Printf("Failed to seed product discounts: %v", err)
		return err
	}

	log.Println("Seeded product discount successfully.")
	return nil
}
