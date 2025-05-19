package seed

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SeedUsers(ctx context.Context, db *pgxpool.Pool) error {
	query := `
		INSERT INTO public.users (id, email, "password", "role", is_verified, created_at, updated_at) VALUES
			('c59b7e07-cbb6-4647-858a-406c3c4b4bd2'::uuid, 'xobeba1095@magpit.com', '$argon2id$v=19$m=65536,t=3,p=2$MH9T37oskorZGmHyuZ5Uyg$yyCmvc+buc5croTd4TmxLtqjaOd1D3LvXWOnyMl8AQ0', 'admin', true, '2025-05-18 18:28:38.854598+07', NULL),
			('171b67f9-a3c2-4115-b75c-1a7816ca1cd2'::uuid, 'sdgi6@dcpa.net', '$argon2id$v=19$m=65536,t=3,p=2$jEpY2BHM0zseUt//j13fNQ$1ORFua0u/0YCfj47Olmuf+t27DEcWPqGB0KvlGCM+bo', 'user', true, '2025-05-18 20:47:26.90677+07', NULL),
			('0ee41132-5592-4626-8f87-1df70db54e3f'::uuid, 'testuser@example.com', '$argon2id$v=19$m=65536,t=3,p=2$abc123$xyz4567890', 'user', true, '2025-05-19 03:49:14.898142+07', NULL),
			('8a60a062-c0a8-4034-b2f1-a832b9adb435'::uuid, 'i08qo@dcpa.net', '$argon2id$v=19$m=65536,t=3,p=2$Nw3mnvc4JCgzbYfctP5WeA$CQJ2RCFXV5SRtQCzAwP3G7wSPkUCAOYrcOR25Ubk7sU', 'user', true, '2025-05-18 20:42:23.05615+07', NULL)
		ON CONFLICT (id) DO NOTHING;
	`

	_, err := db.Exec(ctx, query)
	if err != nil {
		log.Printf("Failed to seed users: %v", err)
		return err
	}

	log.Println("Seeded users successfully.")
	return nil
}
