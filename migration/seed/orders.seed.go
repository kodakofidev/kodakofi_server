package seed

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SeedOrders(ctx context.Context, db *pgxpool.Pool) error {
	query := `
		INSERT INTO public.orders (user_id,fullname,address,delivery_method_id,payment_method_id,status_id,created_at,updated_at) VALUES
			('171b67f9-a3c2-4115-b75c-1a7816ca1cd2'::uuid,'sdgi6','Jl. Merdeka No. 123',1,1,4,'2025-05-18 20:57:22.586127+07',NULL),
			('171b67f9-a3c2-4115-b75c-1a7816ca1cd2'::uuid,'sdgi6','Jl. Merdeka No. 123',1,1,4,'2025-05-18 20:57:53.860762+07',NULL),
			('171b67f9-a3c2-4115-b75c-1a7816ca1cd2'::uuid,'sdgi6','Jl. Merdeka No. 123',1,1,4,'2025-05-18 21:03:24.705226+07',NULL)
		ON CONFLICT (name) DO NOTHING;
	`

	_, err := db.Exec(ctx, query)
	if err != nil {
		log.Printf("Failed to seed orders: %v", err)
		return err
	}

	log.Println("Seeded orders successfully.")
	return nil
}
