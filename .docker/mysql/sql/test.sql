DROP DATABASE IF EXISTS test;
CREATE DATABASE test
    DEFAULT CHARACTER SET utf8mb4
    DEFAULT COLLATE utf8mb4_bin;

use test

CREATE TABLE if not exists roles (
    id int PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,

    UNIQUE INDEX idx_roles_name (name)
) engine=InnoDB;

CREATE TABLE if not exists users (
    id char(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    image_url VARCHAR(255),
    role_id int NOT NULL,
    detail_id char(36),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE INDEX idx_users_email (email),
    UNIQUE INDEX idx_users_detail (detail_id),
    CONSTRAINT fk_role FOREIGN KEY (role_id) REFERENCES roles(id)
) engine=InnoDB;

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id char(36) PRIMARY KEY,
    user_id char(36) NOT NULL,
    token VARCHAR(255) NOT NULL,
    expired_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE INDEX idx_token (token),
    INDEX idx_token_user (user_id),
    CONSTRAINT fk_token_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) engine=InnoDB;

CREATE TABLE IF NOT EXISTS jurusan(
    id int PRIMARY KEY AUTO_INCREMENT,
    nama VARCHAR(255) NOT NULL,

    UNIQUE INDEX idx_jurusan_nama (nama)
) engine=InnoDB;

CREATE TABLE IF NOT EXISTS prodi(
    id int PRIMARY KEY AUTO_INCREMENT,
    nama VARCHAR(255) NOT NULL,
    jurusan_id int NOT NULL,

    UNIQUE INDEX idx_prodi_nama (nama),
    INDEX idx_prodi_jurusan (jurusan_id),
    CONSTRAINT fk_prodi_jurusan FOREIGN KEY (jurusan_id) REFERENCES jurusan(id) ON DELETE CASCADE
) engine=InnoDB;

CREATE TABLE IF NOT EXISTS mata_kuliah(
    id char(36) PRIMARY KEY,
    nama_mk VARCHAR(255) NOT NULL,
    sks int NOT NULL
) engine=InnoDB;

CREATE TABLE IF NOT EXISTS kurikulum(
    id char(36) PRIMARY KEY,
    nama_kurikulum VARCHAR(255) NOT NULL,
    prodi_id int NOT NULL,

    CONSTRAINT fk_kurikulum_prodi FOREIGN KEY (prodi_id) REFERENCES prodi(id) ON DELETE CASCADE
) engine=InnoDB;

CREATE TABLE IF NOT EXISTS kurikulum_mk(
    kurikulum_id char(36) NOT NULL,
    mata_kuliah_id char(36) NOT NULL,

    PRIMARY KEY (kurikulum_id, mata_kuliah_id),
    CONSTRAINT fk_kurikulum_mk_kurikulum FOREIGN KEY (kurikulum_id) REFERENCES kurikulum(id) ON DELETE CASCADE,
    CONSTRAINT fk_kurikulum_mk_mata_kuliah FOREIGN KEY (mata_kuliah_id) REFERENCES mata_kuliah(id) ON DELETE CASCADE
) engine=InnoDB;

CREATE TABLE IF NOT EXISTS tahun_akademik(
    id int PRIMARY KEY,
    tipe_semester ENUM('ganjil', 'genap') NOT NULL,
    tahun_awal date NOT NULL,
    tahun_akhir date NOT NULL,

    UNIQUE INDEX idx_tahun_akademik_tahun (tahun_awal, tahun_akhir)
) engine=InnoDB;

CREATE TABLE IF NOT EXISTS kelas(
    id char(36) PRIMARY KEY,
    nama_kelas VARCHAR(255) NOT NULL,
    tahun_akademik_id int NOT NULL,
    kurikulum_id char(36) NOT NULL,
    prodi_id int NOT NULL,

    CONSTRAINT fk_kelas_tahun_akademik FOREIGN KEY (tahun_akademik_id) REFERENCES tahun_akademik(id) ON DELETE CASCADE,
    CONSTRAINT fk_kelas_kurikulum FOREIGN KEY (kurikulum_id) REFERENCES kurikulum(id) ON DELETE CASCADE,
    CONSTRAINT fk_kelas_prodi FOREIGN KEY (prodi_id) REFERENCES prodi(id) ON DELETE CASCADE
) engine=InnoDB;

CREATE TABLE IF NOT EXISTS kelas_mahasiswa(
    kelas_id char(36) NOT NULL,
    mahasiswa_id char(36) NOT NULL,

    PRIMARY KEY (kelas_id, mahasiswa_id),
    CONSTRAINT fk_kelas_mahasiswa_kelas FOREIGN KEY (kelas_id) REFERENCES kelas(id) ON DELETE CASCADE,
    CONSTRAINT fk_kelas_mahasiswa_mahasiswa FOREIGN KEY (mahasiswa_id) REFERENCES users(detail_id) ON UPDATE CASCADE
) engine=InnoDB;
