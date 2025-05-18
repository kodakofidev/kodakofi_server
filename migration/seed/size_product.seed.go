package seed

// import (
// 	"context"
// 	"log"

// 	"github.com/jackc/pgx/v5/pgxpool"
// )

// func SeedSizeProduct(ctx context.Context, db *pgxpool.Pool) error {
// 	query := `
// 	INSERT INTO size_products (product_id, size_id)
// 	VALUES
// 	('685d83de-a09b-4875-9e21-12d1629e7d61', 1),
// 	('685d83de-a09b-4875-9e21-12d1629e7d61', 2),
// 	('685d83de-a09b-4875-9e21-12d1629e7d61', 3),
// 	('ff340d2f-3b50-42a8-802f-b11b9b81679d', 1),
// 	('ff340d2f-3b50-42a8-802f-b11b9b81679d', 2),
// 	('ff340d2f-3b50-42a8-802f-b11b9b81679d', 3),
// 	('7bc51c9b-f372-4bd0-bccf-9ed732d7c335', 1),
// 	('7bc51c9b-f372-4bd0-bccf-9ed732d7c335', 2),
// 	('7bc51c9b-f372-4bd0-bccf-9ed732d7c335', 3),
// 	('1607a577-6df0-4974-9b93-6dc329cc8368', 1),
// 	('1607a577-6df0-4974-9b93-6dc329cc8368', 2),
// 	('1607a577-6df0-4974-9b93-6dc329cc8368', 3),
// 	('791f7fe7-3952-4724-895a-92f914fec372', 1),
// 	('791f7fe7-3952-4724-895a-92f914fec372', 2),
// 	('791f7fe7-3952-4724-895a-92f914fec372', 3),
// 	('2bbf7f8a-1830-4d30-a924-4ba81f471e27', 1),
// 	('2bbf7f8a-1830-4d30-a924-4ba81f471e27', 2),
// 	('2bbf7f8a-1830-4d30-a924-4ba81f471e27', 3),
// 	('ca0905eb-2088-4a11-af8f-7b76601cc206', 1),
// 	('ca0905eb-2088-4a11-af8f-7b76601cc206', 2),
// 	('ca0905eb-2088-4a11-af8f-7b76601cc206', 3),
// 	('e58daf7e-05e4-4353-84c1-471b46718b05', 1),
// 	('e58daf7e-05e4-4353-84c1-471b46718b05', 2),
// 	('e58daf7e-05e4-4353-84c1-471b46718b05', 3),
// 	('839f121f-8351-4e2c-8e10-d8ac81cc3e45', 1),
// 	('839f121f-8351-4e2c-8e10-d8ac81cc3e45', 2),
// 	('839f121f-8351-4e2c-8e10-d8ac81cc3e45', 3)
// 	`
// 	_, err := db.Exec(ctx, query)
// 	if err != nil {
// 		log.Printf("Failed to seed product size: %v", err)
// 		return err
// 	}

// 	log.Println("Seeded product size successfully.")
// 	return nil
// }
