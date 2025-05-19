package seed

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SeedRatings(ctx context.Context, db *pgxpool.Pool) error {
	query := `
		INSERT INTO public.ratings (user_id, product_id, rating, created_at)
		VALUES
			('8a60a062-c0a8-4034-b2f1-a832b9adb435', '4e841656-596c-434d-b5bb-1f27f5d7418c', 5, '2025-05-01 10:15:00+07'),
			('8a60a062-c0a8-4034-b2f1-a832b9adb435', '5740f5fe-5178-4933-9e7b-63cab72fd79a', 4, '2025-05-03 11:20:00+07'),
			('171b67f9-a3c2-4115-b75c-1a7816ca1cd2', 'e6622f29-2a8b-4110-9a82-8616eed29570', 4, '2025-04-29 09:40:00+07'),
			('171b67f9-a3c2-4115-b75c-1a7816ca1cd2', 'a2b74af4-06cc-4004-989e-6150af06926c', 5, '2025-05-02 12:00:00+07'),
			('171b67f9-a3c2-4115-b75c-1a7816ca1cd2', 'c504c3ea-af18-4bd8-a2e9-d7e773b1ea5d', 3, '2025-05-06 07:15:00+07')
		ON CONFLICT (user_id, product_id) DO NOTHING;
	`

	_, err := db.Exec(ctx, query)
	if err != nil {
		log.Printf("Failed to seed ratings: %v", err)
		return err
	}

	log.Println("Seeded products ratings successfully.")
	return nil
}
