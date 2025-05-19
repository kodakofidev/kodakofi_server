package seed

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SeedCodeOtp(ctx context.Context, db *pgxpool.Pool) error {
	query := `
		INSERT INTO public.code_otp (user_id,code,type_id,expired_at,created_at,updated_at) VALUES
			('8a60a062-c0a8-4034-b2f1-a832b9adb435'::uuid,'324681',1,'2025-05-18 20:57:23.063547+07','2025-05-18 20:42:23.064056+07',NULL),
			('171b67f9-a3c2-4115-b75c-1a7816ca1cd2'::uuid,'901956',1,'2025-05-18 21:02:26.914949+07','2025-05-18 20:47:26.91563+07',NULL)
		ON CONFLICT (id) DO NOTHING;
	`

	_, err := db.Exec(ctx, query)
	if err != nil {
		log.Printf("Failed to seed code_otp: %v", err)
		return err
	}

	log.Println("Seeded code_otp successfully.")
	return nil
}
