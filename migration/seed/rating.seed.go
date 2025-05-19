package seed

// import (
// 	"context"
// 	"log"

// 	"github.com/jackc/pgx/v5/pgxpool"
// )

// func SeedRating(ctx context.Context, db *pgxpool.Pool) error {
// 	query := `
// 	INSERT INTO ratings (user_id, product_id, rating)
// 	VALUES
// 	(gen_random_uuid(), '685d83de-a09b-4875-9e21-12d1629e7d61', TRUE),
// 	(gen_random_uuid(), '685d83de-a09b-4875-9e21-12d1629e7d61', TRUE),
// 	(gen_random_uuid(), '685d83de-a09b-4875-9e21-12d1629e7d61', TRUE),
// 	(gen_random_uuid(), '9a283cef-5a9f-437f-ac24-cf489094c7aa', TRUE),
// 	(gen_random_uuid(), '9a283cef-5a9f-437f-ac24-cf489094c7aa', TRUE),
// 	(gen_random_uuid(), '9a283cef-5a9f-437f-ac24-cf489094c7aa', TRUE),
// 	(gen_random_uuid(), '9a283cef-5a9f-437f-ac24-cf489094c7aa', TRUE),
// 	(gen_random_uuid(), '9a283cef-5a9f-437f-ac24-cf489094c7aa', TRUE),
// 	(gen_random_uuid(), 'c4e99f87-efd0-49a3-bd51-e1ba499f5b0b', TRUE),
// 	(gen_random_uuid(), 'c4e99f87-efd0-49a3-bd51-e1ba499f5b0b', TRUE),
// 	(gen_random_uuid(), 'c4e99f87-efd0-49a3-bd51-e1ba499f5b0b', TRUE),
// 	(gen_random_uuid(), 'c4e99f87-efd0-49a3-bd51-e1ba499f5b0b', TRUE),
// 	(gen_random_uuid(), 'c4e99f87-efd0-49a3-bd51-e1ba499f5b0b', TRUE),
// 	(gen_random_uuid(), '00e87080-5fe9-4fbb-a6b4-7f40ac229723', TRUE),
// 	(gen_random_uuid(), '00e87080-5fe9-4fbb-a6b4-7f40ac229723', TRUE),
// 	(gen_random_uuid(), '00e87080-5fe9-4fbb-a6b4-7f40ac229723', TRUE),
// 	(gen_random_uuid(), '00e87080-5fe9-4fbb-a6b4-7f40ac229723', TRUE),
// 	(gen_random_uuid(), 'ff340d2f-3b50-42a8-802f-b11b9b81679d', TRUE),
// 	(gen_random_uuid(), 'ff340d2f-3b50-42a8-802f-b11b9b81679d', TRUE),
// 	(gen_random_uuid(), 'ff340d2f-3b50-42a8-802f-b11b9b81679d', TRUE),
// 	(gen_random_uuid(), 'ff340d2f-3b50-42a8-802f-b11b9b81679d', TRUE),
// 	(gen_random_uuid(), '2bbf7f8a-1830-4d30-a924-4ba81f471e27', TRUE),
// 	(gen_random_uuid(), '2bbf7f8a-1830-4d30-a924-4ba81f471e27', TRUE),
// 	(gen_random_uuid(), '2bbf7f8a-1830-4d30-a924-4ba81f471e27', TRUE),
// 	(gen_random_uuid(), '2bbf7f8a-1830-4d30-a924-4ba81f471e27', TRUE),
// 	(gen_random_uuid(), 'ec4e8c26-b66f-4964-ab5a-4dd8f7ec0e05', TRUE),
// 	(gen_random_uuid(), 'ec4e8c26-b66f-4964-ab5a-4dd8f7ec0e05', TRUE),
// 	(gen_random_uuid(), 'ec4e8c26-b66f-4964-ab5a-4dd8f7ec0e05', TRUE),
// 	(gen_random_uuid(), 'ec4e8c26-b66f-4964-ab5a-4dd8f7ec0e05', TRUE),
// 	(gen_random_uuid(), 'ec4e8c26-b66f-4964-ab5a-4dd8f7ec0e05', TRUE),
// 	(gen_random_uuid(), 'b64c6043-6b04-4e03-aec2-f8534a539de6', TRUE),
// 	(gen_random_uuid(), 'b64c6043-6b04-4e03-aec2-f8534a539de6', TRUE),
// 	(gen_random_uuid(), 'b64c6043-6b04-4e03-aec2-f8534a539de6', TRUE),
// 	(gen_random_uuid(), 'f30b51f9-9877-4e7f-b81a-3406c20b9629', TRUE),
// 	(gen_random_uuid(), 'f30b51f9-9877-4e7f-b81a-3406c20b9629', TRUE),
// 	(gen_random_uuid(), 'f30b51f9-9877-4e7f-b81a-3406c20b9629', TRUE),
// 	(gen_random_uuid(), 'f30b51f9-9877-4e7f-b81a-3406c20b9629', TRUE),
// 	(gen_random_uuid(), 'f30b51f9-9877-4e7f-b81a-3406c20b9629', TRUE)`

// 	_, err := db.Exec(ctx, query)
// 	if err != nil {
// 		log.Printf("Failed to seed rating: %v", err)
// 		return err
// 	}

// 	log.Println("Seeded products rating successfully.")
// 	return nil
// }
