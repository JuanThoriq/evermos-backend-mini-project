# Evermos backend mini project

REST API e-commerce untuk tugas akhir Rakamin x Evermos. Aplikasi mencakup autentikasi, akun pengguna, toko, alamat pengiriman, kategori, produk, transaksi, dan integrasi data wilayah Indonesia.

## Teknologi

- Go dan Fiber
- GORM dan MySQL
- JWT HS256
- bcrypt

## Fitur sesuai requirement

- Register dan login. Registrasi membuat akun dan toko dalam satu database transaction.
- JWT authentication dan pembatasan data berdasarkan pemilik.
- CRUD alamat, profil, toko, kategori, dan produk.
- Kategori hanya dapat diubah admin.
- Upload gambar produk dan toko.
- Pagination dan filter produk atau toko.
- Transaksi atomik: validasi alamat, kunci stok, snapshot produk ke `log_produk`, simpan detail, lalu kurangi stok.
- Proxy data provinsi dan kota dari API Wilayah Indonesia.

## Menjalankan aplikasi

Prasyarat: Go 1.25+, MySQL 8+, dan database `evermos_mini`.

```sql
CREATE DATABASE IF NOT EXISTS evermos_mini
CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

Salin `.env.example` menjadi `.env`, lalu isi koneksi lokal dan JWT secret. File `.env` diabaikan Git.

```bash
go mod download
go run .
```

Server berjalan di `http://localhost:8080`. GORM membuat tabel dan aplikasi menambahkan tiga kategori awal.

## Pengujian

```bash
go test ./...
go vet ./...
```

Integration test memakai database terpisah agar data pengembangan tidak terhapus. Buat `evermos_mini_test`, lalu set DSN hanya pada terminal tempat test dijalankan:

```bash
TEST_DATABASE_DSN='root:password@tcp(127.0.0.1:3306)/evermos_mini_test?charset=utf8mb4&parseTime=True&loc=Local' go test ./... -count=1
```

Test mencakup hash password, JWT valid/rusak/kedaluwarsa, isolasi alamat antar-user, checkout berhasil, snapshot produk, pengurangan stok, rollback ketika salah satu item gagal, dan penolakan alamat milik user lain.

Import `Rakamin Evermos Virtual Internship.postman_collection.json` ke Postman. Set variable `local` menjadi `http://localhost:8080`. Setelah login, set `token` dari response login.

Urutan smoke test yang disarankan:

1. `POST /auth/register`
2. `POST /auth/login`
3. `GET /user` dan `GET /toko/my`
4. `POST /user/alamat`
5. `POST /product` dengan multipart form-data
6. `POST /trx`
7. Periksa stok melalui `GET /product/:id`

Untuk menguji admin, ubah role langsung di MySQL lalu login ulang:

```sql
UPDATE users SET role = 'admin' WHERE id = 1;
```

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

## Catatan desain

Project ini merupakan modular monolith, bukan microservices. Pemisahan per domain dilakukan pada controller, route, dan model. Untuk skala mini project, satu deployment dan satu database lebih mudah dijalankan serta menjaga transaksi stok tetap konsisten. Pemecahan menjadi microservices baru masuk akal saat domain perlu deployment atau scaling mandiri.
