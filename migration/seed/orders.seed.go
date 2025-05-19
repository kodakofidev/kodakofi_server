package seed

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SeedOrders(ctx context.Context, db *pgxpool.Pool) error {
	query := `
		INSERT INTO public.orders (id, user_id,fullname,address,delivery_method_id,payment_method_id,status_id,created_at,updated_at) VALUES
			(6, '171b67f9-a3c2-4115-b75c-1a7816ca1cd2'::uuid,'Redha Pradana','Jl. Tupai Rumput No. 1',2,3,1,'2025-05-19 15:15:15.906678+07',NULL),
			(9, '171b67f9-a3c2-4115-b75c-1a7816ca1cd2'::uuid,'Redha Pradana','Jl. Tupai Rumput No. 1',1,1,1,'2025-05-19 15:45:59.796077+07',NULL),
			(10, '8a60a062-c0a8-4034-b2f1-a832b9adb435'::uuid,'Imanul Redha','Jl. Tupai Rumput No. 1',2,5,1,'2025-05-19 15:55:19.922704+07',NULL)
		ON CONFLICT (id) DO NOTHING;
	`

	_, err := db.Exec(ctx, query)
	if err != nil {
		log.Printf("Failed to seed orders: %v", err)
		return err
	}

	log.Println("Seeded orders successfully.")
	return nil
}
