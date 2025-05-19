package seed

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SeedProfiles(ctx context.Context, db *pgxpool.Pool) error {
	query := `
		INSERT INTO public.profiles (user_id,fullname,phone,address,image,created_at,updated_at) VALUES
			('c59b7e07-cbb6-4647-858a-406c3c4b4bd2'::uuid,'','','','','2025-05-18 18:28:38.854598+07',NULL),
			('171b67f9-a3c2-4115-b75c-1a7816ca1cd2'::uuid,'','','','','2025-05-18 20:47:26.90677+07',NULL),
			('0ee41132-5592-4626-8f87-1df70db54e3f'::uuid,'','','','','2025-05-19 03:49:14.898142+07',NULL),	 
			('8a60a062-c0a8-4034-b2f1-a832b9adb435'::uuid,'Imanul Amin Test Edit','081211112222','JL. Kutilang No. 123','/img/1747615441996359000_8a60a062-c0a8-4034-b2f1-a832b9adb435_profile.jpg','2025-05-18 20:42:23.05615+07',NULL)
		ON CONFLICT (user_id) DO NOTHING;
	`

	_, err := db.Exec(ctx, query)
	if err != nil {
		log.Printf("Failed to seed profiles: %v", err)
		return err
	}

	log.Println("Seeded profiles successfully.")
	return nil
}
