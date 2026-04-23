CREATE DATABASE if not exists test;

use test

CREATE TABLE if not exists users (
    id binary(16) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    image_url VARCHAR(255),
    role_kode char(4) NOT NULL,
    detail_id int
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE INDEX idx_users_email (email),
    INDEX idx_users_detail (detail_id),
    CONSTRAINT fk_role FOREIGN KEY (role_kode) REFERENCES roles(kode)
);

CREATE TABLE IF NOT EXISTS refresh_token (
    id binary(16) PRIMARY KEY,
    user_id BINARY(16) NOT NULL,
    token VARCHAR(255) NOT NULL,
    expired_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE INDEX idx_token (token),
    INDEX idx_token_user (user_id),
    CONSTRAINT fk_token_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE if not exist role (
    kode char(4) PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);


insert into role (kode, name) values
    ('r-su', 'super admin'),
    ('r-aa', 'admin akademik'),
    ('r-pa', 'admin pegawai'),
    ('r-ka', 'admin keuangan'),
    ('r-ma', 'admin mahasiswa'),
    ('r-ms', 'mahasiswa'),
    ('r-ds', 'dosen');
