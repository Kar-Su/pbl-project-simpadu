# pbl-project-simpadu

Repository proyek PBL (Project Based Learning) - Sistem Informasi Manajemen Padu (SIMPADU).
Berikut adalah panduan untuk menjalankan proyek ini di lingkungan lokal menggunakan Docker.

## Requirement

- [Git](https://git-scm.com/)
- [Docker](https://www.docker.com/)
-  [Docker Compose](https://docs.docker.com/compose/)

## Tutorial Penggunaan Repo

### 1. Clone Repository

    git clone https://github.com/Kar-Su/pbl-project-simpadu.git
    cd pbl-project-simpadu

### 2. Siapkan File Environment

Salin file `.env.example` yang ada di root proyek dan di dalam folder `backend/` menjadi `.env`.

**Linux / macOS:**

    cp .env.example .env
    cp ./backend/.env.example ./backend/.env

**Windows (Command Prompt):**

    copy .env.example .env
    copy .\backend\.env.example .\backend\.env

**Windows (PowerShell):**

    Copy-Item .env.example .env
    Copy-Item .\backend\.env.example .\backend\.env

> **Catatan:** Jika file `.env` sudah ada, Anda dapat menimpanya atau menyesuaikan isinya sesuai kebutuhan.

### 3. Jalankan Docker Containers

    docker compose up -d

Perintah ini akan membangun (jika diperlukan) dan menjalankan semua layanan dalam mode background. Tunggu hingga semua container siap.

### 4. Seed Database

Setelah container berjalan, isi basis data dengan data awal (seeder):

    docker exec -i golang_pbl go run cmd/seeder/main.go --seed

## Data Awal (Seed)

Setelah proses seeding berhasil, berikut data yang tersedia:

- **Daftar Pengguna (Users):** [users.json](https://github.com/Kar-Su/pbl-project-simpadu/blob/main/backend/internal/database/seeders/json/users.json)
- **Daftar Peran (Roles):** [roles.json](https://github.com/Kar-Su/pbl-project-simpadu/blob/main/backend/internal/database/seeders/json/roles.json)

Gunakan data tersebut untuk login atau pengujian API.

## Dokumentasi API

Swagger UI tersedia di:  
[http://localhost/api/swagger/index.html](http://localhost/api/swagger/index.html)

## Video Tutorial

Cara melakukan request API melalui Swagger UI:

https://github.com/user-attachments/assets/bd0f7326-edad-43c3-92aa-148018816216


## Contoh Payload JWT

<img width="1035" height="537" alt="Image" src="https://github.com/user-attachments/assets/fc6cc8c8-1d9c-4b52-b5f3-6e6ab0c653b1" />

---

Jika mengalami masalah, Hubungi Tim 1.


