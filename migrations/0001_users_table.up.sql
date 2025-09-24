-- 0001_users_table.up.sql
CREATE TABLE IF NOT EXISTS users (
  id INT AUTO_INCREMENT PRIMARY KEY,
  nama VARCHAR(255),
  kata_sandi VARCHAR(255),
  notelp VARCHAR(255) UNIQUE,
  `tanggal lahir` DATE,
  `jenis kelamin` VARCHAR(255),
  tentang TEXT,
  pekerjaan VARCHAR(255),
  email VARCHAR(255) UNIQUE,
  id_provinsi VARCHAR(255),
  id_kota VARCHAR(255),
  isAdmin BOOLEAN,
  updated_at DATE,
  created_at DATE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;