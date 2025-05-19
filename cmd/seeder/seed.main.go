package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"
	"github.com/kodakofidev/kodakofi_server/migration/seed"
	"github.com/kodakofidev/kodakofi_server/pkg"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	ctx := context.Background()

	db, err := pkg.Posql()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	log.Println("Starting users seeding...")
	seed.SeedUsers(ctx, db)

	log.Println("Starting otp_type seeding...")
	seed.SeedOtpType(ctx, db)

	log.Println("Starting code_otp seeding...")
	seed.SeedCodeOtp(ctx, db)

	log.Println("Starting profiles seeding...")
	seed.SeedProfiles(ctx, db)

	log.Println("Starting categories seeding...")
	seed.SeedCategories(ctx, db)

	log.Println("Starting products seeding...")
	seed.SeedProducts(ctx, db)

	log.Println("Starting sizes seeding...")
	seed.SeedSizes(ctx, db)

	log.Println("Starting size_products seeding...")
	seed.SeedSizeProducts(ctx, db)

	log.Println("Starting delivery_methods seeding...")
	seed.SeedDeliveryMethods(ctx, db)

	log.Println("Starting payment_methods seeding...")
	seed.SeedPaymentMethods(ctx, db)

	log.Println("Starting status seeding...")
	seed.SeedStatus(ctx, db)

	// log.Println("Starting orders seeding...")
	// seed.SeedOrders(ctx, db)

	// log.Println("Starting transactions seeding...")
	// seed.SeedTrasnsactions(ctx, db)

	// log.Println("Starting products_orders seeding...")
	// seed.SeedPrductsOrders(ctx, db)

	log.Println("Starting tax seeding...")
	seed.SeedTax(ctx, db)

	log.Println("Starting products rating seeding...")
	seed.SeedRating(ctx, db)

	// log.Println("Starting product image seeding...")
	// seed.SeedImage(ctx, db)

	log.Println("Starting discount seeding...")
	seed.SeedDiscount(ctx, db)

	log.Println("Starting product discount seeding...")
	seed.SeedProductDiscount(ctx, db)

	log.Println("Seeding completed successfully.")
}
