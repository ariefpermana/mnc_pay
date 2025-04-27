# mnc_pay

# MNC Payment API

## 🚀 Setup Project

Berikut adalah langkah-langkah untuk menjalankan proyek ini di environment local.

---

## 📦 1. Generate `db.sql`

Untuk membuat struktur database secara otomatis, ikuti langkah-langkah berikut:
1. running `db.sql`

## ⚙️ 2. Setup `.env`

Update file `.env` di root project dan isi dengan konfigurasi database local seperti contoh berikut:

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=mnc_payment

## 3. running services
Jalankan command di terminal:

    ```bash
    go mod tidy
    go run main.go
    ```

   Ini akan menginstall semua dependensi yang diperlukan dan menjalankan program untuk melakukan migrasi ke database sesuai dengan model yang ada.


---
## 4. Testing API
- Import collection `mnc_pay.postman_collection.json` di Postman untuk menguji endpoint API dan melakukan test hit pada service.
