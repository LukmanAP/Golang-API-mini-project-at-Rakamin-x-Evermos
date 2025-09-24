-- 0004_produk_table.up.sql
CREATE TABLE IF NOT EXISTS produk (
  id INT AUTO_INCREMENT PRIMARY KEY,
  nama_produk VARCHAR(255),
  slug VARCHAR(255),
  `harga reseller` VARCHAR(255),
  `harga konsumen` VARCHAR(255),
  stok INT,
  deskripsi TEXT,
  created_at DATE,
  updated_at DATE,
  id_toko INT,
  id_category INT,
  CONSTRAINT fk_produk_toko
    FOREIGN KEY (id_toko) REFERENCES toko(id)
    ON UPDATE CASCADE ON DELETE SET NULL,
  CONSTRAINT fk_produk_category
    FOREIGN KEY (id_category) REFERENCES category(id)
    ON UPDATE CASCADE ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;