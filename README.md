# ☕️ Koda Kofi

## 📘 Database Setup — Migration & Seeding

Proyek ini menggunakan PostgreSQL untuk manajemen data. Terdapat dua komponen utama yang perlu dijalankan:

- **Migration** – untuk membuat/mengubah struktur tabel di database
- **Seeding** – untuk mengisi data awal ke dalam tabel

---

## ⚙️ Prasyarat

1. Pastikan sudah ada file `.env` di root project yang berisi:

    ```env
    DB_URL=postgres://username:password@localhost:5432/db_name?sslmode=disable
    ```

2. Install CLI `golang-migrate`:

    ### 🛠️ Cara install via `go install` (direkomendasikan)

    ```bash
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
    ```

    > 🔔 Pastikan `$GOBIN` sudah masuk ke `PATH`. Tambahkan ke `.bashrc` / `.zshrc` jika belum:

    bash

    ```bash
    export PATH=$PATH:$(go env GOPATH)/bin
    ```
    atau dengan Powershell
    ```
    $env:PATH += ";$(go env GOPATH)\bin"
    ```


---

## 🔨 Migration

### 📁 Struktur Folder

```
migration/
├── 20250516120000_create_users_table.up.sql
├── 20250516120000_create_users_table.down.sql
└── ...
```

### 🛠️ Perintah Makefile untuk Migration

- **Buat file migration baru**:

    ```bash
    make migrate-init name=create_table_name
    ```

- **Jalankan seluruh migration ke atas**:

    ```bash
    make migrate-up
    ```

- **Rollback migration (turun 1 file)**:

    ```bash
    make migrate-down
    ```
    Membatalkan/mengembalikan perubahan migrasi terakhir yang telah diterapkan.

- **Reset status migration (force 0)**:

    ```bash
    make migrate-fix
    ```

---

## 🌱 Seeding

### 📁 Struktur Folder

```
cmd/
└── seeder/
└── seed.main.go // Entry point seeder

migration/
└── seed/
├── delivery_methods.seed.go
├── payment_methods.seed.go
└── status.seed.go
```

### ▶️ Jalankan Seeding

- Gunakan Makefile:

    ```bash
    make seed
    ```

- Atau jalankan langsung:

    ```bash
    go run ./cmd/seeder/seed.main.go
    ```

Seeder akan mengisi data awal untuk:

- Delivery Methods
- Payment Methods
- Status

### ✅ Output Berhasil

Contoh output jika berhasil:

```
Starting delivery_methods seeding...
Seeded delivery_methods successfully.
Starting payment_methods seeding...
Seeded payment_methods successfully.
Starting status seeding...
Seeded status successfully.
Seeding completed successfully.
```

### ✅ Migrate Ulang + Seeding

```bash
make migrate-reset
```
Perintah tersebut akan melakukan migrate ulang sekaligus menjalankan seeding

---

## 📦 Tips

- Buat dahulu database di local

    ```CREATE DASATABE <database_name>;```
- Jalankan `migrate-up` terlebih dahulu sebelum menjalankan `seed`
- Seeder aman untuk dijalankan berulang karena menggunakan

    `ON CONFLICT DO NOTHING`
- Semua koneksi database menggunakan `pgxpool` (bukan `sqlx` atau `database/sql`)
- Jika sudah pernah melakukan migrate, maka gunakan:

    ```make migrate-reset```

    untuk migrate ulang dan langsung menjalankan seeding