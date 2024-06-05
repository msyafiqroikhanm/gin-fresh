# API Aplikasi E-Procurement

Selamat datang di aplikasi peminjaman kendaraan operasional! Berikut adalah panduan instalasi dan penggunaan aplikasi ini yang menggunakan Gin dan GORM di Go.

## Prasyarat

Pastikan Anda sudah menginstal:

- [Go](https://golang.org/dl/) (Versi 1.16 ke atas)
- [Git](https://git-scm.com/)
- [PostgreSQL](https://www.postgresql.org/) (Aplikasi ini dikembangankan dengan menggunakan PostgreSQL, harap disesuaikan dengan relational database yang digunakan)

## Instalasi

1. Clone repository ini ke dalam direktori lokal Anda.
   ```bash
   git clone https://github.com/msyafiqroikhanm/jxb-eprocurement.git
   ```
2. Masuk ke direktori aplikasi.
   ```bash
   cd jxb-eprocurement
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

## Lisensi

Aplikasi ini dilisensikan di bawah MIT License.

Kritik saran kami perlukan untuk memperbaiki dan meningkatkan pengembangan aplikasi di masa mendatang.
Terima kasih telah menggunakan aplikasi kami!
