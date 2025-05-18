package seed

// import (
// 	"context"
// 	"log"

// 	"github.com/jackc/pgx/v5/pgxpool"
// )

// func SeedSize(ctx context.Context, db *pgxpool.Pool) error {
// 	query := `
// 	INSERT INTO sizes (size, added_price)
// 	VALUES
// 	('Regular', 0),
// 	('Medium', 3000),
// 	('Large', 5000)
// 	`
// 	_, err := db.Exec(ctx, query)
// 	if err != nil {
// 		log.Printf("Failed to seed size: %v", err)
// 		return err
// 	}

// 	log.Println("Seeded size successfully.")
// 	return nil
// }
