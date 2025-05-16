package seed

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SeedOtpType(ctx context.Context, db *pgxpool.Pool) error {
	query := `
		INSERT INTO otp_type (name)
		VALUES 
			('email_verification'),
			('password_reset')
		ON CONFLICT (name) DO NOTHING;
	`

	_, err := db.Exec(ctx, query)
	if err != nil {
		log.Printf("Failed to seed otp_type: %v", err)
		return err
	}

	log.Println("Seeded otp_type successfully.")
	return nil
}
