-- 0005_foto_produk_table.up.sql
CREATE TABLE IF NOT EXISTS foto_produk (
  id INT AUTO_INCREMENT PRIMARY KEY,
  id_produk INT,
  url VARCHAR(255),
  updated_at DATE,
  created_at DATE,
  CONSTRAINT fk_foto_produk_produk
    FOREIGN KEY (id_produk) REFERENCES produk(id)
    ON UPDATE CASCADE ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;