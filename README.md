# Project Evermos API

REST API berbasis Go Fiber + GORM + MySQL dengan JWT Authentication dan dokumentasi Swagger.

## Tech Stack
- Go + Fiber v2
- GORM (MySQL)
- JWT (HS256)
- Swagger (swaggo + gofiber/swagger)
- Migrations SQL (folder `./migrations`)

## Fitur Utama (Modules)
- Auth: login, register
- Users: profil, alamat kirim
- Toko: profil toko, update toko (upload foto)
- Product: CRUD dengan upload foto
- Category: CRUD (admin only)
- Address: list provinces/cities (EMSIFA API + caching)
- Transaction: list, detail, create

## Prasyarat
- Go 1.21+
- MySQL (aktif dan dapat diakses)
- Git
- Swag CLI (untuk generate Swagger): `go install github.com/swaggo/swag/cmd/swag@latest`

## Setup Cepat (Windows)
1) Salin file env
- Copy `.env.example` ke `.env` lalu sesuaikan nilai environment.

2) Instal swag CLI (untuk generate Swagger)
- Pastikan `GOBIN` sudah ada di PATH Windows Anda.

3) Konfigurasi .env
- Wajib isi:
  - DB_HOST, DB_PORT, DB_USER, DB_PASS, DB_NAME
  - JWT_SECRET (gunakan string acak yang panjang dan kuat)
  - APP_PORT (misal: 8080)
  - BASE_FILE_URL (misal: http://127.0.0.1:8080)
  - UPLOAD_DIR_PRODUCT (default: uploads/products)

Contoh `.env` minimal:
```
APP_PORT=8080
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASS=yourpassword
DB_NAME=evermos
JWT_SECRET=your-very-strong-random-secret
BASE_FILE_URL=http://127.0.0.1:8080
UPLOAD_DIR_PRODUCT=uploads/products
```

4) Generate dokumentasi Swagger
```bash
swag init -g cmd/main.go -o docs
```

5) Jalankan aplikasi
```bash
go run ./cmd/main.go
```
- Proses migrasi database akan berjalan otomatis saat start.

## Perintah Umum
- Install swag CLI:
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```
- Generate Swagger:
```bash
swag init -g cmd/main.go -o docs
```
- Run server:
```bash
go run ./cmd/main.go
```
- Build binary:
```bash
go build -o evermos.exe ./cmd/main.go
```

## URL Penting
- Swagger UI: http://127.0.0.1:8080/swagger/index.html
- Health Check: http://127.0.0.1:8080/health

## Autentikasi JWT
- Login menghasilkan token JWT (HS256) ditandatangani dengan `JWT_SECRET`.
- Kirimkan token pada setiap request:
  - Header utama: `token: <JWT>`
  - Alternatif: `Authorization: Bearer <JWT>`

Catatan Keamanan:
- Jaga kerahasiaan `JWT_SECRET`.
- Gunakan nilai random yang kuat, jangan hardcode atau nilai lemah.
- Ganti sekret jika terindikasi bocor.

## Upload Files
- Product photos: `POST /product` (multipart form, field `photos`)
- Update toko dengan foto: `PUT /toko/{id_toko}` (multipart form, field `photo`)
- File disimpan di folder `./uploads` (URL publik bergantung `BASE_FILE_URL`).

## Database & Migrasi
- File migrasi ada di folder `./migrations`.
- Migrasi dijalankan otomatis saat server start.
- Pastikan database sudah ada dan kredensial benar di `.env`.

## Menjalankan Server Kedua (Port Berbeda)
- Aplikasi membaca `APP_PORT` dari `.env`.
- Jika ingin menjalankan instance kedua, ubah sementara `APP_PORT` (mis. 8081) lalu jalankan instance kedua.
- Alternatif: jalankan salinan proyek dengan `.env` berbeda untuk berjalan bersamaan.

## Struktur Direktori (ringkas)
- `cmd/main.go` — entrypoint aplikasi
- `api/http` — router Fiber dan anotasi Swagger
- `internal/config` — loader konfigurasi .env
- `internal/db` — koneksi DB dan migrasi
- `internal/todo` — handlers, services, repositories, models
- `migrations` — file SQL migrasi
- `docs` — hasil generate Swagger (docs.go/json/yaml)
- `uploads` — penyimpanan file

## Arsitektur Singkat
Diagram (Mermaid):
```mermaid
flowchart TD
    Client[Clients\n(Web / Mobile / Postman)] -->|HTTP/JSON| Fiber[Go Fiber HTTP Server]

    subgraph Fiber Layer
        Router[Router + Swagger UI]
        MW[Middleware\n- JWT Auth\n- Logging\n- Recovery]
    end

    Fiber --> Router
    Router --> MW
    MW --> Handlers[Handlers (api/http/*)]
    Handlers --> Services[Services (business logic)]
    Services --> Repos[Repositories (DB access)]
    Repos --> MySQL[(MySQL)]
    Services --> EMSIFA[External API: EMSIFA Provinces/Cities]
    Services --> Files[File Storage\n(uploads/)]
    Config[Config (.env -> internal/config)] --> Fiber

    %% Envelope response standar
    %% status, message, errors, data
```

Fallback ASCII (jika Mermaid tidak tersedia):
```
Clients -> Fiber HTTP Server
           -> Router + Swagger UI
              -> Middleware (JWT, Logging, Recovery)
                 -> Handlers (api/http/*)
                    -> Services (business logic)
                       -> Repositories -> MySQL (DB)
                       -> External API (EMSIFA)
                       -> File Storage (uploads/)
Config (.env -> internal/config) -> Fiber init
```

Alur singkat request:
- Client mengirim HTTP/JSON ke server Fiber.
- Router menyalurkan ke route; Swagger UI tersedia untuk dokumentasi.
- Middleware memverifikasi JWT, logging, dan error recovery.
- Handler parsing/validasi input, memanggil Service.
- Service menjalankan bisnis proses, orkestrasi transaksi, dan memanggil Repository/External API/File storage sesuai kebutuhan.
- Repository mengakses MySQL via GORM.
- Response dikembalikan dalam amplop standar: `status`, `message`, `errors`, `data`.

## Troubleshooting
- “missing required JWT env var”: set `JWT_SECRET` di `.env`.
- Koneksi DB gagal: cek host, port, user, password, dan DB sudah dibuat.
- Swagger tidak update: jalankan `swag init -g cmd/main.go -o docs`.
- Port bentrok: ubah `APP_PORT` di `.env`.

## Lisensi
Internal use.