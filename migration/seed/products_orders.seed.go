package seed

// import (
// 	"context"
// 	"log"

// 	"github.com/jackc/pgx/v5/pgxpool"
// )

// func SeedPrductsOrders(ctx context.Context, db *pgxpool.Pool) error {
// 	query := `
// 		INSERT INTO public.products_orders (order_id,product_id,base_price,size,qty,added_price,sub_total)
// 		VALUES
// 			(1,'4e841656-596c-434d-b5bb-1f27f5d7418c'::uuid,25000,Reguler,2,0,50000),
// 			(1,'8511cb06-1612-4b5a-9ef5-143abdc2077a'::uuid,3000,Reguler,4,0,12000),
// 			(2,'9a7b950f-3664-4df3-9da5-249e91de2b31'::uuid,22000,Reguler,2,0,44000),
// 			(2,'8511cb06-1612-4b5a-9ef5-143abdc2077a'::uuid,3000,Reguler,3,0,9000),
// 			(3,'9a7b950f-3664-4df3-9da5-249e91de2b31'::uuid,22000,Reguler,2,0,44000),
// 			(3,'8511cb06-1612-4b5a-9ef5-143abdc2077a'::uuid,3000,Reguler,4,0,12000)
// 		ON CONFLICT (order_id, product_id, size) DO NOTHING;
// 	`

// 	_, err := db.Exec(ctx, query)
// 	if err != nil {
// 		log.Printf("Failed to seed products_orders: %v", err)
// 		return err
// 	}

// 	log.Println("Seeded products_orders successfully.")
// 	return nil
// }
