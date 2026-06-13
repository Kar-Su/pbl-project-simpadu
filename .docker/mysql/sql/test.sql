-- DROP DATABASE IF EXISTS test;
-- CREATE DATABASE test
--     DEFAULT CHARACTER SET utf8mb4
--     DEFAULT COLLATE utf8mb4_bin;

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
    detail_id char(36) DEFAULT NULL,
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
    name VARCHAR(255) NOT NULL,

    UNIQUE INDEX idx_jurusan_name (name)
) engine=InnoDB;

CREATE TABLE IF NOT EXISTS prodi(
    id int PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    jenjang enum('D3', 'D4') NOT NULL,
    jurusan_id int NOT NULL,

    UNIQUE INDEX idx_prodi_name (name),
    INDEX idx_prodi_jurusan (jurusan_id),
    CONSTRAINT fk_prodi_jurusan FOREIGN KEY (jurusan_id) REFERENCES jurusan(id) ON DELETE CASCADE
) engine=InnoDB;

CREATE TABLE IF NOT EXISTS mata_kuliah(
    id char(36) PRIMARY KEY,
    kode char(12) NOT NULL,
    name VARCHAR(255) NOT NULL,
    sks int NOT NULL,

    UNIQUE INDEX idx_mata_kuliah_kode (kode)
) engine=InnoDB;

CREATE TABLE IF NOT EXISTS kurikulum(
    id char(36) PRIMARY KEY,
    kode char(12) NOT NULL,
    name VARCHAR(255) NOT NULL,
    prodi_id int NOT NULL,

    UNIQUE INDEX idx_kurikulum_kode (kode),
    CONSTRAINT fk_kurikulum_prodi FOREIGN KEY (prodi_id) REFERENCES prodi(id) ON DELETE CASCADE
) engine=InnoDB;

CREATE TABLE IF NOT EXISTS kurikulum_mk(
    kurikulum_kode char(12) NOT NULL,
    mk_kode char(12) NOT NULL,
    semester int NOT NULL,
    wajib boolean NOT NULL DEFAULT false,

    PRIMARY KEY (kurikulum_kode, mk_kode),
    CONSTRAINT fk_kurikulum_mk_kurikulum FOREIGN KEY (kurikulum_kode) REFERENCES kurikulum(kode) ON DELETE CASCADE,
    CONSTRAINT fk_kurikulum_mk_mata_kuliah FOREIGN KEY (mk_kode) REFERENCES mata_kuliah(kode) ON DELETE CASCADE
) engine=InnoDB;

CREATE TABLE IF NOT EXISTS tahun_akademik(
    id int PRIMARY KEY,
    tipe_semester ENUM('ganjil', 'genap') NOT NULL,
    tahun_awal date NOT NULL,
    tahun_akhir date NOT NULL,
    status ENUM('aktif', 'nonaktif') NOT NULL DEFAULT 'aktif',

    UNIQUE INDEX idx_tahun_akademik_tahun (tahun_awal, tahun_akhir)
) engine=InnoDB;

