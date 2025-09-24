-- 0006_alamat_table.up.sql
CREATE TABLE IF NOT EXISTS alamat (
  id INT AUTO_INCREMENT PRIMARY KEY,
  id_user INT,
  `judul alamat` VARCHAR(255),
  `nama penerima` VARCHAR(255),
  `no telp` VARCHAR(255),
  detail_alamat VARCHAR(255),
  updated_at DATE,
  created_at DATE,
  CONSTRAINT fk_alamat_user 
    FOREIGN KEY (id_user) REFERENCES users(id)
    ON UPDATE CASCADE ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;