package seed

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SeedPrductsOrders(ctx context.Context, db *pgxpool.Pool) error {
	query := `
		INSERT INTO public.products_orders (order_id,product_id,base_price,qty,added_price,sub_total)
		VALUES
			(1,'4e841656-596c-434d-b5bb-1f27f5d7418c'::uuid,25000,2,0,50000),
			(1,'4e841656-596c-434d-b5bb-1f27f5d7418c'::uuid,25000,1,6250,31250),
			(1,'9a7b950f-3664-4df3-9da5-249e91de2b31'::uuid,22000,1,11000,33000),
			(1,'8511cb06-1612-4b5a-9ef5-143abdc2077a'::uuid,3000,4,0,12000),
			(2,'4e841656-596c-434d-b5bb-1f27f5d7418c'::uuid,25000,1,6250,31250),
			(2,'9a7b950f-3664-4df3-9da5-249e91de2b31'::uuid,22000,2,0,44000),
			(2,'8511cb06-1612-4b5a-9ef5-143abdc2077a'::uuid,3000,3,0,9000),
			(3,'4e841656-596c-434d-b5bb-1f27f5d7418c'::uuid,25000,1,6250,31250),
			(3,'9a7b950f-3664-4df3-9da5-249e91de2b31'::uuid,22000,2,0,44000),
			(3,'5740f5fe-5178-4933-9e7b-63cab72fd79a'::uuid,30000,1,0,30000),
			(3,'8511cb06-1612-4b5a-9ef5-143abdc2077a'::uuid,3000,4,0,12000)
		ON CONFLICT (name) DO NOTHING;
	`

	_, err := db.Exec(ctx, query)
	if err != nil {
		log.Printf("Failed to seed products_orders: %v", err)
		return err
	}

	log.Println("Seeded products_orders successfully.")
	return nil
}