CREATE TABLE IF NOT EXISTS kelas(
    id char(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    semester int NOT NULL DEFAULT 1,
    tahun_akademik_id int NOT NULL,
    kurikulum_kode char(12) NOT NULL,
    prodi_id int NOT NULL,

    INDEX idx_name (name),
    CONSTRAINT fk_kelas_tahun_akademik FOREIGN KEY (tahun_akademik_id) REFERENCES tahun_akademik(id) ON DELETE CASCADE,
    CONSTRAINT fk_kelas_kurikulum FOREIGN KEY (kurikulum_kode) REFERENCES kurikulum(kode) ON DELETE CASCADE,
    CONSTRAINT fk_kelas_prodi FOREIGN KEY (prodi_id) REFERENCES prodi(id) ON DELETE CASCADE
) engine=InnoDB;

CREATE TABLE IF NOT EXISTS kelas_mahasiswa(
    kelas_id char(36) NOT NULL,
    mahasiswa_id char(36) NOT NULL,

    PRIMARY KEY (kelas_id, mahasiswa_id),
    CONSTRAINT fk_kelas_mahasiswa_kelas FOREIGN KEY (kelas_id) REFERENCES kelas(id) ON DELETE CASCADE,
    CONSTRAINT fk_kelas_mahasiswa_mahasiswa FOREIGN KEY (mahasiswa_id) REFERENCES users(detail_id) ON UPDATE CASCADE
) engine=InnoDB;

CREATE TABLE IF NOT EXISTS pengampu(
    id char(36) PRIMARY KEY,
    kelas_id char(36) NOT NULL,
    mk_kode char(12) NOT NULL,
    dosen_id char(36) NOT NULL,

    CONSTRAINT idx_kelas_mk UNIQUE (kelas_id, mk_kode),
    INDEX idx_dosen_id (dosen_id),

    CONSTRAINT fk_pengampu_kelas FOREIGN KEY (kelas_id) REFERENCES kelas(id) ON DELETE CASCADE,
    CONSTRAINT fk_pengampu_mk FOREIGN KEY (mk_kode) REFERENCES mata_kuliah(kode) ON DELETE CASCADE,
    CONSTRAINT fk_pengampu_dosen FOREIGN KEY (dosen_id) REFERENCES users(detail_id) ON UPDATE CASCADE
) engine=InnoDB;

CREATE TABLE IF NOT EXISTS presensi (
    id char(36) PRIMARY KEY,
    tipe enum('mahasiswa', 'pegawai') NOT NULL,
    pengampu_id char(36) DEFAULT NULL,
    created_at date DEFAULT (CURRENT_DATE),
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_created_at (created_at),
    CONSTRAINT fk_presensi_pengampu FOREIGN KEY (pengampu_id) REFERENCES pengampu(id) ON UPDATE CASCADE
) engine=InnoDB;

CREATE TABLE IF NOT EXISTS presensi_mahasiswa (
    presensi_id char(36) NOT NULL,
    mahasiswa_id char(36) NOT NULL,
    status enum('hadir', 'sakit', 'izin', 'alpha') DEFAULT 'alpha',

    PRIMARY KEY (presensi_id, mahasiswa_id),
    CONSTRAINT fk_presensi_mahasiswa_presensi FOREIGN KEY (presensi_id) REFERENCES presensi(id) ON DELETE CASCADE,
    CONSTRAINT fk_presensi_mahasiswa_mahasiswa FOREIGN KEY (mahasiswa_id) REFERENCES users(detail_id) ON UPDATE CASCADE
) engine=InnoDB;

CREATE TABLE IF NOT EXISTS presensi_pegawai (
    presensi_id char(36) NOT NULL,
    pegawai_id char(36) NOT NULL,
    status enum('hadir', 'sakit', 'izin', 'alpha') DEFAULT 'alpha',

    PRIMARY KEY (presensi_id, pegawai_id),
    CONSTRAINT fk_presensi_pegawai_presensi FOREIGN KEY (presensi_id) REFERENCES presensi(id) ON DELETE CASCADE,
    CONSTRAINT fk_presensi_pegawai_pegawai FOREIGN KEY (pegawai_id) REFERENCES users(detail_id) ON UPDATE CASCADE
) engine=InnoDB;

CREATE TABLE IF NOT EXISTS khs (
    id char(36) PRIMARY KEY,
    semester int NOT NULL,
    ips decimal(3,2) default 0,
    ipk decimal(3,2) default 0,
    mahasiswa_id char(36) NOT NULL,

    CONSTRAINT fk_khs_mahasiswa FOREIGN KEY (mahasiswa_id) REFERENCES users(detail_id) ON UPDATE CASCADE,

    UNIQUE INDEX idx_khs_semester_mahasiswa (semester, mahasiswa_id)
) engine=InnoDB;

CREATE TABLE IF NOT EXISTS nilai_mk (
    id char(36) PRIMARY KEY,
    khs_id char(36) NOT NULL,
    total_nilai decimal(5,2) UNSIGNED CHECK (total_nilai >= 0 AND total_nilai <= 100) DEFAULT 0,
    grade_nilai enum('A', 'B', 'C', 'D', 'E'),
    pengampu_id char(36) NOT NULL,

    CONSTRAINT fk_nilai_mk_khs FOREIGN KEY (khs_id) REFERENCES khs(id) ON DELETE CASCADE,
    CONSTRAINT fk_nilai_mk_pengampu FOREIGN KEY (pengampu_id) REFERENCES pengampu(id) ON UPDATE CASCADE,

    UNIQUE INDEX idx_khs_pengampu (khs_id, pengampu_id)
) engine=InnoDB;
