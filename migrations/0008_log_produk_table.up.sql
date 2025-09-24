-- 0008_log_produk_table.up.sql
CREATE TABLE IF NOT EXISTS log_produk (
  id INT AUTO_INCREMENT PRIMARY KEY,
  id_produk INT,
  nama_produk VARCHAR(255),
  slug VARCHAR(255),
  `harga reseller` VARCHAR(255),
  `harga konsumen` VARCHAR(255),
  deskripsi TEXT,
  created_at DATE,
  updated_at DATE,
  id_toko INT,
  id_category INT,
  CONSTRAINT fk_log_produk_toko
    FOREIGN KEY (id_toko) REFERENCES toko(id)
    ON UPDATE CASCADE ON DELETE SET NULL,
  CONSTRAINT fk_log_produk_category
    FOREIGN KEY (id_category) REFERENCES category(id)
    ON UPDATE CASCADE ON DELETE SET NULL
  -- Catatan: id_produk tidak dijadikan FK sesuai deskripsi. Tambahkan jika diperlukan.
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;