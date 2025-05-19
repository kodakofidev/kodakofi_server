package seed

// import (
// 	"context"
// 	"log"

// 	"github.com/jackc/pgx/v5/pgxpool"
// )

// func SeedImage(ctx context.Context, db *pgxpool.Pool) error {
// 	query := `
// 	INSERT INTO product_images (product_id, path)
// 	VALUES
// 	('685d83de-a09b-4875-9e21-12d1629e7d61', '/public/product-image/image5.jpg'),
// 	('685d83de-a09b-4875-9e21-12d1629e7d61', '/public/product-image/image2.jpg'),
// 	('685d83de-a09b-4875-9e21-12d1629e7d61', '/public/product-image/image7.jpg'),
// 	('ff340d2f-3b50-42a8-802f-b11b9b81679d', '/public/product-image/image5.jpg'),
// 	('ff340d2f-3b50-42a8-802f-b11b9b81679d', '/public/product-image/image2.jpg'),
// 	('ff340d2f-3b50-42a8-802f-b11b9b81679d', '/public/product-image/image7.jpg'),
// 	('7bc51c9b-f372-4bd0-bccf-9ed732d7c335', '/public/product-image/image5.jpg'),
// 	('7bc51c9b-f372-4bd0-bccf-9ed732d7c335', '/public/product-image/image2.jpg'),
// 	('7bc51c9b-f372-4bd0-bccf-9ed732d7c335', '/public/product-image/image7.jpg'),
// 	('1607a577-6df0-4974-9b93-6dc329cc8368', '/public/product-image/image5.jpg'),
// 	('1607a577-6df0-4974-9b93-6dc329cc8368', '/public/product-image/image2.jpg'),
// 	('1607a577-6df0-4974-9b93-6dc329cc8368', '/public/product-image/image7.jpg'),
// 	('791f7fe7-3952-4724-895a-92f914fec372', '/public/product-image/image5.jpg'),
// 	('791f7fe7-3952-4724-895a-92f914fec372', '/public/product-image/image2.jpg'),
// 	('791f7fe7-3952-4724-895a-92f914fec372', '/public/product-image/image7.jpg'),
// 	('2bbf7f8a-1830-4d30-a924-4ba81f471e27', '/public/product-image/image3.png'),
// 	('2bbf7f8a-1830-4d30-a924-4ba81f471e27', '/public/product-image/image4.png'),
// 	('2bbf7f8a-1830-4d30-a924-4ba81f471e27', '/public/product-image/image6.jpg'),
// 	('ca0905eb-2088-4a11-af8f-7b76601cc206', '/public/product-image/image3.png'),
// 	('ca0905eb-2088-4a11-af8f-7b76601cc206', '/public/product-image/image4.png'),
// 	('ca0905eb-2088-4a11-af8f-7b76601cc206', '/public/product-image/image6.jpg'),
// 	('c4e99f87-efd0-49a3-bd51-e1ba499f5b0b', '/public/product-image/image3.png'),
// 	('c4e99f87-efd0-49a3-bd51-e1ba499f5b0b', '/public/product-image/image4.png'),
// 	('c4e99f87-efd0-49a3-bd51-e1ba499f5b0b', '/public/product-image/image6.jpg'),
// 	('e58daf7e-05e4-4353-84c1-471b46718b05', '/public/product-image/image3.png'),
// 	('e58daf7e-05e4-4353-84c1-471b46718b05', '/public/product-image/image4.png'),
// 	('e58daf7e-05e4-4353-84c1-471b46718b05', '/public/product-image/image6.jpg'),
// 	('839f121f-8351-4e2c-8e10-d8ac81cc3e45', '/public/product-image/image3.png'),
// 	('839f121f-8351-4e2c-8e10-d8ac81cc3e45', '/public/product-image/image4.png'),
// 	('839f121f-8351-4e2c-8e10-d8ac81cc3e45', '/public/product-image/image6.jpg'),
// 	('9a283cef-5a9f-437f-ac24-cf489094c7aa', '/public/product-image/image1.png'),
// 	('9a283cef-5a9f-437f-ac24-cf489094c7aa', '/public/product-image/image3.png'),
// 	('9a283cef-5a9f-437f-ac24-cf489094c7aa', '/public/product-image/image4.png'),
// 	('47b80d70-0512-4c11-9e3f-364563b40e4a', '/public/product-image/image1.png'),
// 	('47b80d70-0512-4c11-9e3f-364563b40e4a', '/public/product-image/image3.png'),
// 	('47b80d70-0512-4c11-9e3f-364563b40e4a', '/public/product-image/image4.png')
// 	`

// 	_, err := db.Exec(ctx, query)
// 	if err != nil {
// 		log.Printf("Failed to seed image: %v", err)
// 		return err
// 	}

// 	log.Println("Seeded product image successfully.")
// 	return nil
// }
