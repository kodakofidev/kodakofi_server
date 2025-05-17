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

	log.Println("Starting delivery_methods seeding...")
	seed.SeedDeliveryMethods(ctx, db)

	log.Println("Starting payment_methods seeding...")
	seed.SeedPaymentMethods(ctx, db)

	log.Println("Starting status seeding...")
	seed.SeedStatus(ctx, db)

	log.Println("Starting otp_type seeding...")
	seed.SeedOtpType(ctx, db)

	log.Println("Starting categories seeding...")
	seed.SeedCategories(ctx, db)

	log.Println("Seeding completed successfully.")
}
