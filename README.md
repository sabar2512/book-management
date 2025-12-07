# Book Management API

## 📋 Daftar Isi

- [Deskripsi Singkat](#-deskripsi-singkat)
- [Fitur Utama](#-fitur-utama)
- [Teknologi](#-teknologi)
- [Prerequisites](#-prerequisites)
- [Instalasi](#-instalasi)
- [Konfigurasi](#-konfigurasi)
- [Database](#-database)
- [API Endpoints](#-api-endpoints)
- [Contoh Penggunaan](#-contoh-penggunaan)
- [Deployment Railway](#-deployment-railway)
- [Struktur Project](#-struktur-project)
- [Testing](#-testing)
- [Contributing](#-contributing)
- [License](#-license)

---

## 🎯 Deskripsi Singkat

**Book Management API** adalah layanan backend berbasis Golang yang menyediakan fitur lengkap untuk manajemen buku dan kategori. Aplikasi ini dilengkapi dengan sistem autentikasi JWT, validasi input otomatis, dan perhitungan ketebalan buku berdasarkan jumlah halaman.

Dirancang untuk deployment di cloud platform seperti Railway, aplikasi ini menggunakan Docker multi-stage build untuk optimasi performa dan ukuran image.

**Seluruh endpoint utama berada di jalur `/api/*`**, sementara endpoint root (`/`) memberikan dokumentasi ringkas mengenai status dan daftar endpoint yang tersedia.

---

## ✨ Fitur Utama

### 🔐 Autentikasi (JWT)
- ✅ Login menggunakan username dan password
- ✅ JWT token dengan expiry 24 jam
- ✅ Middleware otomatis memvalidasi token
- ✅ Proteksi semua endpoint sensitif

### 📚 Manajemen Buku
- ✅ CRUD lengkap untuk buku
- ✅ Auto-calculate ketebalan buku (tipis/tebal berdasarkan halaman)
- ✅ Validasi release year (1980-2024)
- ✅ Validasi kategori harus exist
- ✅ Support image URL
- ✅ Audit trail (created_by, modified_by)

### 🏷️ Manajemen Kategori
- ✅ CRUD lengkap untuk kategori
- ✅ Relasi one-to-many dengan buku
- ✅ Endpoint khusus untuk list buku per kategori
- ✅ Cascade delete (hapus kategori = hapus buku terkait)

---

## 🛠 Teknologi

### Backend
- **Go** 1.21+ - Programming language
- **Gin** - Web framework
- **PostgreSQL** - Relational database
- **JWT** - Authentication & authorization
- **sql-migrate** - Database migration tool

### Deployment
- **Docker** - Containerization
- **Railway** - Cloud platform deployment
- **Multi-stage build** - Optimized Docker image

### Libraries
```go
github.com/gin-gonic/gin           // Web framework
github.com/lib/pq                  // PostgreSQL driver
github.com/golang-jwt/jwt/v5       // JWT implementation
github.com/joho/godotenv           // Environment variables
```

---

## 📦 Prerequisites

Sebelum memulai, pastikan sudah terinstall:

- **Go** 1.21 atau lebih tinggi ([Download](https://go.dev/dl/))
- **PostgreSQL** 12 atau lebih tinggi ([Download](https://www.postgresql.org/download/))
- **Git** ([Download](https://git-scm.com/downloads))
- **sql-migrate** (akan diinstall saat setup)

---

## 🚀 Instalasi

### 1. Clone Repository
```bash
git clone <repository-url>
cd book-management
```

### 2. Install Dependencies
```bash
# Download Go modules
go mod download

# Install sql-migrate tool
go install github.com/rubenv/sql-migrate/...@latest
```

### 3. Setup Database
```bash
# Start PostgreSQL service
sudo systemctl start postgresql  # Linux
# atau
brew services start postgresql   # macOS

# Create database
sudo -u postgres psql
CREATE DATABASE bookdb;
\q
```

### 4. Konfigurasi Environment
```bash
# Copy template environment
cp .env.example .env

# Edit file .env dengan kredensial database Anda
nano .env
```

### 5. Jalankan Migration
```bash
# Run database migrations
sql-migrate up

# Expected output:
# Applied 3 migrations
# +-- 001_create_users_table.sql
# +-- 002_create_categories_table.sql
# +-- 003_create_books_table.sql
```

### 6. Seed Data (Opsional)
```bash
# Populate database dengan sample data
go run seed/main.go
```

### 7. Jalankan Aplikasi
```bash
# Run server
go run main.go

# Output:
# Database connected successfully
# Server running on port 8080
```

**✅ Aplikasi berjalan di `http://localhost:8080`**

---

## ⚙️ Konfigurasi

### Environment Variables (.env)

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=bookdb
DB_SSLMODE=disable

# Untuk Railway, gunakan DATABASE_URL
# DATABASE_URL=postgresql://user:pass@host:port/dbname

# Server Configuration
PORT=8080
GIN_MODE=debug              # Gunakan 'release' untuk production

# JWT Secret Key (WAJIB diganti untuk production!)
JWT_SECRET=your-super-secret-key-change-this
```

**⚠️ PENTING:**
- Ganti `DB_PASSWORD` dengan password PostgreSQL Anda
- Ganti `JWT_SECRET` dengan string random yang aman untuk production
- Untuk Railway, `DATABASE_URL` akan di-set otomatis

---

## 🗄️ Database

### Struktur Tabel

#### 1. Tabel Users
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(100),
    modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_by VARCHAR(100)
);
```

#### 2. Tabel Categories
```sql
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(100),
    modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_by VARCHAR(100)
);
```

#### 3. Tabel Books
```sql
CREATE TABLE books (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    image_url TEXT,
    release_year INTEGER NOT NULL,
    price NUMERIC NOT NULL,
    total_page INTEGER NOT NULL,
    thickness VARCHAR(50) NOT NULL,  -- auto-calculated: 'tipis' atau 'tebal'
    category_id INTEGER NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(100),
    modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_by VARCHAR(100)
);

CREATE INDEX idx_books_category_id ON books(category_id);
```

### Aturan Business Logic

**Thickness Calculation:**
- `total_page ≤ 100` → `"tipis"`
- `total_page > 100` → `"tebal"`

**Validasi:**
- `release_year`: 1980 - 2024
- `price`: ≥ 0
- `total_page`: ≥ 1
- `category_id`: harus exist di tabel categories

---

## 🔌 API Endpoints

### Root Endpoint
```
GET /
```
Menampilkan dokumentasi ringkas API dan daftar endpoint.

**Response:**
```json
{
  "message": "Book Management API is running 🚀",
  "version": "1.0.0",
  "endpoints": {
    "Authentication": {
      "POST /api/login": "Get JWT token"
    },
    "Books": {
      "GET /api/books": "List all books",
      "POST /api/books": "Create new book",
      "GET /api/books/:id": "Get book detail",
      "PUT /api/books/:id": "Update book",
      "DELETE /api/books/:id": "Delete book"
    },
    "Categories": {
      "GET /api/categories": "List all categories",
      "POST /api/categories": "Create new category",
      "GET /api/categories/:id": "Get category detail",
      "PUT /api/categories/:id": "Update category",
      "DELETE /api/categories/:id": "Delete category",
      "GET /api/categories/:id/books": "Get books by category"
    }
  }
}
```

### Authentication

#### Login
```http
POST /api/login
Content-Type: application/json
```

**Request Body:**
```json
{
  "username": "admin@example.com",
  "password": "12345"
}
```

**Response:**
```json
{
  "message": "Login successful",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "username": "admin@example.com"
}
```

**💡 Gunakan token untuk semua request protected:**
```
Authorization: Bearer <JWT_TOKEN>
```

---

### Books Endpoints

Semua endpoint books memerlukan JWT token dalam header `Authorization`.

#### 1. Get All Books
```http
GET /api/books
Authorization: Bearer <token>
```

**Response:**
```json
{
  "data": [
    {
      "id": 1,
      "title": "Bumi Manusia",
      "description": "Novel sejarah Indonesia",
      "image_url": "https://example.com/book.jpg",
      "release_year": 1980,
      "price": 120000,
      "total_page": 500,
      "thickness": "tebal",
      "category_id": 1,
      "created_at": "2024-01-01T10:00:00Z",
      "created_by": "admin",
      "modified_at": "2024-01-01T10:00:00Z",
      "modified_by": "admin"
    }
  ]
}
```

#### 2. Create Book
```http
POST /api/books
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "title": "Laskar Pelangi",
  "description": "Novel karya Andrea Hirata",
  "image_url": "https://example.com/laskar-pelangi.jpg",
  "release_year": 2005,
  "price": 85000,
  "total_page": 529,
  "category_id": 1
}
```

**Response:**
```json
{
  "message": "Book created successfully",
  "id": 2,
  "thickness": "tebal"
}
```

#### 3. Get Book by ID
```http
GET /api/books/:id
Authorization: Bearer <token>
```

#### 4. Update Book
```http
PUT /api/books/:id
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:** (sama seperti Create Book)

#### 5. Delete Book
```http
DELETE /api/books/:id
Authorization: Bearer <token>
```

**Response:**
```json
{
  "message": "Book deleted successfully"
}
```

---

### Categories Endpoints

Semua endpoint categories memerlukan JWT token.

#### 1. Get All Categories
```http
GET /api/categories
Authorization: Bearer <token>
```

#### 2. Create Category
```http
POST /api/categories
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "Novel"
}
```

#### 3. Get Category by ID
```http
GET /api/categories/:id
Authorization: Bearer <token>
```

#### 4. Update Category
```http
PUT /api/categories/:id
Authorization: Bearer <token>
Content-Type: application/json
```

#### 5. Delete Category
```http
DELETE /api/categories/:id
Authorization: Bearer <token>
```

#### 6. Get Books by Category
```http
GET /api/categories/:id/books
Authorization: Bearer <token>
```

**Response:**
```json
{
  "data": [
    {
      "id": 1,
      "title": "Bumi Manusia",
      "category_id": 1,
      ...
    }
  ]
}
```

---

## 💡 Contoh Penggunaan

### Scenario: Menambah Buku Baru

```bash
# 1. Login untuk mendapatkan token
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin@example.com",
    "password": "12345"
  }'

# Response: {"token": "eyJhbG..."}

# 2. Simpan token
export TOKEN="eyJhbG..."

# 3. Buat kategori baru
curl -X POST http://localhost:8080/api/categories \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Fiksi Ilmiah"
  }'

# 4. Tambah buku dengan kategori tersebut
curl -X POST http://localhost:8080/api/books \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "Negeri 5 Menara",
    "description": "Novel inspiratif karya A. Fuadi",
    "image_url": "https://example.com/negeri5menara.jpg",
    "release_year": 2009,
    "price": 95000,
    "total_page": 432,
    "category_id": 1
  }'

# Response: {"id": 3, "thickness": "tebal"}

# 5. Lihat semua buku
curl http://localhost:8080/api/books \
  -H "Authorization: Bearer $TOKEN"
```

---

## 🚢 Deployment Railway

### Persiapan

Aplikasi ini sudah dilengkapi dengan `Dockerfile` yang dioptimasi untuk deployment di Railway:

```dockerfile
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
```

### Langkah Deployment

#### 1. Via Railway CLI

```bash
# Install Railway CLI
npm i -g @railway/cli

# Login
railway login

# Initialize project
railway init

# Add PostgreSQL database
# (Lakukan dari Railway dashboard)

# Set environment variables
railway variables set JWT_SECRET=your-super-secret-key
railway variables set GIN_MODE=release

# Deploy
railway up
```

#### 2. Via GitHub Integration

1. Push code ke GitHub repository
2. Buka [Railway Dashboard](https://railway.app)
3. Click "New Project" → "Deploy from GitHub repo"
4. Pilih repository Anda
5. Railway akan auto-detect Dockerfile
6. Tambahkan PostgreSQL database:
   - Click "New" → "Database" → "PostgreSQL"
7. Set environment variables:
   - `JWT_SECRET`: Secret key untuk JWT
   - `GIN_MODE`: `release`
   - `DATABASE_URL`: (auto-set oleh Railway)
8. Deploy!

#### 3. Run Migrations di Railway

```bash
# Connect to Railway project
railway link

# Run migrations
railway run sql-migrate up -env=production
```

### Environment Variables di Railway

Railway akan otomatis menyediakan:
- ✅ `DATABASE_URL` - PostgreSQL connection string
- ✅ `PORT` - Port untuk aplikasi

Yang perlu Anda set manual:
- ⚙️ `JWT_SECRET` - Secret key untuk JWT
- ⚙️ `GIN_MODE` - Set ke `release` untuk production

### Generate Public URL

1. Go to Railway dashboard
2. Click project → Settings
3. Click "Generate Domain"
4. Aplikasi dapat diakses di `https://your-app.up.railway.app`

---

## 📁 Struktur Project

```
book-management/
├── main.go                    # Entry point aplikasi
├── go.mod                     # Go module dependencies
├── go.sum                     # Dependencies checksum
├── Dockerfile                 # Docker configuration
├── docker-compose.yml         # Docker Compose setup
├── railway.json              # Railway deployment config
├── .env.example              # Environment template
├── .gitignore                # Git ignore rules
│
├── config/                   # Konfigurasi
│   └── database.go          # Database connection
│
├── middleware/               # HTTP middleware
│   └── auth.go              # JWT authentication
│
├── models/                   # Data models
│   ├── book.go              # Book model
│   ├── category.go          # Category model
│   └── user.go              # User model
│
├── handlers/                 # Request handlers
│   ├── auth.go              # Login handler
│   ├── book.go              # Book CRUD handlers
│   └── category.go          # Category CRUD handlers
│
├── routes/                   # Route definitions
│   └── routes.go            # API routes setup
│
├── migrations/               # Database migrations
│   ├── 001_create_users_table.sql
│   ├── 002_create_categories_table.sql
│   └── 003_create_books_table.sql
│
├── seed/                     # Database seeding
│   └── main.go              # Seed script
│
└── docs/                     # Documentation
    ├── README.md            # This file
    ├── CARA_MENJALANKAN.md  # Setup guide (ID)
    ├── TESTING.md           # Testing guide
    └── API.md               # API documentation
```

---

## 🧪 Testing

### Manual Testing dengan cURL

Lihat file [TESTING.md](TESTING.md) untuk panduan lengkap testing.

### Import Postman Collection

1. Import file `postman_collection.json`
2. Set environment variables:
   - `base_url`: `http://localhost:8080`
   - `token`: (akan di-set otomatis setelah login)
3. Run collection

### Automated Testing

```bash
# Run tests
go test -v ./...

# Run with coverage
go test -v -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## 🤝 Contributing

Kontribusi sangat diterima! Berikut cara berkontribusi:

1. Fork repository ini
2. Buat branch baru (`git checkout -b feature/AmazingFeature`)
3. Commit perubahan (`git commit -m 'Add: amazing feature'`)
4. Push ke branch (`git push origin feature/AmazingFeature`)
5. Buat Pull Request

Lihat [CONTRIBUTING.md](CONTRIBUTING.md) untuk panduan lengkap.

### Commit Message Format

- `Add:` untuk fitur baru
- `Fix:` untuk bug fixes
- `Update:` untuk perubahan fitur existing
- `Docs:` untuk dokumentasi
- `Refactor:` untuk refactoring

---

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.

```
MIT License

Copyright (c) 2024 Book Management API

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction...
```

---

## 📞 Support

Jika menemukan bug atau ada pertanyaan:

- 🐛 [Report Bug](https://github.com/username/repo/issues)
- 💡 [Request Feature](https://github.com/username/repo/issues)
- 📧 Email: support@example.com

---

## 🙏 Acknowledgments

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [PostgreSQL](https://www.postgresql.org/)
- [JWT-Go](https://github.com/golang-jwt/jwt)
- [Railway](https://railway.app)

---
