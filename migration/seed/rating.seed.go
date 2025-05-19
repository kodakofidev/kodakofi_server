package seed

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SeedRating(ctx context.Context, db *pgxpool.Pool) error {
	query := `
	INSERT INTO ratings (user_id, product_id, rating)
	VALUES
	('8a60a062-c0a8-4034-b2f1-a832b9adb435', '5740f5fe-5178-4933-9e7b-63cab72fd79a', 1),
	('171b67f9-a3c2-4115-b75c-1a7816ca1cd2', '5740f5fe-5178-4933-9e7b-63cab72fd79a', 1),
	('8a60a062-c0a8-4034-b2f1-a832b9adb435', '06bd407e-5fab-4590-9e3f-7dc442af3b42', 1),
	('171b67f9-a3c2-4115-b75c-1a7816ca1cd2', '06bd407e-5fab-4590-9e3f-7dc442af3b42', 1),
	('171b67f9-a3c2-4115-b75c-1a7816ca1cd2', '31b9935a-7bd3-4ae0-898d-386e4cffb82e', 1),
	('8a60a062-c0a8-4034-b2f1-a832b9adb435', '31b9935a-7bd3-4ae0-898d-386e4cffb82e', 1),
	('171b67f9-a3c2-4115-b75c-1a7816ca1cd2', '287ec09d-928c-4562-9f29-86ad95dce6f6', 1),
	('8a60a062-c0a8-4034-b2f1-a832b9adb435', '287ec09d-928c-4562-9f29-86ad95dce6f6', 1),
	('171b67f9-a3c2-4115-b75c-1a7816ca1cd2', '0076aee4-9db2-4941-a69d-e07ff562dc3b', 1),
	('8a60a062-c0a8-4034-b2f1-a832b9adb435', '0076aee4-9db2-4941-a69d-e07ff562dc3b', 1),
	('171b67f9-a3c2-4115-b75c-1a7816ca1cd2', 'dfd726cf-8c5c-44ac-8d06-82378bb4c31c', 1),
	('8a60a062-c0a8-4034-b2f1-a832b9adb435', 'dfd726cf-8c5c-44ac-8d06-82378bb4c31c', 1),
	('171b67f9-a3c2-4115-b75c-1a7816ca1cd2', '9a7b950f-3664-4df3-9da5-249e91de2b31', 1),
	('8a60a062-c0a8-4034-b2f1-a832b9adb435', '9a7b950f-3664-4df3-9da5-249e91de2b31', 1),
	('171b67f9-a3c2-4115-b75c-1a7816ca1cd2', '7af264d9-fa31-45a4-8948-b2db4c267fd6', 1),
	('8a60a062-c0a8-4034-b2f1-a832b9adb435', '7af264d9-fa31-45a4-8948-b2db4c267fd6', 1),
	('171b67f9-a3c2-4115-b75c-1a7816ca1cd2', '9c4817f0-455e-415f-9c79-8a7e6f4fc1ab', 1),
	('8a60a062-c0a8-4034-b2f1-a832b9adb435', '9c4817f0-455e-415f-9c79-8a7e6f4fc1ab', 1)
	`

	_, err := db.Exec(ctx, query)
	if err != nil {
		log.Printf("Failed to seed rating: %v", err)
		return err
	}

	log.Println("Seeded products rating successfully.")
	return nil
}
