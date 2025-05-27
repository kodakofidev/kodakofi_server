package seed

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SeedRatings(ctx context.Context, db *pgxpool.Pool) error {
	query := `
		INSERT INTO public.ratings (user_id, product_id, created_at, updated_at)
		VALUES
			('171b67f9-a3c2-4115-b75c-1a7816ca1cd2', '4e841656-596c-434d-b5bb-1f27f5d7418c', '2025-05-19 10:15:00+07', NULL),
			('171b67f9-a3c2-4115-b75c-1a7816ca1cd2', '5740f5fe-5178-4933-9e7b-63cab72fd79a', '2025-05-19 11:00:00+07', NULL),
			('0ee41132-5592-4626-8f87-1df70db54e3f', 'e6622f29-2a8b-4110-9a82-8616eed29570', '2025-05-19 12:30:00+07', NULL),
			('0ee41132-5592-4626-8f87-1df70db54e3f', '9a7b950f-3664-4df3-9da5-249e91de2b31', '2025-05-19 13:00:00+07', NULL),
			('8a60a062-c0a8-4034-b2f1-a832b9adb435', '8afc2e72-ef45-45b0-a936-62d34bd626bf', '2025-05-19 14:00:00+07', NULL),
			('8a60a062-c0a8-4034-b2f1-a832b9adb435', '7af264d9-fa31-45a4-8948-b2db4c267fd6', '2025-05-19 14:30:00+07', NULL),
			('171b67f9-a3c2-4115-b75c-1a7816ca1cd2', 'f866f4f6-f89c-4395-90ca-241dfb52951c', '2025-05-19 15:00:00+07', NULL),
			('0ee41132-5592-4626-8f87-1df70db54e3f', '95a70a1a-10c8-4a3f-80ad-6c430b74ef3e', '2025-05-19 15:30:00+07', NULL),
			('8a60a062-c0a8-4034-b2f1-a832b9adb435', '40425c10-f932-4b44-97a6-681b56a5ddfa', '2025-05-19 16:00:00+07', NULL),
			('171b67f9-a3c2-4115-b75c-1a7816ca1cd2', 'aead3bb8-eaa3-4408-a192-cc36b227f464', '2025-05-19 16:30:00+07', NULL),
			('0ee41132-5592-4626-8f87-1df70db54e3f', 'a2b74af4-06cc-4004-989e-6150af06926c', '2025-05-19 17:00:00+07', NULL),
			('8a60a062-c0a8-4034-b2f1-a832b9adb435', 'dfd726cf-8c5c-44ac-8d06-82378bb4c31c', '2025-05-19 17:30:00+07', NULL)
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
