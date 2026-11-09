# dc-transaction-service
- service config get value from vault and put on file named application_env (e.g. application_local)
- install swagger on local machine : make install-swag
- install mockery on local machine : make install-mockery
- download go dependencies : go mod tidy
- generate docs : make docs or swag init

## Setup Database (Migration & Seeding)

Service memakai 1 database Postgres: `warehouse_db` untuk semua
operasi (read & write, tanpa replica).
Migrasi/seeding dijalankan lewat subcommand di `main.go` — tidak perlu
eksekusi SQL manual ke DB.

### 1. Siapkan config

```bash
cp template_application_local.json application_local.json
```

Isi di `application_local.json`:

- `database.warehouse_db` →
  `username`, `password`, `hostname` (`host:port`), `dbname`
- `secretkey` → base64 string (dipakai signing JWT)

`ReadConfig` membaca file `application_<env>.json` sesuai flag `-env`
(mis. `-env local` → `application_local.json`). `InitDBConnection`
membuka koneksi `pgx` ke tiap key di map `database` dengan DSN
`postgres://<username>:<password>@<hostname>/<dbname>`.

### 2. Jalankan migration

```bash
# lihat status migrasi (applied/pending)
go run main.go -env local migrate status

# apply semua migrasi yang pending
go run main.go -env local migrate up

# rollback — default 1 step, atau n step
go run main.go -env local migrate down
go run main.go -env local migrate down 2
```

Yang dibuat migrasi saat ini (folder `migration/`):

| Version | Name | Isi |
|---|---|---|
| 202611090001 | `create_user_schema_tables` | schema `"User"` + `"User".roles` + `"User".users` |

Tiap migrasi tercatat di tabel `schema_migrations` (auto-create) di DB
target, jadi `migrate up` aman dijalankan berulang — yang sudah applied
di-skip. Target DB ditentukan field `DBKey` di tiap file migrasi
(selalu writer, bukan replica).

Menambah tabel/kolom baru = buat file `migration/<version>_<name>.go`,
register `Migration{Version, Name, DBKey, Up, Down}` via `init()`.
Pakai timestamp-style version, mis. `202611100001`.

### 3. Jalankan seeding

```bash
# semua seeder
go run main.go -env local seed

# satu seeder tertentu
go run main.go -env local seed roles
go run main.go -env local seed users
```

Seeder yang tersedia (folder `seed/`, idempotent — aman diulang):

| Name | Isi |
|---|---|
| `roles` | `Superadmin`(1), `Admin`(2), `Staff`(3) → `"User".roles` |
| `users` | `superadmin`/`admin123`, `admin`/`admin123`, `staff`/`staff123` → `"User".users` (password SHA256) |

Catatan: jalankan `roles` sebelum `users` (FK `role_id`). `seed` tanpa
argumen menjalankan semua seeder terurut by nama — `roles` otomatis
lebih dulu.

### 4. Start service

```bash
go run main.go -env local
```

Shortcut Makefile: `make migrate-up`, `make migrate-down STEPS=n`,
`make migrate-status`, `make seed NAME=roles` (`ENV=local` default).

## Testing

Unit test memakai stub hand-written (tanpa DB/redis/network), contoh:
`app/service/auth_service_test.go`.

```bash
# semua test di package service
go test ./app/service/

# verbose — tampilkan tiap test & subtest
go test ./app/service/ -v

# semua test di project
go test ./...

# hanya test tertentu (regex)
go test ./app/service/ -run TestAuth -v

# hanya satu subtest
go test ./app/service/ -run 'TestAuthLogin_Success' -v

# coverage
go test ./app/service/ -coverprofile=cover.out && go tool cover -func=cover.out
```

Catatan: `TestMain` meng-init logger ke `logs/*.log` — pastikan folder `logs/` ada.
