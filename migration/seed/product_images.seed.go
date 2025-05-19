package seed

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SeedProductImages(ctx context.Context, db *pgxpool.Pool) error {
	query := `
		INSERT INTO product_images (product_id, path)
		VALUES
			('4e841656-596c-434d-b5bb-1f27f5d7418c', '1_1.jpg'),
			('4e841656-596c-434d-b5bb-1f27f5d7418c', '1_2.jpg'),
			('4e841656-596c-434d-b5bb-1f27f5d7418c', '1_3.jpg'),
			('5740f5fe-5178-4933-9e7b-63cab72fd79a', '2_1.jpg'),
			('5740f5fe-5178-4933-9e7b-63cab72fd79a', '2_2.jpg'),
			('5740f5fe-5178-4933-9e7b-63cab72fd79a', '2_3.jpg'),
			('e6622f29-2a8b-4110-9a82-8616eed29570', '3_1.jpg'),
			('e6622f29-2a8b-4110-9a82-8616eed29570', '3_2.jpg'),
			('e6622f29-2a8b-4110-9a82-8616eed29570', '3_3.jpg'),
			('9a7b950f-3664-4df3-9da5-249e91de2b31', '4_1.jpg'),
			('9a7b950f-3664-4df3-9da5-249e91de2b31', '4_2.jpg'),
			('9a7b950f-3664-4df3-9da5-249e91de2b31', '4_3.jpg'),
			('8afc2e72-ef45-45b0-a936-62d34bd626bf', '5_1.jpg'),
			('8afc2e72-ef45-45b0-a936-62d34bd626bf', '5_2.jpg'),
			('8afc2e72-ef45-45b0-a936-62d34bd626bf', '5_3.jpg'),
			('7af264d9-fa31-45a4-8948-b2db4c267fd6', '6_1.jpg'),
			('7af264d9-fa31-45a4-8948-b2db4c267fd6', '6_2.jpg'),
			('7af264d9-fa31-45a4-8948-b2db4c267fd6', '6_3.jpg'),
			('f866f4f6-f89c-4395-90ca-241dfb52951c', '7_1.jpg'),
			('f866f4f6-f89c-4395-90ca-241dfb52951c', '7_2.jpg'),
			('f866f4f6-f89c-4395-90ca-241dfb52951c', '7_3.jpg'),
			('287ec09d-928c-4562-9f29-86ad95dce6f6', '8_1.jpg'),
			('287ec09d-928c-4562-9f29-86ad95dce6f6', '8_2.jpg'),
			('287ec09d-928c-4562-9f29-86ad95dce6f6', '8_3.jpg'),
			('95a70a1a-10c8-4a3f-80ad-6c430b74ef3e', '9_1.jpg'),
			('95a70a1a-10c8-4a3f-80ad-6c430b74ef3e', '9_2.jpg'),
			('95a70a1a-10c8-4a3f-80ad-6c430b74ef3e', '9_3.jpg'),
			('40425c10-f932-4b44-97a6-681b56a5ddfa', '10_1.jpg'),
			('40425c10-f932-4b44-97a6-681b56a5ddfa', '10_2.jpg'),
			('40425c10-f932-4b44-97a6-681b56a5ddfa', '10_3.jpg'),
			('0076aee4-9db2-4941-a69d-e07ff562dc3b', '11_1.jpg'),
			('0076aee4-9db2-4941-a69d-e07ff562dc3b', '11_2.jpg'),
			('0076aee4-9db2-4941-a69d-e07ff562dc3b', '11_3.jpg'),
			('aead3bb8-eaa3-4408-a192-cc36b227f464', '12_1.jpg'),
			('aead3bb8-eaa3-4408-a192-cc36b227f464', '12_2.jpg'),
			('aead3bb8-eaa3-4408-a192-cc36b227f464', '12_3.jpg'),
			('a2b74af4-06cc-4004-989e-6150af06926c', '13_1.jpg'),
			('a2b74af4-06cc-4004-989e-6150af06926c', '13_2.jpg'),
			('a2b74af4-06cc-4004-989e-6150af06926c', '13_3.jpg'),
			('dfd726cf-8c5c-44ac-8d06-82378bb4c31c', '14_1.jpg'),
			('dfd726cf-8c5c-44ac-8d06-82378bb4c31c', '14_2.jpg'),
			('dfd726cf-8c5c-44ac-8d06-82378bb4c31c', '14_3.jpg'),
			('31b9935a-7bd3-4ae0-898d-386e4cffb82e', '15_1.jpg'),
			('31b9935a-7bd3-4ae0-898d-386e4cffb82e', '15_2.jpg'),
			('31b9935a-7bd3-4ae0-898d-386e4cffb82e', '15_3.jpg'),
			('06bd407e-5fab-4590-9e3f-7dc442af3b42', '16_1.jpg'),
			('06bd407e-5fab-4590-9e3f-7dc442af3b42', '16_2.jpg'),
			('06bd407e-5fab-4590-9e3f-7dc442af3b42', '16_3.jpg'),
			('978feadd-1c68-479f-a99a-b831b732464b', '17_1.jpg'),
			('978feadd-1c68-479f-a99a-b831b732464b', '17_2.jpg'),
			('978feadd-1c68-479f-a99a-b831b732464b', '17_3.jpg'),
			('fad473ac-7dbe-470a-af07-836726e9b1a6', '18_1.jpg'),
			('fad473ac-7dbe-470a-af07-836726e9b1a6', '18_2.jpg'),
			('fad473ac-7dbe-470a-af07-836726e9b1a6', '18_3.jpg'),
			('e07cf18a-b204-455d-95e3-eaf489a805b2', '19_1.jpg'),
			('e07cf18a-b204-455d-95e3-eaf489a805b2', '19_2.jpg'),
			('e07cf18a-b204-455d-95e3-eaf489a805b2', '19_3.jpg'),
			('9c4817f0-455e-415f-9c79-8a7e6f4fc1ab', '20_1.jpg'),
			('9c4817f0-455e-415f-9c79-8a7e6f4fc1ab', '20_2.jpg'),
			('9c4817f0-455e-415f-9c79-8a7e6f4fc1ab', '20_3.jpg'),
			('40cb38bd-ba7f-4d40-b5e5-0a013b45e4c8', '21_1.jpg'),
			('40cb38bd-ba7f-4d40-b5e5-0a013b45e4c8', '21_2.jpg'),
			('40cb38bd-ba7f-4d40-b5e5-0a013b45e4c8', '21_3.jpg'),
			('97db8997-77ad-4793-8a83-e030ff84d4dd', '22_1.jpg'),
			('97db8997-77ad-4793-8a83-e030ff84d4dd', '22_2.jpg'),
			('97db8997-77ad-4793-8a83-e030ff84d4dd', '22_3.jpg'),
			('d6a142e4-9960-444f-9cca-041ceb595a9a', '23_1.jpg'),
			('d6a142e4-9960-444f-9cca-041ceb595a9a', '23_2.jpg'),
			('d6a142e4-9960-444f-9cca-041ceb595a9a', '23_3.jpg'),
			('9577276b-d4cb-4de2-a371-c4a25e790c63', '24_1.jpg'),
			('9577276b-d4cb-4de2-a371-c4a25e790c63', '24_2.jpg'),
			('9577276b-d4cb-4de2-a371-c4a25e790c63', '24_3.jpg'),
			('042875a5-da48-456a-a6e7-b6746f43ab02', '25_1.jpg'),
			('042875a5-da48-456a-a6e7-b6746f43ab02', '25_2.jpg'),
			('042875a5-da48-456a-a6e7-b6746f43ab02', '25_3.jpg'),
			('fb374b9a-ddde-4a0a-8739-325dcf9543dc', '26_1.jpg'),
			('fb374b9a-ddde-4a0a-8739-325dcf9543dc', '26_2.jpg'),
			('fb374b9a-ddde-4a0a-8739-325dcf9543dc', '26_3.jpg'),
			('ed6d3a40-7760-45cd-a8df-106daeef0227', '27_1.jpg'),
			('ed6d3a40-7760-45cd-a8df-106daeef0227', '27_2.jpg'),
			('ed6d3a40-7760-45cd-a8df-106daeef0227', '27_3.jpg'),
			('c504c3ea-af18-4bd8-a2e9-d7e773b1ea5d', '28_1.jpg'),
			('c504c3ea-af18-4bd8-a2e9-d7e773b1ea5d', '28_2.jpg'),
			('c504c3ea-af18-4bd8-a2e9-d7e773b1ea5d', '28_3.jpg'),
			('986f09ea-9843-45a3-91a2-da3193c12a63', '29_1.jpg'),
			('986f09ea-9843-45a3-91a2-da3193c12a63', '29_2.jpg'),
			('986f09ea-9843-45a3-91a2-da3193c12a63', '29_3.jpg'),
			('4c8c4d2e-6f9a-4c91-9f5d-0660ce5c3e6e', '30_1.jpg'),
			('4c8c4d2e-6f9a-4c91-9f5d-0660ce5c3e6e', '30_2.jpg'),
			('4c8c4d2e-6f9a-4c91-9f5d-0660ce5c3e6e', '30_3.jpg'),
			('8511cb06-1612-4b5a-9ef5-143abdc2077a', '31_1.jpg'),
			('8511cb06-1612-4b5a-9ef5-143abdc2077a', '31_2.jpg'),
			('8511cb06-1612-4b5a-9ef5-143abdc2077a', '31_3.jpg'),
			('40a2a975-9223-4a0e-9379-cb055fe8ae98', '32_1.jpg'),
			('40a2a975-9223-4a0e-9379-cb055fe8ae98', '32_2.jpg'),
			('40a2a975-9223-4a0e-9379-cb055fe8ae98', '32_3.jpg')
		ON CONFLICT (product_id, path) DO NOTHING;
	`

	_, err := db.Exec(ctx, query)
	if err != nil {
		log.Printf("Failed to seed product_images: %v", err)
		return err
	}

	log.Println("Seeded product_images successfully.")
	return nil
}
