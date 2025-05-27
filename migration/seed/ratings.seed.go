package seed

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SeedRatings(ctx context.Context, db *pgxpool.Pool) error {
	query := `
		INSERT INTO ratings (user_id, product_id, created_at)
		VALUES 
			('171b67f9-a3c2-4115-b75c-1a7816ca1cd2', '4e841656-596c-434d-b5bb-1f27f5d7418c', '2025-05-01 10:00:00+07'),
			('171b67f9-a3c2-4115-b75c-1a7816ca1cd2', '5740f5fe-5178-4933-9e7b-63cab72fd79a', '2025-05-02 10:15:00+07'),
			('171b67f9-a3c2-4115-b75c-1a7816ca1cd2', 'e6622f29-2a8b-4110-9a82-8616eed29570', '2025-05-03 09:30:00+07'),
			('0ee41132-5592-4626-8f87-1df7011fd3f5', '9a7b950f-3664-4df3-9da5-249e91de2b31', '2025-05-01 08:20:00+07'),
			('0ee41132-5592-4626-8f87-1df7011fd3f5', '8afc2e72-ef45-45b0-a936-62d34bd626bf', '2025-05-03 11:00:00+07'),
			('0ee41132-5592-4626-8f87-1df7011fd3f5', '7af264d9-fa31-45a4-8948-b2db4c267fd6', '2025-05-04 13:00:00+07'),
			('c59b7e07-cbb6-4647-858a-406c3c4b4bd2', '287ec09d-928c-4562-9f29-86ad95dce6f6', '2025-05-02 14:10:00+07'),
			('c59b7e07-cbb6-4647-858a-406c3c4b4bd2', '95a70a1a-10c8-4a3f-80ad-6c430b74ef3e', '2025-05-04 16:30:00+07'),
			('c59b7e07-cbb6-4647-858a-406c3c4b4bd2', '40425c10-f932-4b44-97a6-681b56a5ddfa', '2025-05-05 10:40:00+07'),
			('171b67f9-a3c2-4115-b75c-1a7816ca1cd2', '0076aee4-9db2-4941-a69d-e07ff562dc3b', '2025-05-05 12:00:00+07'),
			('0ee41132-5592-4626-8f87-1df7011fd3f5', 'aead3bb8-eaa3-4408-a192-cc36b227f464', '2025-05-06 09:45:00+07'),
			('c59b7e07-cbb6-4647-858a-406c3c4b4bd2', 'a2b74af4-06cc-4004-989e-6150af06926c', '2025-05-07 11:20:00+07'),
			('171b67f9-a3c2-4115-b75c-1a7816ca1cd2', 'dfd726cf-8c5c-44ac-8d06-82378bb4c31c', '2025-05-08 13:30:00+07'),
			('0ee41132-5592-4626-8f87-1df7011fd3f5', '31b9935a-7bd3-4ae0-898d-386e4cffb82e', '2025-05-09 15:00:00+07'),
			('c59b7e07-cbb6-4647-858a-406c3c4b4bd2', '06bd407e-5fab-4590-9e3f-7dc442af3b42', '2025-05-10 16:45:00+07')
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
