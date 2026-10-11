# Evermos backend mini project

REST API marketplace untuk mini project Rakamin x Evermos. Aplikasi menangani autentikasi, akun, toko, alamat pengiriman, kategori, produk, transaksi, dan data wilayah Indonesia.

## Fitur

- Register dan login menggunakan bcrypt serta JWT.
- Toko dibuat otomatis saat user mendaftar.
- Ownership authorization untuk akun, alamat, toko, produk, dan transaksi.
- Role-based access control untuk kategori admin.
- Upload gambar JPG/PNG maksimal 5 MB.
- Pagination dan filter produk atau toko.
- Checkout atomik dengan MySQL transaction dan row-level locking.
- Snapshot produk pada `log_produk` agar histori transaksi tidak berubah.
- Integrasi API Wilayah Indonesia.

## Teknologi

- Go 1.25+
- Fiber 2
- GORM
- MySQL 8
- JWT HS256 dan bcrypt

## Arsitektur

Project memakai modular monolith. Route meneruskan request ke middleware dan controller. Controller menjalankan aturan bisnis dan mengakses MySQL melalui GORM.

```text
Client → Fiber route → JWT/RBAC middleware → Controller → GORM → MySQL
```

Struktur folder:

```text
config/       koneksi dan seed database
controllers/  request, aturan bisnis, query, response
middleware/   JWT dan pemeriksaan admin
models/       model tabel dan relasi
routes/       registrasi endpoint
utils/        bcrypt, JWT, dan validasi
uploads/      penyimpanan gambar lokal
```

## Menjalankan aplikasi

Prasyarat: Go, MySQL 8, dan database `evermos_mini`.

```sql
CREATE DATABASE IF NOT EXISTS evermos_mini
CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

Salin `.env.example` menjadi `.env`, lalu isi koneksi lokal dan JWT secret. Jangan commit `.env`.

```bash
go mod download
go run .
```

Server berjalan di `http://localhost:8080`. GORM menjalankan migrasi dan aplikasi menambahkan kategori awal jika tabel kategori masih kosong.

## Konfigurasi

| Variable | Kegunaan | Contoh |
|---|---|---|
| `DB_HOST` | Host MySQL | `127.0.0.1` |
| `DB_PORT` | Port MySQL | `3306` |
| `DB_USER` | User MySQL | `root` |
| `DB_PASS` | Password MySQL | rahasia |
| `DB_NAME` | Nama database | `evermos_mini` |
| `APP_PORT` | Port HTTP | `:8080` |
| `JWT_SECRET` | Secret penandatangan token | string acak panjang |

## Endpoint utama

| Area | Method dan path | Akses |
|---|---|---|
| Auth | `POST /auth/register`, `POST /auth/login` | Publik |
| Profil | `GET /user`, `PUT /user` | User |
| Alamat | `GET/POST /user/alamat`, `GET/PUT/DELETE /user/alamat/:id` | Pemilik |
| Toko | `GET /toko`, `GET /toko/:id_toko` | Publik |
| Toko | `GET /toko/my`, `PUT /toko/:id_toko` | Pemilik |
| Kategori | `GET /category`, `GET /category/:id` | Publik |
| Kategori | `POST /category`, `PUT/DELETE /category/:id` | Admin |
| Produk | `GET /product`, `GET /product/:id` | Publik |
| Produk | `POST /product`, `PUT/DELETE /product/:id` | Pemilik toko |
| Transaksi | `GET /trx`, `GET /trx/:id`, `POST /trx` | Pemilik |
| Wilayah | `/provcity/*` | Publik |

Filter produk: `page`, `limit`, `nama_produk`, `category_id`, `toko_id`, `min_harga`, dan `max_harga`.

### Contoh register

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "nama":"Juan",
    "email":"juan@example.com",
    "no_telp":"081234567890",
    "kata_sandi":"password123",
    "tanggal_Lahir":"02/01/2000",
    "pekerjaan":"Developer",
    "id_provinsi":"11",
    "id_kota":"1101"
  }'
```

Endpoint protected menerima header:

```http
Authorization: Bearer <token>
```

## Postman

Import `Rakamin Evermos Virtual Internship.postman_collection.json`. Set variable `local` menjadi `http://localhost:8080`. Register dan login, lalu isi variable `token` dari response login.

Urutan smoke test:

1. Register dan login.
2. Buka profil dan toko otomatis.
3. Buat alamat.
4. Buat produk dengan multipart form-data.
5. Buat transaksi.
6. Periksa stok dan snapshot transaksi.

Untuk menguji kategori admin, ubah role langsung di MySQL lalu login ulang:

```sql
UPDATE users SET role = 'admin' WHERE id = 1;
```

## Pengujian

Unit test tidak membutuhkan database:

```bash
go test ./...
go vet ./...
```

Integration test memakai database terpisah agar data development tidak terhapus:

```sql
CREATE DATABASE IF NOT EXISTS evermos_mini_test
CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

```bash
TEST_DATABASE_DSN='root:password@tcp(127.0.0.1:3306)/evermos_mini_test?charset=utf8mb4&parseTime=True&loc=Local' \
go test ./... -count=1
```

Test mencakup bcrypt, JWT valid/rusak/kedaluwarsa, isolasi alamat antar-user, checkout, snapshot produk, pengurangan stok, rollback, dan penolakan alamat milik user lain.

## Catatan desain

Project ini bukan microservices. Satu aplikasi dan satu database cukup untuk scope mini project serta membuat transaction stok lebih sederhana. Pada deployment multi-instance, file upload perlu dipindahkan ke object storage. Sistem finansial produksi juga sebaiknya menyimpan uang sebagai integer satuan terkecil atau tipe decimal, bukan `float64`.
