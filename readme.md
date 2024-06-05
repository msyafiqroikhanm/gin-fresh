# API Aplikasi Peminjaman Kendaraan Operasional

Selamat datang di aplikasi peminjaman kendaraan operasional! Berikut adalah panduan instalasi dan penggunaan aplikasi ini yang menggunakan Gin dan GORM di Go.

## Prasyarat

Pastikan Anda sudah menginstal:

- [Go](https://golang.org/dl/) (Versi 1.16 ke atas)
- [Git](https://git-scm.com/)
- [PostgreSQL](https://www.postgresql.org/) (Aplikasi ini dikembangankan dengan menggunakan PostgreSQL, harap disesuaikan dengan relational database yang digunakan)

## Instalasi

1. Clone repository ini ke dalam direktori lokal Anda.
   ```bash
   git clone https://github.com/msyafiqroikhanm/go-vehicle-loans.git
   ```
2. Masuk ke direktori aplikasi.
   ```bash
   cd go-vehicle-loans
   ```
3. Install dependensi dengan menggunakan `go mod`.
   ```bash
   go mod tidy
   ```

## Konfigurasi

1. Salin file `.env.example` menjadi `.env`.
   ```bash
   cp .env.example .env
   ```
2. Edit file `.env` dan sesuaikan dengan konfigurasi database Anda. Contoh:

   ```env
    APP_NAME="Operational Vehicle Loan"
    PORT=5000
    JWT_SECRET=your_jwt_secret

    PGHOST=localhost
    PGPORT=5432
    PGUSER=postgres
    PGPASSWORD=root
    PGDATABASE=go-vehicle-loan
   ```

## Menjalankan Aplikasi

Untuk menjalankan aplikasi, gunakan perintah berikut:

```bash
go run main.go
```

Aplikasi akan berjalan di http://localhost:8080.

## Dokumentasi API

Untuk informasi lebih lanjut mengenai API yang digunakan oleh aplikasi ini, silakan kunjungi dokumentasi API di [Postman](https://documenter.getpostman.com/view/25285573/2sA3Qv8WFk).

### Daftar API

#### API User

- **Registrasi Pengguna:** `POST /register`
- **Login Pengguna:** `POST /login`
- **Profil Pengguna:** `GET /user/:id`
- **Update Profil Pengguna:** `PUT /user/:id`
- **Hapus Pengguna:** `DELETE /user/:id`

#### API Role

- **Daftar Role:** `GET /roles`
- **Tambah Role:** `POST /roles`
- **Detail Role:** `GET /roles/:id`
- **Update Role:** `PUT /roles/:id`
- **Hapus Role:** `DELETE /roles/:id`

#### API Loan

- **Daftar Loan:** `GET /loans`
- **Detail Loan:** `GET /loans/:id`
- **Tambah Loan:** `POST /loans`
- **Update Loan:** `PUT /loans/:id`
- **Hapus Loan:** `DELETE /loans/:id`
- **Return Loan:** `POST /loans/:id/return`

#### API Vehicle

- **Daftar Vehicle:** `GET /vehicles`
- **Detail Vehicle:** `GET /vehicles/:id`
- **Tambah Vehicle:** `POST /vehicles`
- **Update Vehicle:** `PUT /vehicles/:id`
- **Hapus Vehicle:** `DELETE /vehicles/:id`

#### API Vehicle Types

- **Daftar Vehicle Types:** `GET /vehicles/types`
- **Detail Vehicle Type:** `GET /vehicles/types/:id`
- **Tambah Vehicle Type:** `POST /vehicles/types`
- **Update Vehicle Type:** `PUT /vehicles/types/:id`
- **Hapus Vehicle Type:** `DELETE /vehicles/types/:id`

## Contoh Pengguna

Berikut adalah sample username dan password yang dapat Anda gunakan untuk login:

| Role  | Email             | Password |
| ----- | ----------------- | -------- |
| Admin | admin@example.com | password |
| User  | user@example.com  | password |

## Lisensi

Aplikasi ini dilisensikan di bawah MIT License.

Kritik saran kami perlukan untuk memperbaiki dan meningkatkan pengembangan aplikasi di masa mendatang.
Terima kasih telah menggunakan aplikasi kami!
