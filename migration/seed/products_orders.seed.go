package seed

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SeedPrductsOrders(ctx context.Context, db *pgxpool.Pool) error {
	query := `
		INSERT INTO public.products_orders (order_id,product_id,base_price,size,is_iced,qty,added_price,sub_total)
		VALUES
			(6,'4e841656-596c-434d-b5bb-1f27f5d7418c'::uuid,25000,'Medium',true,1,6250,31250),
			(6,'4e841656-596c-434d-b5bb-1f27f5d7418c'::uuid,25000,'Reguler',true,1,0,25000),
			(6,'4e841656-596c-434d-b5bb-1f27f5d7418c'::uuid,25000,'Large',true,1,12500,37500),
			(6,'4e841656-596c-434d-b5bb-1f27f5d7418c'::uuid,25000,'Reguler',false,1,0,25000),
			(6,'9a7b950f-3664-4df3-9da5-249e91de2b31'::uuid,22000,'Reguler',true,1,0,22000),
			(6,'978feadd-1c68-479f-a99a-b831b732464b'::uuid,45000,'Not Drink',false,1,0,45000),
			(6,'8511cb06-1612-4b5a-9ef5-143abdc2077a'::uuid,3000,'Not Drink',false,4,0,12000),
			(9,'4e841656-596c-434d-b5bb-1f27f5d7418c'::uuid,25000,'Medium',true,1,6250,31250),
			(9,'4e841656-596c-434d-b5bb-1f27f5d7418c'::uuid,25000,'Large',true,2,12500,75000),
			(9,'4e841656-596c-434d-b5bb-1f27f5d7418c'::uuid,25000,'Reguler',true,1,0,25000),
			(9,'4e841656-596c-434d-b5bb-1f27f5d7418c'::uuid,25000,'Reguler',false,1,0,25000),
			(9,'9a7b950f-3664-4df3-9da5-249e91de2b31'::uuid,22000,'Medium',false,1,5500,27500),
			(9,'978feadd-1c68-479f-a99a-b831b732464b'::uuid,45000,'Not Drink',false,2,0,90000),
			(9,'8511cb06-1612-4b5a-9ef5-143abdc2077a'::uuid,3000,'Not Drink',false,4,0,12000),
			(10,'4e841656-596c-434d-b5bb-1f27f5d7418c'::uuid,25000,'Reguler',true,5,0,125000),
			(10,'978feadd-1c68-479f-a99a-b831b732464b'::uuid,45000,'Not Drink',false,2,0,90000),
			(10,'8511cb06-1612-4b5a-9ef5-143abdc2077a'::uuid,3000,'Not Drink',false,5,0,15000)
		ON CONFLICT (order_id, product_id, size, is_iced) DO NOTHING;
	`

	_, err := db.Exec(ctx, query)
	if err != nil {
		log.Printf("Failed to seed products_orders: %v", err)
		return err
	}

	log.Println("Seeded products_orders successfully.")
	return nil
}
