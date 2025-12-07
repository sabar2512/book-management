# Book Management API

Book Management API adalah layanan backend berbasis Golang yang menyediakan fitur autentikasi JWT, manajemen buku, dan manajemen kategori. Aplikasi ini menggunakan PostgreSQL sebagai database utama dan dideploy melalui Railway menggunakan Dockerfile multi-stage.

---

## 1. Deskripsi Singkat

API ini dirancang untuk kebutuhan pengelolaan buku dan kategori dengan fitur login dan proteksi akses melalui JWT. Seluruh endpoint utama berada di jalur `/api/*`, sementara endpoint root (`/`) memberikan dokumentasi ringkas mengenai layanan.

Aplikasi ini dioptimalkan untuk lingkungan cloud, menggunakan Docker sebagai runtime dan Railway sebagai platform deployment.

---

## 2. Fitur Utama

### 2.1 Autentikasi (JWT)
- Login menggunakan username dan password.
- Mendapatkan JWT token untuk akses endpoint yang dilindungi.
- Middleware memvalidasi token sebelum request diproses.

### 2.2 Manajemen Buku
Endpoint:
- GET `/api/books`
- GET `/api/books/:id`
- POST `/api/books`
- PUT `/api/books/:id`
- DELETE `/api/books/:id`

Model buku mendukung informasi detail seperti deskripsi, gambar, harga, total halaman, tingkat ketebalan, dan metadata pembuatan/perubahan.

### 2.3 Manajemen Kategori
Endpoint:
- GET `/api/categories`
- GET `/api/categories/:id`
- POST `/api/categories`
- PUT `/api/categories/:id`
- DELETE `/api/categories/:id`
- GET `/api/categories/:id/books`

Kategori dapat digunakan sebagai pengelompokan buku dan menjadi foreign key dalam tabel buku.

---

## 3. Struktur Folder Project

book-management/
│── main.go
│── go.mod
│── Dockerfile
│── config/ : koneksi database
│── routes/ : konfigurasi routing
│── handlers/ : handler untuk auth, books, categories
│── middleware/ : middleware JWT
│── models/ : struktur data (Book dan Category)

yaml
Salin kode

Struktur ini memisahkan logika aplikasi, membuat kode bersih, mudah dipahami, dan mudah dikembangkan lebih lanjut.

---

## 4. Struktur Database

### 4.1 Tabel Users
Kolom:
- id (SERIAL, PK)
- username (TEXT, UNIQUE)
- password (TEXT)
- created_at
- created_by
- modified_at
- modified_by

### 4.2 Tabel Categories
Kolom:
- id (SERIAL, PK)
- name (TEXT)
- created_at
- created_by
- modified_at
- modified_by

### 4.3 Tabel Books
Kolom:
- id (SERIAL, PK)
- title (TEXT)
- description (TEXT)
- image_url (TEXT)
- release_year (INTEGER)
- price (NUMERIC)
- total_page (INTEGER)
- thickness (INTEGER)
- category_id (INTEGER, FK → categories.id)
- created_at
- created_by
- modified_at
- modified_by

Struktur tabel ini disesuaikan penuh dengan handler Golang sehingga tidak terjadi error saat pemetaan data.

---

## 5. Contoh Penggunaan API

### 5.1 Login (mendapatkan JWT)
POST /api/login

css
Salin kode

Body:
```json
{
  "username": "admin@example.com",
  "password": "12345"
}
Response:

json
Salin kode
{
  "message": "Login successful",
  "token": "<JWT_TOKEN>",
  "username": "admin@example.com"
}
Gunakan token untuk semua request protected:

makefile
Salin kode
Authorization: Bearer <JWT_TOKEN>
5.2 Menambah Buku
bash
Salin kode
POST /api/books
Header:

pgsql
Salin kode
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json
Body:

json
Salin kode
{
  "title": "Bumi Manusia",
  "description": "Novel sejarah Indonesia",
  "image_url": "https://example.com/book.jpg",
  "release_year": 1980,
  "price": 120000,
  "total_page": 500,
  "category_id": 1
}
5.3 Mendapatkan Semua Buku
sql
Salin kode
GET /api/books
Authorization: Bearer <JWT_TOKEN>
6. Deployment Railway
Aplikasi dideploy menggunakan Dockerfile berikut:

dockerfile
Salin kode
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

FROM alpine:3.19
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
Environment Variable Wajib:
DATABASE_URL → URL PostgreSQL dari Railway

JWT_SECRET → secret key JWT

7. Dokumentasi Root
Akses:

sql
Salin kode
GET /
Output memberikan ringkasan status aplikasi dan daftar endpoint yang tersedia.

Contoh:

json
Salin kode
{
  "message": "Book Management API is running 🚀",
  "endpoints": {
    "Books": { "...": "..." },
    "Categories": { "...": "..." }
  }
}
8. Lisensi
MIT — dapat digunakan dan dimodifikasi bebas.

9. Kontribusi
Pull request dan masukan sangat diterima.