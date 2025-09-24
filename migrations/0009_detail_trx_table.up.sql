-- 0009_detail_trx_table.up.sql
CREATE TABLE IF NOT EXISTS detail_trx (
  id INT AUTO_INCREMENT PRIMARY KEY,
  id_trx INT,
  id_log_produk INT,
  id_toko INT,
  kuantitas INT,
  harga_total INT,
  updated_at DATE,
  created_at DATE,
  CONSTRAINT fk_detail_trx_trx
    FOREIGN KEY (id_trx) REFERENCES trx(id)
    ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT fk_detail_trx_log
    FOREIGN KEY (id_log_produk) REFERENCES log_produk(id)
    ON UPDATE CASCADE ON DELETE RESTRICT,
  CONSTRAINT fk_detail_trx_toko
    FOREIGN KEY (id_toko) REFERENCES toko(id)
    ON UPDATE CASCADE ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;